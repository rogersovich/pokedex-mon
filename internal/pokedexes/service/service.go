package service

import (
	"context"
	"fmt"
	"log"
	"pokedex/internal/pokedexes/model"
	"pokedex/internal/pokedexes/repository"
	"pokedex/internal/shared/pokeapi"
	"strconv"
	"sync"
	"time"
)

type PokedexesService interface {
	SyncAllPokedexes(ctx context.Context) error
	GetPokedexDetail(ctx context.Context, identifier string) (model.PokedexesDetail, error)
	GetPokedexList(ctx context.Context, limit, offset int, baseUrl string) (model.ListPokedexes, error)
}

type pokedexesServiceImpl struct {
	pokedexesRepo repository.PokedexesRepository
	pokeAPIClient *pokeapi.Client
}

func NewPokedexesService(repo repository.PokedexesRepository, api *pokeapi.Client) PokedexesService {
	return &pokedexesServiceImpl{
		pokedexesRepo: repo,
		pokeAPIClient: api,
	}
}

func (s *pokedexesServiceImpl) SyncAllPokedexes(ctx context.Context) error {
	log.Println("Starting full Pokedexes data synchronization...")

	limit := 25
	offset := 0
	totalSynced := 0

	for {
		listCtx, cancelList := context.WithTimeout(ctx, 30*time.Second)
		listResponse, err := s.pokeAPIClient.FetchPokedexesList(listCtx, limit, offset)
		cancelList()

		if err != nil {
			if err.Error() == "rate_limit_hit" { // Check for the string error, can be improved with custom error types
				log.Println("Rate limit hit during list fetch, retrying after a delay...")
				time.Sleep(5 * time.Second)
				continue
			}
			return fmt.Errorf("failed to fetch pokedexes list from PokeAPI: %w", err)
		}

		if len(listResponse.Results) == 0 {
			break // No more data to fetch
		}

		var wg sync.WaitGroup
		resultsChan := make(chan struct {
			Detail model.PokedexesDetail
			Err    error
		}, len(listResponse.Results))

		// Enqueue each detail fetch through the shared client
		for _, item := range listResponse.Results {
			wg.Add(1)
			// Fungsi goroutine ini bertugas untuk mengambil detail data pokedex berdasarkan URL yang diberikan.
			// Setiap item pada batch list akan diproses secara paralel (concurrent) menggunakan goroutine.
			// Setelah detail berhasil diambil (atau terjadi error), hasilnya dikirim ke channel resultsChan
			// dalam bentuk struct yang berisi data detail dan error (jika ada).
			go func(item model.ListPokedexesItem) {
				defer wg.Done()
				detail, err := s.pokeAPIClient.FetchPokedexesDetail(ctx, item.URL)
				resultsChan <- struct {
					Detail model.PokedexesDetail
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
				log.Printf("Error fetching detail for a pokedexes: %v\n", res.Err)
				if res.Err.Error() == "rate_limit_hit" {
					log.Println("Rate limit hit during detail fetch, consider re-queuing or pausing sync.")
				}
				continue
			}

			// Save to Repository
			err := s.pokedexesRepo.SavePokedexes(ctx, res.Detail)
			if err != nil {
				log.Printf("Failed to save pokedexes %s (ID: %d) to repository: %v\n", res.Detail.Name, res.Detail.ID, err)
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

func (s *pokedexesServiceImpl) GetPokedexDetail(ctx context.Context, identifier string) (model.PokedexesDetail, error) {
	id, err := strconv.Atoi(identifier)

	if err != nil {
		return model.PokedexesDetail{}, fmt.Errorf("invalid id: %s must number", identifier)
	}

	return s.pokedexesRepo.GetPokedexByID(ctx, id)
}

func (s *pokedexesServiceImpl) GetPokedexList(ctx context.Context, limit, offset int, baseUrl string) (model.ListPokedexes, error) {
	var list_types []model.ListPokedexesItem
	var totalCount int64
	var err error

	list_types, totalCount, err = s.pokedexesRepo.GetPokedexList(ctx, limit, offset, baseUrl)

	if err != nil {
		return model.ListPokedexes{}, err
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
		list_types = make([]model.ListPokedexesItem, 0)
	}

	return model.ListPokedexes{
		Count:    int(totalCount),
		Next:     nextURL,
		Previous: previousURL,
		Results:  list_types,
	}, nil
}
