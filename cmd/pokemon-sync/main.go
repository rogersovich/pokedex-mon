package main

import (
	"context"
	"log"
	"os"
	"time"

	"pokedex/config"
	"pokedex/database"
	"pokedex/internal/pokemon/repository"
	"pokedex/internal/pokemon/service"
	"pokedex/internal/shared/pokeapi"
)

func main() {
	log.Println("Starting Pokemon Sync Job...")

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	database.ConnectDB(cfg)
	defer database.DisconnectDB()

	// Initialize shared PokeAPI client
	pokeAPIClient := pokeapi.NewClient(cfg)
	defer pokeAPIClient.CloseClient() // Penting untuk menutup klien setelah selesai

	// Initialize Pokemon Module Components needed for sync
	pokemonRepo := repository.NewMongoPokemonRepository()
	pokemonService := service.NewPokemonService(pokemonRepo, pokeAPIClient)

	// Run the synchronization
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute) // Beri waktu yang cukup
	defer cancel()

	err := pokemonService.SyncAllPokemons(ctx)
	if err != nil {
		log.Fatalf("Pokemon sync job failed: %v", err)
		os.Exit(1) // Keluar dengan status error
	}

	log.Println("Pokemon Sync Job completed successfully.")
	os.Exit(0) // Keluar dengan status sukses
}
