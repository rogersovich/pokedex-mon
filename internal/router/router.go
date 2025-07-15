package router

import (
	ability_handler "pokedex/internal/ability/handler"
	evolution_handler "pokedex/internal/evolution/handler"
	generation_handler "pokedex/internal/generation/handler"
	item_handler "pokedex/internal/items/handler"
	pokedexes_handler "pokedex/internal/pokedexes/handler"
	pokemon_species_handler "pokedex/internal/pokemon-species/handler"
	pokemon_type_handler "pokedex/internal/pokemon-type/handler"
	pokemon_handler "pokedex/internal/pokemon/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitAPIRoutes initializes all API routes for the Gin engine.
func InitAPIRoutes(
	router *gin.Engine,
	pokemonHandler *pokemon_handler.PokemonHandler,
	abilityHandler *ability_handler.AbilityHandler,
	pokemonSpeciesHandler *pokemon_species_handler.PokemonSpeciesHandler,
	evolutionHandler *evolution_handler.EvolutionHandler,
	pokemonTypeHandler *pokemon_type_handler.PokemonTypeHandler,
	pokedexesHandler *pokedexes_handler.PokedexesHandler,
	itemHandler *item_handler.ItemHandler,
	generationHandler *generation_handler.GenerationHandler,
) {

	// Configure CORS options
	corsConfig := cors.DefaultConfig()
	// corsConfig.AllowAllOrigins = true
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",
		"http://localhost:3001",
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true

	// Apply CORS middleware
	router.Use(cors.New(corsConfig))

	v1 := router.Group("/api/v1")
	{
		// Register routes for each module
		pokemon_handler.RegisterPokemonRoutes(v1, pokemonHandler)
		ability_handler.RegisterAbilityRoutes(v1, abilityHandler)
		pokemon_species_handler.RegisterPokemonSpeciesRoutes(v1, pokemonSpeciesHandler)
		evolution_handler.RegisterEvolutionRoutes(v1, evolutionHandler)
		pokemon_type_handler.RegisterPokemonTypeRoutes(v1, pokemonTypeHandler)
		pokedexes_handler.RegisterPokedexesRoutes(v1, pokedexesHandler)
		item_handler.RegisterItemRoutes(v1, itemHandler)
		generation_handler.RegisterGenerationRoutes(v1, generationHandler)
	}
}
