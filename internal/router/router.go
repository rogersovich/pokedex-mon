package router

import (
	ability_handler "pokedex/internal/ability/handler"
	evolution_handler "pokedex/internal/evolution/handler"
	pokemon_species_handler "pokedex/internal/pokemon-species/handler"
	pokemon_handler "pokedex/internal/pokemon/handler"

	"github.com/gin-gonic/gin"
)

// InitAPIRoutes initializes all API routes for the Gin engine.
func InitAPIRoutes(
	router *gin.Engine,
	pokemonHandler *pokemon_handler.PokemonHandler,
	abilityHandler *ability_handler.AbilityHandler,
	pokemonSpeciesHandler *pokemon_species_handler.PokemonSpeciesHandler,
	evolutionHandler *evolution_handler.EvolutionHandler,
) {
	v1 := router.Group("/api/v1")
	{
		pokemonGroup := v1.Group("/pokemon")
		{
			pokemonGroup.GET("/", pokemonHandler.GetPokemonList)
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
	}
}
