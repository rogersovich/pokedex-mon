package service

import (
	"context"
	"fmt"
	"log"
	"pokedex/internal/items/model"
	"pokedex/internal/items/repository"
	"pokedex/internal/shared/pokeapi"
	"pokedex/utils"
	"strconv"
	"sync"
	"time"
)

type ItemsService interface {
	SyncAllItems(ctx context.Context) error
	GetItemDetail(ctx context.Context, identifier string, baseUrl string) (model.ItemDetail, *utils.BaseResourceNavigation, *utils.BaseResourceNavigation, error)
	GetItemList(ctx context.Context, limit, offset int, baseUrl string) (model.ListItem, error)
}

type itemsServiceImpl struct {
	itemsRepo     repository.ItemsRepository
	pokeAPIClient *pokeapi.Client
}

func NewItemsService(repo repository.ItemsRepository, api *pokeapi.Client) ItemsService {
	return &itemsServiceImpl{
		itemsRepo:     repo,
		pokeAPIClient: api,
	}
}

func (s *itemsServiceImpl) SyncAllItems(ctx context.Context) error {
	log.Println("Starting full Items data synchronization...")

	limit := 50
	offset := 0
	totalSynced := 0

	for {
		listCtx, cancelList := context.WithTimeout(ctx, 30*time.Second)
		listResponse, err := s.pokeAPIClient.FetchItemList(listCtx, limit, offset)
		cancelList()

		if err != nil {
			if err.Error() == "rate_limit_hit" { // Check for the string error, can be improved with custom error types
				log.Println("Rate limit hit during list fetch, retrying after a delay...")
				time.Sleep(5 * time.Second)
				continue
			}
			return fmt.Errorf("failed to fetch items list from PokeAPI: %w", err)
		}

		if len(listResponse.Results) == 0 {
			break // No more data to fetch
		}

		var wg sync.WaitGroup
		resultsChan := make(chan struct {
			Detail model.ItemDetail
			Err    error
		}, len(listResponse.Results))

		// Enqueue each detail fetch through the shared client
		for _, item := range listResponse.Results {
			wg.Add(1)
			go func(item model.ListItemItem) {
				defer wg.Done()
				detail, err := s.pokeAPIClient.FetchItemDetail(ctx, item.URL)
				resultsChan <- struct {
					Detail model.ItemDetail
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
				log.Printf("Error fetching detail for a item: %v\n", res.Err)
				if res.Err.Error() == "rate_limit_hit" {
					log.Println("Rate limit hit during detail fetch, consider re-queuing or pausing sync.")
				}
				continue
			}

			// Save to Repository
			err := s.itemsRepo.SaveItem(ctx, res.Detail)
			if err != nil {
				log.Printf("Failed to save item %s (ID: %d) to repository: %v\n", res.Detail.Name, res.Detail.ID, err)
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

func (s *itemsServiceImpl) GetItemDetail(ctx context.Context, identifier string, baseUrl string) (data model.ItemDetail, next, prev *utils.BaseResourceNavigation, err error) {
	var nextData *utils.BaseResourceNavigation
	var prevData *utils.BaseResourceNavigation

	id, err := strconv.Atoi(identifier)

	if err == nil {
		data, _ = s.itemsRepo.GetItemByID(ctx, id)
	} else {
		data, _ = s.itemsRepo.GetItemByName(ctx, identifier)
	}

	_, totalCount, err := s.itemsRepo.GetItemList(ctx, 1, 0, baseUrl)

	if err != nil {
		return model.ItemDetail{}, nil, nil, err
	}

	totalCountInt := int(totalCount)

	if data.ID >= totalCountInt {
		nextData = nil
	} else {
		nextID := data.ID + 1
		nextURL := fmt.Sprintf("%s/%s", baseUrl, strconv.Itoa(nextID))

		getNextData, _ := s.itemsRepo.GetItemByID(ctx, nextID)

		nextData = &utils.BaseResourceNavigation{
			ID:   nextID,
			Name: getNextData.Name,
			URL:  nextURL,
		}
	}

	if data.ID <= 1 {
		prevData = nil
	} else {
		prevID := data.ID - 1
		prevURL := fmt.Sprintf("%s/%s", baseUrl, strconv.Itoa(prevID))

		getPrevData, _ := s.itemsRepo.GetItemByID(ctx, prevID)

		prevData = &utils.BaseResourceNavigation{
			ID:   prevID,
			Name: getPrevData.Name,
			URL:  prevURL,
		}
	}

	return data, nextData, prevData, nil
}

func (s *itemsServiceImpl) GetItemList(ctx context.Context, limit, offset int, baseUrl string) (model.ListItem, error) {
	var data []model.ListItemItem
	var totalCount int64
	var err error

	data, totalCount, err = s.itemsRepo.GetItemList(ctx, limit, offset, baseUrl)

	if err != nil {
		return model.ListItem{}, err
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
	if data == nil {
		data = make([]model.ListItemItem, 0)
	}

	return model.ListItem{
		Count:    int(totalCount),
		Next:     nextURL,
		Previous: previousURL,
		Results:  data,
	}, nil
}
