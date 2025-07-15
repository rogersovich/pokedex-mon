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
		pokemonGroup := v1.Group("/pokemon")
		{
			pokemonGroup.GET("", pokemonHandler.GetPokemonList)
			pokemonGroup.GET("/:identifier", pokemonHandler.GetPokemonDetail)
		}
		abilityGroup := v1.Group("/ability")
		{
			// abilityGroup.GET("/", abilityHandler.GetPokeGemonList)
			abilityGroup.GET("/:identifier", abilityHandler.GetAbilityDetail)
		}
		pokemonSpeciesGroup := v1.Group("/pokemon-species")
		{
			// pokemonSpeciesGroup.GET("/", pokemonHandler.GetPokemonList)
			pokemonSpeciesGroup.GET("/:identifier", pokemonSpeciesHandler.GetPokemonSpeciesDetail)
		}
		evolutionGroup := v1.Group("/evolution")
		{
			// evolutionGroup.GET("/", pokemonHandler.GetPokemonList)
			evolutionGroup.GET("/:identifier", evolutionHandler.GetEvolutionDetail)
			evolutionGroup.GET("/pokemon-type/:pokemon-id", evolutionHandler.GetEvolutionPokemonType)
		}
		pokemonTypeGroup := v1.Group("/type")
		{
			pokemonTypeGroup.GET("", pokemonTypeHandler.GetPokemonTypeList)
			pokemonTypeGroup.GET("/:identifier", pokemonTypeHandler.GetPokemonTypeDetail)
			pokemonTypeGroup.GET("/weakness/:pokemon-id", pokemonTypeHandler.GetWeaknessPokemonTypes)
		}
		pokedexesGroup := v1.Group("/pokedex")
		{
			pokedexesGroup.GET("", pokedexesHandler.GetPokemonTypeList)
			pokedexesGroup.GET("/:identifier", pokedexesHandler.GetPokedexDetail)
		}
		itemGroup := v1.Group("/item")
		{
			itemGroup.GET("", itemHandler.GetItemList)
			itemGroup.GET("/:identifier", itemHandler.GetItemDetail)
		}
		generationGroup := v1.Group("/generation")
		{
			generationGroup.GET("", generationHandler.GetGenerationList)
			generationGroup.GET("/:identifier", generationHandler.GetGenerationDetail)
		}
	}
}
