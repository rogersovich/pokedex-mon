package handler

import (
	"github.com/gin-gonic/gin"
)

func RegisterEvolutionRoutes(rg *gin.RouterGroup, handler *EvolutionHandler) {
	evolutionGroup := rg.Group("/evolution")
	{
		evolutionGroup.GET("/:identifier", handler.GetEvolutionDetail)
		evolutionGroup.GET("/pokemon-type/:pokemon-id", handler.GetEvolutionPokemonType)
	}
}
