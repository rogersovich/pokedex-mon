package handler

import (
	"github.com/gin-gonic/gin"
)

func RegisterPokemonSpeciesRoutes(rg *gin.RouterGroup, handler *PokemonSpeciesHandler) {
	pokemonSpeciesGroup := rg.Group("/pokemon-species")
	{
		pokemonSpeciesGroup.GET("/:identifier", handler.GetPokemonSpeciesDetail)
	}
}
