package service

import (
	"context"
	"fmt"
	"log"
	"pokedex/internal/pokemon-type/model"
	"pokedex/internal/pokemon-type/repository"
	"pokedex/internal/shared/pokeapi"
	"strconv"
	"sync"
	"time"
)

type PokemonTypeService interface {
	SyncAllPokemonType(ctx context.Context) error
	GetPokemonType(ctx context.Context, identifier string) (model.PokemonTypeDetailResponse, error)
	GetPokemonTypeList(ctx context.Context, limit, offset int, baseUrl string) (model.PokemonListTypeResponse, error)
}

type pokemonTypeServiceImpl struct {
	pokemonTypeRepo repository.PokemonTypeRepository
	pokeAPIClient   *pokeapi.Client
}

func NewPokemonTypeService(repo repository.PokemonTypeRepository, api *pokeapi.Client) PokemonTypeService {
	return &pokemonTypeServiceImpl{
		pokemonTypeRepo: repo,
		pokeAPIClient:   api,
	}
}

func (s *pokemonTypeServiceImpl) SyncAllPokemonType(ctx context.Context) error {
	log.Println("Starting full data synchronization...")

	limit := 50 // Fetch 100 pokemons at a time from PokeAPI
	offset := 0
	totalSynced := 0

	for {
		listCtx, cancelList := context.WithTimeout(ctx, 30*time.Second)
		listResponse, err := s.pokeAPIClient.FetchPokemonTypeList(listCtx, limit, offset)
		cancelList()

		if err != nil {
			if err.Error() == "rate_limit_hit" { // Check for the string error, can be improved with custom error types
				log.Println("Rate limit hit during list fetch, retrying after a delay...")
				time.Sleep(5 * time.Second)
				continue
			}
			return fmt.Errorf("failed to fetch data list from PokeAPI: %w", err)
		}

		if len(listResponse.Results) == 0 {
			break // No more data to fetch
		}

		var wg sync.WaitGroup
		resultsChan := make(chan struct {
			Detail model.PokemonTypeDetailResponse
			Err    error
		}, len(listResponse.Results))

		// Enqueue each detail fetch through the shared client
		for _, item := range listResponse.Results {
			wg.Add(1)
			go func(item model.PokemonTypeListItem) {
				defer wg.Done()
				detail, err := s.pokeAPIClient.FetchPokemonTypeDetail(ctx, item.URL)
				resultsChan <- struct {
					Detail model.PokemonTypeDetailResponse
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
				log.Printf("Error fetching detail for a data: %v\n", res.Err)
				if res.Err.Error() == "rate_limit_hit" {
					log.Println("Rate limit hit during detail fetch, consider re-queuing or pausing sync.")
				}
				continue
			}

			// Save to Repository
			err := s.pokemonTypeRepo.SavePokemonType(ctx, res.Detail)
			if err != nil {
				log.Printf("Failed to save data (ID: %d) to repository: %v\n", res.Detail.TypeID, err)
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

	log.Printf("Full Type Pokémon data synchronization completed. Total unique data synced: %d\n", totalSynced)
	return nil
}

func (s *pokemonTypeServiceImpl) GetPokemonType(ctx context.Context, identifier string) (model.PokemonTypeDetailResponse, error) {
	id, err := strconv.Atoi(identifier)

	var res model.PokemonTypeDetailResponse
	if err == nil {
		res, _ = s.pokemonTypeRepo.GetPokemonTypeByID(ctx, id)
	} else {
		res, _ = s.pokemonTypeRepo.GetPokemonTypeByName(ctx, identifier)
	}

	return res, err
}

func (s *pokemonTypeServiceImpl) GetPokemonTypeList(ctx context.Context, limit, offset int, baseUrl string) (model.PokemonListTypeResponse, error) {
	var list_types []model.PokemonTypeListItem
	var totalCount int64
	var err error

	list_types, totalCount, err = s.pokemonTypeRepo.GetPokemonTypeList(ctx, limit, offset, baseUrl)

	if err != nil {
		return model.PokemonListTypeResponse{}, err
	}

	// --- LOGIKA PEMBANGUNAN URL NEXT DAN PREVIOUS ---
	var nextURL *string
	var previousURL *string

	// Next URL
	if offset+limit < int(totalCount) {
		nextOffset := offset + limit
		url := fmt.Sprintf("%s?limit=%d&offset=%d", baseUrl, limit, nextOffset)
		nextURL = &url
	}

	// Previous URL
	if offset > 0 {
		prevOffset := offset - limit
		if prevOffset < 0 {
			prevOffset = 0 // Pastikan offset tidak negatif
		}
		url := fmt.Sprintf("%s?limit=%d&offset=%d", baseUrl, limit, prevOffset)
		previousURL = &url
	}
	// --- AKHIR LOGIKA PEMBANGUNAN URL NEXT DAN PREVIOUS ---

	// Ensure Results is an empty slice (not nil) if there are no items
	if list_types == nil {
		list_types = make([]model.PokemonTypeListItem, 0)
	}

	return model.PokemonListTypeResponse{
		Count:    int(totalCount),
		Next:     nextURL,
		Previous: previousURL,
		Results:  list_types,
	}, nil
}
