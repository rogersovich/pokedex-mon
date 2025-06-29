package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"pokedex/internal/pokemon/model"
	"pokedex/internal/pokemon/repository"
	"pokedex/internal/shared/pokeapi"
)

// PokemonService defines the business logic for Pokemon operations.
type PokemonService interface {
	SyncAllPokemons(ctx context.Context) error
	GetPokemon(ctx context.Context, identifier string) (model.PokemonDetailResponse, error)
	GetPokemonList(ctx context.Context, limit, offset int, baseUrl string, searchQuery string) (model.PokemonListResponse, error)
}

// pokemonServiceImpl implements the PokemonService interface.
type pokemonServiceImpl struct {
	pokemonRepo   repository.PokemonRepository
	pokeAPIClient *pokeapi.Client // Use the shared client
}

// NewPokemonService creates a new instance of PokemonService.
func NewPokemonService(repo repository.PokemonRepository, api *pokeapi.Client) PokemonService {
	return &pokemonServiceImpl{
		pokemonRepo:   repo,
		pokeAPIClient: api,
	}
}

// SyncAllPokemons fetches all pokemon list and their details from PokeAPI
// and stores them in the local repository. This should be run as a background job.
func (s *pokemonServiceImpl) SyncAllPokemons(ctx context.Context) error {
	log.Println("Starting full Pokémon data synchronization...")

	limit := 100 // Fetch 100 pokemons at a time from PokeAPI
	offset := 0
	totalSynced := 0

	for {
		listCtx, cancelList := context.WithTimeout(ctx, 30*time.Second)
		listResponse, err := s.pokeAPIClient.FetchPokemonList(listCtx, limit, offset)
		cancelList()

		if err != nil {
			if err.Error() == "rate_limit_hit" { // Check for the string error, can be improved with custom error types
				log.Println("Rate limit hit during list fetch, retrying after a delay...")
				time.Sleep(5 * time.Second)
				continue
			}
			return fmt.Errorf("failed to fetch pokemon list from PokeAPI: %w", err)
		}

		if len(listResponse.Results) == 0 {
			break // No more pokemons to fetch
		}

		var wg sync.WaitGroup
		resultsChan := make(chan struct {
			Detail model.PokemonDetail
			Err    error
		}, len(listResponse.Results))

		// Enqueue each detail fetch through the shared client
		for _, item := range listResponse.Results {
			wg.Add(1)
			go func(item model.PokemonListItem) {
				defer wg.Done()
				detail, err := s.pokeAPIClient.FetchPokemonDetail(ctx, item.URL)
				resultsChan <- struct {
					Detail model.PokemonDetail
					Err    error
				}{Detail: detail, Err: err}
			}(item)
		}

		// Wait for all detail fetches for the current batch to complete
		go func() {
			wg.Wait()
			close(resultsChan)
		}()

		for res := range resultsChan {
			if res.Err != nil {
				log.Printf("Error fetching detail for a pokemon: %v\n", res.Err)
				if res.Err.Error() == "rate_limit_hit" {
					log.Println("Rate limit hit during detail fetch, consider re-queuing or pausing sync.")
				}
				continue
			}

			// Save to Repository
			err := s.pokemonRepo.SavePokemon(ctx, res.Detail)
			if err != nil {
				log.Printf("Failed to save Pokemon %s (ID: %d) to repository: %v\n", res.Detail.Name, res.Detail.ID, err)
			} else {
				totalSynced++
			}
		}

		log.Printf("Batch processed. Total synced so far: %d\n", totalSynced)

		offset += limit
		if offset >= listResponse.Count {
			break
		}
	}

	log.Printf("Full Pokémon data synchronization completed. Total unique pokemons synced: %d\n", totalSynced)
	return nil
}

// GetPokemon retrieves a single pokemon detail by ID or Name from the local repository.
func (s *pokemonServiceImpl) GetPokemon(ctx context.Context, identifier string) (model.PokemonDetailResponse, error) {
	id, err := strconv.Atoi(identifier)
	if err == nil {
		return s.pokemonRepo.GetPokemonByID(ctx, id)
	} else {
		return s.pokemonRepo.GetPokemonByName(ctx, identifier)
	}
}

func (s *pokemonServiceImpl) GetPokemonList(ctx context.Context, limit, offset int, baseUrl string, searchQuery string) (model.PokemonListResponse, error) {
	var pokemons []model.PokemonDetail
	var totalCount int64
	var err error

	if searchQuery != "" {
		pokemons, totalCount, err = s.pokemonRepo.SearchPokemons(ctx, searchQuery, limit, offset)
	} else {
		pokemons, totalCount, err = s.pokemonRepo.GetPokemonList(ctx, limit, offset)
	}

	if err != nil {
		return model.PokemonListResponse{}, err
	}

	var listItems []model.PokemonListItem
	for _, p := range pokemons {
		listItems = append(listItems, model.PokemonListItem{
			Name:    p.Name,
			URL:     fmt.Sprintf("%s/%d", baseUrl, p.ID),
			Types:   p.Types,
			Sprites: p.Sprites,
		})
	}

	// --- LOGIKA PEMBANGUNAN URL NEXT DAN PREVIOUS ---
	var nextURL *string
	var previousURL *string

	// Next URL
	if offset+limit < int(totalCount) {
		nextOffset := offset + limit
		url := fmt.Sprintf("%s?limit=%d&offset=%d", baseUrl, limit, nextOffset)
		if searchQuery != "" {
			url = fmt.Sprintf("%s&q=%s", url, searchQuery) // Tambahkan q parameter
		}
		nextURL = &url
	}

	// Previous URL
	if offset > 0 {
		prevOffset := offset - limit
		if prevOffset < 0 {
			prevOffset = 0 // Pastikan offset tidak negatif
		}
		url := fmt.Sprintf("%s?limit=%d&offset=%d", baseUrl, limit, prevOffset)
		if searchQuery != "" {
			url = fmt.Sprintf("%s&q=%s", url, searchQuery) // Tambahkan q parameter
		}
		previousURL = &url
	}
	// --- AKHIR LOGIKA PEMBANGUNAN URL NEXT DAN PREVIOUS ---

	return model.PokemonListResponse{
		Count:    int(totalCount),
		Next:     nextURL,
		Previous: previousURL,
		Results:  listItems,
	}, nil
}
