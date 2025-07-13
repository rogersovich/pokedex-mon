package main

import (
	"context"
	"log"
	"os"
	"time"

	"pokedex/config"
	"pokedex/database"
	"pokedex/internal/pokedexes/repository"
	"pokedex/internal/pokedexes/service"
	"pokedex/internal/shared/pokeapi"
)

func main() {
	log.Println("Starting Pokedexes Sync Job...")

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	database.ConnectDB(cfg)
	defer database.DisconnectDB()

	// Initialize shared PokeAPI client
	pokeAPIClient := pokeapi.NewClient(cfg)
	defer pokeAPIClient.CloseClient()

	typeRepo := repository.NewMongoPokedexesRepository()
	typeService := service.NewPokedexesService(typeRepo, pokeAPIClient)

	// Run the synchronization
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute) // Beri waktu yang cukup
	defer cancel()

	err := typeService.SyncAllPokedexes(ctx)
	if err != nil {
		log.Fatalf("Type Pokedexes sync job failed: %v", err)
		os.Exit(1) // Keluar dengan status error
	}

	log.Println("Type Pokedexes Sync Job completed successfully.")
	os.Exit(0) // Keluar dengan status sukses
}
