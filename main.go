package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"pokedex/config"
	"pokedex/database"
	"pokedex/internal/router"
	"pokedex/internal/shared/pokeapi"

	ability_handler "pokedex/internal/ability/handler"
	ability_repo "pokedex/internal/ability/repository"
	ability_service "pokedex/internal/ability/service"
	evolution_handler "pokedex/internal/evolution/handler"
	evolution_repo "pokedex/internal/evolution/repository"
	evolution_service "pokedex/internal/evolution/service"
	pokemon_species_handler "pokedex/internal/pokemon-species/handler"
	pokemon_species_repo "pokedex/internal/pokemon-species/repository"
	pokemon_species_service "pokedex/internal/pokemon-species/service"
	pokemon_type_handler "pokedex/internal/pokemon-type/handler"
	pokemon_type_repo "pokedex/internal/pokemon-type/repository"
	pokemon_type_service "pokedex/internal/pokemon-type/service"
	pokemon_handler "pokedex/internal/pokemon/handler"
	pokemon_repo "pokedex/internal/pokemon/repository"
	pokemon_service "pokedex/internal/pokemon/service"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB (shared by all modules)
	database.ConnectDB(cfg)
	defer database.DisconnectDB()

	// Initialize shared PokeAPI client (handles rate limiting for all modules)
	pokeAPIClient := pokeapi.NewClient(cfg)
	defer pokeAPIClient.CloseClient()

	// --- Initialize Pokemon Module Components ---
	evolutionRepo := evolution_repo.NewMongoEvolutionRepository()
	evolutionService := evolution_service.NewEvolutionService(evolutionRepo, pokeAPIClient)
	evolutionHandler := evolution_handler.NewEvolutionHandler(evolutionService)

	pokemonRepo := pokemon_repo.NewMongoPokemonRepository()
	pokemonService := pokemon_service.NewPokemonService(pokemonRepo, pokeAPIClient, evolutionService)
	pokemonHandler := pokemon_handler.NewPokemonHandler(pokemonService)

	abilityRepo := ability_repo.NewMongoAbilityRepository()
	abilityService := ability_service.NewAbilityService(abilityRepo, pokeAPIClient)
	abilityHandler := ability_handler.NewAbilityHandler(abilityService)

	pokemonSpeciesRepo := pokemon_species_repo.NewMongoPokemonSpeciesRepository()
	pokemonSpeciesService := pokemon_species_service.NewPokemonSpeciesService(pokemonSpeciesRepo, pokeAPIClient)
	pokemonSpeciesHandler := pokemon_species_handler.NewPokemonSpeciesHandler(pokemonSpeciesService)

	pokemonTypeRepo := pokemon_type_repo.NewMongoPokemonTypeRepository()
	pokemonTypeService := pokemon_type_service.NewPokemonTypeService(pokemonTypeRepo, pokeAPIClient)
	pokemonTypeHandler := pokemon_type_handler.NewPokemonTypeHandler(pokemonTypeService)

	// --- End Pokemon Module Components ---

	// Initialize Gin router
	routerEngine := gin.New() // Menggunakan Gin barebones

	routerEngine.Use(gin.Logger())   // Tambahkan logger
	routerEngine.Use(gin.Recovery()) // Tambahkan recovery

	// Setup API routes for all modules
	router.InitAPIRoutes(routerEngine, pokemonHandler, abilityHandler, pokemonSpeciesHandler, evolutionHandler, pokemonTypeHandler)

	// Start Gin server
	serverPort := ":" + cfg.Port

	if err := routerEngine.Run(serverPort); err != nil {
		log.Fatalf("Gin server failed to start: %v", err)
	}

	log.Printf("Server listening on %s\n", serverPort)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	log.Println("Server gracefully stopped.")
}
