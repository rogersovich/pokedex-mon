package main

import (
	"context"
	"log"
	"os"
	"time"

	"pokedex/config"
	"pokedex/database"
	"pokedex/internal/pokemon-type/repository"
	"pokedex/internal/pokemon-type/service"
	"pokedex/internal/shared/pokeapi"
)

func main() {
	log.Println("Starting Pokemon Type Sync Job...")

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	database.ConnectDB(cfg)
	defer database.DisconnectDB()

	// Initialize shared PokeAPI client
	pokeAPIClient := pokeapi.NewClient(cfg)
	defer pokeAPIClient.CloseClient()

	typeRepo := repository.NewMongoPokemonTypeRepository()
	typeService := service.NewPokemonTypeService(typeRepo, pokeAPIClient)

	// Run the synchronization
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute) // Beri waktu yang cukup
	defer cancel()

	err := typeService.SyncAllPokemonType(ctx)
	if err != nil {
		log.Fatalf("Type Pokemon sync job failed: %v", err)
		os.Exit(1) // Keluar dengan status error
	}

	log.Println("Type Pokemon Sync Job completed successfully.")
	os.Exit(0) // Keluar dengan status sukses
}
