package service

import (
	"context"
	"fmt"
	"log"
	"pokedex/internal/generation/model"
	"pokedex/internal/generation/repository"
	"pokedex/internal/shared/pokeapi"
	"sync"
	"time"
)

type GenerationService interface {
	SyncAllGeneration(ctx context.Context) error
}

type generationServiceImpl struct {
	generationRepo repository.GenerationRepository
	pokeAPIClient  *pokeapi.Client
}

func NewGenerationService(repo repository.GenerationRepository, api *pokeapi.Client) GenerationService {
	return &generationServiceImpl{
		generationRepo: repo,
		pokeAPIClient:  api,
	}
}

func (s *generationServiceImpl) SyncAllGeneration(ctx context.Context) error {
	log.Println("Starting full generation data synchronization...")

	limit := 25
	offset := 0
	totalSynced := 0

	for {
		listCtx, cancelList := context.WithTimeout(ctx, 30*time.Second)
		listResponse, err := s.pokeAPIClient.FetchGenerationList(listCtx, limit, offset)
		cancelList()

		if err != nil {
			if err.Error() == "rate_limit_hit" { // Check for the string error, can be improved with custom error types
				log.Println("Rate limit hit during list fetch, retrying after a delay...")
				time.Sleep(5 * time.Second)
				continue
			}
			return fmt.Errorf("failed to fetch generation list from PokeAPI: %w", err)
		}

		if len(listResponse.Results) == 0 {
			break // No more data to fetch
		}

		var wg sync.WaitGroup
		resultsChan := make(chan struct {
			Detail model.GenerationDetail
			Err    error
		}, len(listResponse.Results))

		// Enqueue each detail fetch through the shared client
		for _, item := range listResponse.Results {
			wg.Add(1)
			go func(item model.ListGenerationItem) {
				defer wg.Done()
				detail, err := s.pokeAPIClient.FetchGenerationDetail(ctx, item.URL)
				resultsChan <- struct {
					Detail model.GenerationDetail
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
				log.Printf("Error fetching detail for a generation: %v\n", res.Err)
				if res.Err.Error() == "rate_limit_hit" {
					log.Println("Rate limit hit during detail fetch, consider re-queuing or pausing sync.")
				}
				continue
			}

			// Save to Repository
			err := s.generationRepo.SaveGenerations(ctx, res.Detail)
			if err != nil {
				log.Printf("Failed to save generation %s (ID: %d) to repository: %v\n", res.Detail.Name, res.Detail.ID, err)
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
