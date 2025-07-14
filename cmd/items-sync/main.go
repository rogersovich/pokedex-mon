package main

import (
	"context"
	"log"
	"os"
	"time"

	"pokedex/config"
	"pokedex/database"
	"pokedex/internal/items/repository"
	"pokedex/internal/items/service"
	"pokedex/internal/shared/pokeapi"
)

func main() {
	log.Println("Starting Items Sync Job...")

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	database.ConnectDB(cfg)
	defer database.DisconnectDB()

	// Initialize shared PokeAPI client
	pokeAPIClient := pokeapi.NewClient(cfg)
	defer pokeAPIClient.CloseClient()

	itemsRepo := repository.NewMongoItemsRepository()
	itemsService := service.NewItemsService(itemsRepo, pokeAPIClient)

	// Run the synchronization
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute) // Beri waktu yang cukup
	defer cancel()

	err := itemsService.SyncAllItems(ctx)
	if err != nil {
		log.Fatalf("Type Items sync job failed: %v", err)
		os.Exit(1) // Keluar dengan status error
	}

	log.Println("Type Items Sync Job completed successfully.")
	os.Exit(0) // Keluar dengan status sukses
}
