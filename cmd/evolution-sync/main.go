package main

import (
	"context"
	"log"
	"os"
	"time"

	"pokedex/config"
	"pokedex/database"
	"pokedex/internal/evolution/repository"
	"pokedex/internal/evolution/service"
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

	evolutionRepo := repository.NewMongoEvolutionRepository()
	evolutionService := service.NewEvolutionService(evolutionRepo, pokeAPIClient)

	// Run the synchronization
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute) // Beri waktu yang cukup
	defer cancel()

	err := evolutionService.SyncAllEvolution(ctx)
	if err != nil {
		log.Fatalf("Evolution Chain sync job failed: %v", err)
		os.Exit(1) // Keluar dengan status error
	}

	log.Println("Evolution Chain Sync Job completed successfully.")
	os.Exit(0) // Keluar dengan status sukses
}
