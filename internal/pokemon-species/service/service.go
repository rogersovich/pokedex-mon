package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"pokedex/internal/pokemon-species/model"
	"pokedex/internal/pokemon-species/repository"
	"pokedex/internal/shared/pokeapi"
)

type PokemonSpeciesService interface {
	SyncAllPokemonSpecies(ctx context.Context) error
	GetPokemonSpecies(ctx context.Context, identifier string) (model.PokemonSpeciesDetail, error)
}

type pokemonSpeciesServiceImpl struct {
	pokemonSpeciesRepo repository.PokemonSpeciesRepository
	pokeAPIClient      *pokeapi.Client
}

func NewPokemonSpeciesService(repo repository.PokemonSpeciesRepository, api *pokeapi.Client) PokemonSpeciesService {
	return &pokemonSpeciesServiceImpl{
		pokemonSpeciesRepo: repo,
		pokeAPIClient:      api,
	}
}

// SyncAllPokemonSpecies fetches all pokemon list and their details from PokeAPI
// and stores them in the local repository. This should be run as a background job.
func (s *pokemonSpeciesServiceImpl) SyncAllPokemonSpecies(ctx context.Context) error {
	log.Println("Starting full data synchronization...")

	limit := 100 // Fetch 100 pokemons at a time from PokeAPI
	offset := 0
	totalSynced := 0

	for {
		listCtx, cancelList := context.WithTimeout(ctx, 30*time.Second)
		listResponse, err := s.pokeAPIClient.FetchPokemonSpeciesList(listCtx, limit, offset)
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
			Detail model.PokemonSpeciesDetail
			Err    error
		}, len(listResponse.Results))

		// Enqueue each detail fetch through the shared client
		for _, item := range listResponse.Results {
			wg.Add(1)
			go func(item model.PokemonSpeciesListItem) {
				defer wg.Done()
				detail, err := s.pokeAPIClient.FetchPokemonSpeciesDetail(ctx, item.URL)
				resultsChan <- struct {
					Detail model.PokemonSpeciesDetail
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
			err := s.pokemonSpeciesRepo.SavePokemonSpecies(ctx, res.Detail)
			if err != nil {
				log.Printf("Failed to save data %s (ID: %d) to repository: %v\n", res.Detail.Name, res.Detail.PokeAPIID, err)
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

	log.Printf("Full Pokémon data synchronization completed. Total unique data synced: %d\n", totalSynced)
	return nil
}

func (s *pokemonSpeciesServiceImpl) GetPokemonSpecies(ctx context.Context, identifier string) (model.PokemonSpeciesDetail, error) {
	id, err := strconv.Atoi(identifier)
	if err == nil {
		return s.pokemonSpeciesRepo.GetPokemonSpeciesByID(ctx, id)
	} else {
		return s.pokemonSpeciesRepo.GetPokemonSpeciesByName(ctx, identifier)
	}
}
