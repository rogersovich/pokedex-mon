package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"pokedex/config"
	"pokedex/database"
	"pokedex/internal/app"
	"pokedex/internal/router"
	"pokedex/internal/shared/pokeapi"
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

	// Initialize all application components (handlers, services, repos) via the app package
	application := app.NewApplication(cfg, pokeAPIClient)

	// Initialize Gin router
	routerEngine := gin.New()

	routerEngine.Use(gin.Logger())   // Add logger middleware
	routerEngine.Use(gin.Recovery()) // Add recovery middleware

	// Setup API routes for all modules, passing handlers from the application struct
	router.InitAPIRoutes(
		routerEngine,
		application.PokemonHandler,
		application.AbilityHandler,
		application.PokemonSpeciesHandler,
		application.EvolutionHandler,
		application.PokemonTypeHandler,
		application.PokedexesHandler,
		application.ItemHandler,
		application.GenerationHandler,
	)

	// Start Gin server in a goroutine so it doesn't block the main thread
	serverPort := ":" + cfg.Port
	go func() {
		if err := routerEngine.Run(serverPort); err != nil {
			log.Fatalf("Gin server failed to start: %v", err)
		}
	}()

	log.Printf("Server listening on %s\n", serverPort)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	log.Println("Server gracefully stopped.")
}
