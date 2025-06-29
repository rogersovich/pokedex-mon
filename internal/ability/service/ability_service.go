package service

import (
	"context"
	"fmt"
	"log"
	"pokedex/internal/ability/model"
	"pokedex/internal/ability/repository"
	"pokedex/internal/shared/pokeapi"
	"strconv"
	"strings"
	"sync"
	"time"
)

// AbilityService defines the business logic for Ability operations.
type AbilityService interface {
	SyncAllAbilities(ctx context.Context) error
	GetAbility(ctx context.Context, identifier string) (model.AbilityDetail, error) // Untuk mengambil dari DB
}

// abilityServiceImpl implements the AbilityService interface.
type abilityServiceImpl struct {
	abilityRepo   repository.AbilityRepository
	pokeAPIClient *pokeapi.Client
}

// NewAbilityService creates a new instance of AbilityService.
func NewAbilityService(repo repository.AbilityRepository, api *pokeapi.Client) AbilityService {
	return &abilityServiceImpl{
		abilityRepo:   repo,
		pokeAPIClient: api,
	}
}

// SyncAllAbilities fetches all abilities from PokeAPI and saves them to the repository.
func (s *abilityServiceImpl) SyncAllAbilities(ctx context.Context) error {
	log.Println("Starting full Ability data synchronization...")

	limit := 50
	offset := 0
	totalSynced := 0

	for {
		listCtx, cancelList := context.WithTimeout(ctx, 30*time.Second)
		listResponse, err := s.pokeAPIClient.FetchAbilityList(listCtx, limit, offset)
		cancelList()

		if err != nil {
			if err.Error() == "rate_limit_hit" { // Check for the string error, can be improved with custom error types
				log.Println("Rate limit hit during list fetch, retrying after a delay...")
				time.Sleep(5 * time.Second)
				continue
			}
			return fmt.Errorf("failed to fetch ability list from PokeAPI: %w", err)
		}

		if len(listResponse.Results) == 0 {
			break // No more data to fetch
		}

		var wg sync.WaitGroup
		resultsChan := make(chan struct {
			Detail model.AbilityDetail
			Err    error
		}, len(listResponse.Results))

		// Enqueue each detail fetch through the shared client
		for _, item := range listResponse.Results {
			wg.Add(1)
			go func(item model.AbilityListItem) {
				defer wg.Done()
				detail, err := s.pokeAPIClient.FetchAbilityDetail(ctx, item.URL)
				resultsChan <- struct {
					Detail model.AbilityDetail
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
				log.Printf("Error fetching detail for a ability: %v\n", res.Err)
				if res.Err.Error() == "rate_limit_hit" {
					log.Println("Rate limit hit during detail fetch, consider re-queuing or pausing sync.")
				}
				continue
			}

			// Save to Repository
			err := s.abilityRepo.SaveAbility(ctx, res.Detail)
			if err != nil {
				log.Printf("Failed to save ability %s (ID: %d) to repository: %v\n", res.Detail.Name, res.Detail.ID, err)
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

	log.Printf("Full data synchronization completed. Total unique abilities synced: %d\n", totalSynced)
	return nil
}

// GetAbility retrieves an ability by ID or name from the repository.
func (s *abilityServiceImpl) GetAbility(ctx context.Context, identifier string) (model.AbilityDetail, error) {
	id, err := strconv.Atoi(identifier)
	if err == nil {
		return s.abilityRepo.GetAbilityByID(ctx, id)
	} else {
		return s.abilityRepo.GetAbilityByName(ctx, strings.ToLower(identifier))
	}
}
