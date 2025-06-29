package pokeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"pokedex/config"
	modelability "pokedex/internal/ability/model"
	modelevolution "pokedex/internal/evolution/model"
	modelpokemonspecies "pokedex/internal/pokemon-species/model"
	modelpokemon "pokedex/internal/pokemon/model"
)

var (
	pokeAPIBaseURL string
	httpClient     *http.Client
	requestQueue   chan func()
	wg             sync.WaitGroup
	once           sync.Once
	rateLimitMs    time.Duration
)

type Client struct{} // Our shared PokeAPI Client

func NewClient(cfg *config.Config) *Client {
	once.Do(func() {
		pokeAPIBaseURL = cfg.PokeAPIURL
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}
		requestQueue = make(chan func(), 100)
		rateLimitMs = time.Duration(cfg.RateLimitMs) * time.Millisecond

		go processQueue()
	})
	return &Client{}
}

func processQueue() {
	ticker := time.NewTicker(rateLimitMs)
	defer ticker.Stop()

	for {
		select {
		case task, ok := <-requestQueue:
			if !ok {
				log.Println("PokeAPI request queue closed. Stopping processor.")
				return
			}
			task()
			<-ticker.C
		}
	}
}

func (c *Client) CloseClient() {
	close(requestQueue)
	wg.Wait()
	log.Println("All queued PokeAPI requests processed and client closed.")
}

func (c *Client) enqueueAndFetch(ctx context.Context, url string, target interface{}) error {
	resultChan := make(chan error, 1)
	wg.Add(1)
	requestQueue <- func() {
		defer wg.Done()
		err := c.fetchAndDecode(ctx, url, target)
		resultChan <- err
	}

	select {
	case err := <-resultChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Client) fetchAndDecode(ctx context.Context, url string, target interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.Canceled {
			return ctx.Err()
		}
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		waitDuration := 3 * time.Second
		if sleepTime, parseErr := strconv.Atoi(retryAfter); parseErr == nil && sleepTime > 0 {
			waitDuration = time.Duration(sleepTime) * time.Second
		}
		log.Printf("Rate limit (429) hit for %s. Waiting %v before potentially retrying.\n", url, waitDuration)
		return fmt.Errorf("rate_limit_hit") // Use string for error since we don't have domain errors here. Can be improved.
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API responded with status %d: %s", resp.StatusCode, resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return nil
}

// --- FUNGSI BARU UNTUK POKEMON ---

// FetchPokemonList fetches a list of Pokémon from PokeAPI.
func (c *Client) FetchPokemonList(ctx context.Context, limit, offset int) (modelpokemon.PokemonListResponse, error) {
	url := pokeAPIBaseURL + "/pokemon?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset)
	log.Printf("Enqueueing list fetch from PokeAPI: %s\n", url)

	var response modelpokemon.PokemonListResponse
	err := c.enqueueAndFetch(ctx, url, &response)
	return response, err
}

// FetchPokemonDetail fetches details for a single Pokemon from PokeAPI.
func (c *Client) FetchPokemonDetail(ctx context.Context, url string) (modelpokemon.PokemonDetail, error) {
	log.Printf("Enqueueing detail fetch from PokeAPI: %s\n", url)

	var response modelpokemon.PokemonDetail
	err := c.enqueueAndFetch(ctx, url, &response)
	return response, err
}

// --- FUNGSI BARU UNTUK POKEMON SPECIES ---

// FetchPokemonList fetches a list of Pokémon from PokeAPI.
func (c *Client) FetchPokemonSpeciesList(ctx context.Context, limit, offset int) (modelpokemonspecies.PokemonSpeciesListResponse, error) {
	url := pokeAPIBaseURL + "/pokemon-species?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset)
	log.Printf("Enqueueing list fetch from PokeAPI: %s\n", url)

	var response modelpokemonspecies.PokemonSpeciesListResponse
	err := c.enqueueAndFetch(ctx, url, &response)
	return response, err
}

// FetchPokemonDetail fetches details for a single Pokemon from PokeAPI.
func (c *Client) FetchPokemonSpeciesDetail(ctx context.Context, url string) (modelpokemonspecies.PokemonSpeciesDetail, error) {
	log.Printf("Enqueueing detail fetch from PokeAPI: %s\n", url)

	var response modelpokemonspecies.PokemonSpeciesDetail
	err := c.enqueueAndFetch(ctx, url, &response)
	return response, err
}

// --- FUNGSI BARU UNTUK ABILITY ---

// FetchAbilityList fetches a paginated list of abilities.
func (c *Client) FetchAbilityList(ctx context.Context, limit, offset int) (modelability.AbilityListResponse, error) {
	url := pokeAPIBaseURL + "/ability?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset)
	log.Printf("Enqueueing list fetch from PokeAPI: %s\n", url)

	var response modelability.AbilityListResponse
	err := c.enqueueAndFetch(ctx, url, &response)
	return response, err
}

// FetchAbilityDetail fetches a single ability detail by its URL.
func (c *Client) FetchAbilityDetail(ctx context.Context, url string) (modelability.AbilityDetail, error) {
	log.Printf("Enqueueing detail fetch from PokeAPI: %s\n", url)

	var response modelability.AbilityDetail
	err := c.enqueueAndFetch(ctx, url, &response)
	return response, err
}

// --- FUNGSI BARU UNTUK EVOLUTION ---

func (c *Client) FetchEvolutionList(ctx context.Context, limit, offset int) (modelevolution.EvolutionListResponse, error) {
	url := pokeAPIBaseURL + "/evolution-chain?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset)
	log.Printf("Enqueueing list fetch from PokeAPI: %s\n", url)

	var response modelevolution.EvolutionListResponse
	err := c.enqueueAndFetch(ctx, url, &response)
	return response, err
}

func (c *Client) FetchEvolutionDetail(ctx context.Context, url string) (modelevolution.EvolutionChain, error) {
	log.Printf("Enqueueing detail fetch from PokeAPI: %s\n", url)

	var response modelevolution.EvolutionChain
	err := c.enqueueAndFetch(ctx, url, &response)
	return response, err
}
