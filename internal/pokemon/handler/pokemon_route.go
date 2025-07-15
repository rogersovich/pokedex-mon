package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterPokemonRoutes registers all Pokémon related routes.
func RegisterPokemonRoutes(rg *gin.RouterGroup, handler *PokemonHandler) {
	pokemonGroup := rg.Group("/pokemon")
	{
		pokemonGroup.GET("", handler.GetPokemonList)
		pokemonGroup.GET("/:identifier", handler.GetPokemonDetail)
	}
}
