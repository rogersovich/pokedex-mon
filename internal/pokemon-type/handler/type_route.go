package handler

import (
	"github.com/gin-gonic/gin"
)

func RegisterPokemonTypeRoutes(rg *gin.RouterGroup, handler *PokemonTypeHandler) {
	pokemonTypeGroup := rg.Group("/type")
	{
		pokemonTypeGroup.GET("", handler.GetPokemonTypeList)
		pokemonTypeGroup.GET("/:identifier", handler.GetPokemonTypeDetail)
		pokemonTypeGroup.GET("/weakness/:pokemon-id", handler.GetWeaknessPokemonTypes)
	}
}
