package handler

import (
	"github.com/gin-gonic/gin"
)

func RegisterPokedexesRoutes(rg *gin.RouterGroup, handler *PokedexesHandler) {
	pokedexesGroup := rg.Group("/pokedex")
	{
		pokedexesGroup.GET("", handler.GetPokedexDetail)
		pokedexesGroup.GET("/:identifier", handler.GetPokedexDetail)
	}
}
