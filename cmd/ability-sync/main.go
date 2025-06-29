package main

import (
	"context"
	"log"
	"os"
	"time"

	"pokedex/config"
	"pokedex/database"

	"pokedex/internal/ability/repository"
	"pokedex/internal/ability/service"
	"pokedex/internal/shared/pokeapi"
)

func main() {
	log.Println("Starting Ability Sync Job...")

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	database.ConnectDB(cfg)
	defer database.DisconnectDB()

	// Initialize shared PokeAPI client
	pokeAPIClient := pokeapi.NewClient(cfg)
	defer pokeAPIClient.CloseClient()

	// Initialize Ability Module
	abilityRepo := repository.NewMongoAbilityRepository()
	abilityService := service.NewAbilityService(abilityRepo, pokeAPIClient)

	// Run the synchronization
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	err := abilityService.SyncAllAbilities(ctx)
	if err != nil {
		log.Fatalf("Ability data sync failed: %v", err)
		os.Exit(1) // Keluar dengan status error
	}

	log.Println("Ability data sync completed successfully.")
	os.Exit(0) // Keluar dengan status sukses
}
