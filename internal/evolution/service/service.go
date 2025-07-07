package service

import (
	"context"
	"fmt"
	"log"
	"pokedex/internal/evolution/model"
	"pokedex/internal/evolution/repository"
	"pokedex/internal/shared/pokeapi"
	"pokedex/utils"
	"strconv"
	"sync"
	"time"
)

type EvolutionService interface {
	SyncAllEvolution(ctx context.Context) error
	GetEvolution(ctx context.Context, identifier string) (model.EvolutionChain, error)
	GetEvolutionPokemonType(ctx context.Context, id int) (model.EvolutionPokemonResponse, error)
}

type evolutionServiceImpl struct {
	evolutionRepo repository.EvolutionRepository
	pokeAPIClient *pokeapi.Client
}

func NewEvolutionService(repo repository.EvolutionRepository, api *pokeapi.Client) EvolutionService {
	return &evolutionServiceImpl{
		evolutionRepo: repo,
		pokeAPIClient: api,
	}
}

func (s *evolutionServiceImpl) SyncAllEvolution(ctx context.Context) error {
	log.Println("Starting full data synchronization...")

	limit := 100 // Fetch 100 pokemons at a time from PokeAPI
	offset := 0
	totalSynced := 0

	for {
		listCtx, cancelList := context.WithTimeout(ctx, 30*time.Second)
		listResponse, err := s.pokeAPIClient.FetchEvolutionList(listCtx, limit, offset)
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
			Detail model.EvolutionChain
			Err    error
		}, len(listResponse.Results))

		// Enqueue each detail fetch through the shared client
		for _, item := range listResponse.Results {
			wg.Add(1)
			go func(item model.EvolutionListItem) {
				defer wg.Done()
				detail, err := s.pokeAPIClient.FetchEvolutionDetail(ctx, item.URL)
				resultsChan <- struct {
					Detail model.EvolutionChain
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
			err := s.evolutionRepo.SaveEvolution(ctx, res.Detail)
			if err != nil {
				log.Printf("Failed to save data (ID: %d) to repository: %v\n", res.Detail.ID, err)
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

func (s *evolutionServiceImpl) GetEvolution(ctx context.Context, identifier string) (model.EvolutionChain, error) {
	id, err := strconv.Atoi(identifier)

	var evoDetail model.EvolutionChain
	if err == nil {
		evoDetail, _ = s.evolutionRepo.GetEvolutionByID(ctx, id)
	} else {
		evoDetail, _ = s.evolutionRepo.GetEvolutionByName(ctx, identifier)
	}

	err = s.populateEvolutionChainDetails(ctx, &evoDetail.Chain)
	if err != nil {
		return model.EvolutionChain{}, fmt.Errorf("failed to populate evolution chain details: %w", err)
	}

	return evoDetail, err
}

func (s *evolutionServiceImpl) GetEvolutionPokemonType(ctx context.Context, pokemon_id int) (model.EvolutionPokemonResponse, error) {
	return s.evolutionRepo.GetEvolutionPokemonType(ctx, pokemon_id)
}

func (s *evolutionServiceImpl) GetPokemonInfo(ctx context.Context, pokemon_name string) (model.EvolutionPokemonInfo, error) {
	return s.evolutionRepo.GetPokemonInfo(ctx, pokemon_name)
}

func (s *evolutionServiceImpl) populateEvolutionChainDetails(ctx context.Context, chainLink *model.ChainLink) error {
	if chainLink == nil || chainLink.Species.URL == "" {
		return nil
	}

	//? 1. Get Pokemon Name and Set Pokemon Info
	pokemonName := chainLink.Species.Name
	pokemonInfoResponse, _ := s.evolutionRepo.GetPokemonInfo(ctx, pokemonName)
	pokemonID := pokemonInfoResponse.ID
	thumbnailImg := utils.GetThumbnailPokemon(pokemonID)

	chainLink.PokemonInfo = model.EvolutionPokemonInfoResponse{
		ID:        pokemonID,
		Name:      pokemonName,
		Thumbnail: thumbnailImg,
	}

	// 2. Fetch the EvolutionPokemonResponse using the extracted ID from your repository.
	evolutionPokemonResponse, err := s.evolutionRepo.GetEvolutionPokemonType(ctx, pokemonID)
	if err != nil {
		fmt.Printf("Warning: Failed to get EvolutionPokemonType for ID %d (Species: %s): %v\n", pokemonID, chainLink.Species.Name, err)
		chainLink.EvolutionType = model.EvolutionPokemonResponse{} // Assign empty struct
		return fmt.Errorf("failed to fetch evolution pokemon type for ID %d: %w", pokemonID, err)
	}

	// 3. Populate the EvolutionType field in the current ChainLink node.
	chainLink.EvolutionType = evolutionPokemonResponse

	// 4. Recursively call this function for each subsequent evolution in EvolvesTo.
	for i := range chainLink.EvolvesTo {
		// Pass the address of the element in the slice so modifications are applied directly.
		if err := s.populateEvolutionChainDetails(ctx, &chainLink.EvolvesTo[i]); err != nil {
			// If an error in a sub-chain should stop the entire process, return it.
			// Otherwise, log and continue.
			fmt.Printf("Error processing sub-evolution chain link for species %s: %v\n", chainLink.EvolvesTo[i].Species.Name, err)
			return err // Returning error here will stop the entire chain population on the first sub-error
		}
	}

	return nil // No error encountered in this branch
}
