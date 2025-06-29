package main

import (
	"context"
	"log"
	"os"
	"time"

	"pokedex/config"
	"pokedex/database"
	"pokedex/internal/pokemon-species/repository"
	"pokedex/internal/pokemon-species/service"
	"pokedex/internal/shared/pokeapi"
)

func main() {
	log.Println("Starting Pokemon Species Sync Job...")

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	database.ConnectDB(cfg)
	defer database.DisconnectDB()

	// Initialize shared PokeAPI client
	pokeAPIClient := pokeapi.NewClient(cfg)
	defer pokeAPIClient.CloseClient()

	pokemonSpeciesRepo := repository.NewMongoPokemonSpeciesRepository()
	pokemonSpeciesService := service.NewPokemonSpeciesService(pokemonSpeciesRepo, pokeAPIClient)

	// Run the synchronization
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute) // Beri waktu yang cukup
	defer cancel()

	err := pokemonSpeciesService.SyncAllPokemonSpecies(ctx)
	if err != nil {
		log.Fatalf("Pokemon species sync job failed: %v", err)
		os.Exit(1) // Keluar dengan status error
	}

	log.Println("Pokemon species Sync Job completed successfully.")
	os.Exit(0) // Keluar dengan status sukses
}
