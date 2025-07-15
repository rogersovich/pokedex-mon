package handler

import (
	"github.com/gin-gonic/gin"
)

func RegisterGenerationRoutes(rg *gin.RouterGroup, handler *GenerationHandler) {
	generationGroup := rg.Group("/generation")
	{
		generationGroup.GET("", handler.GetGenerationList)
		generationGroup.GET("/:identifier", handler.GetGenerationDetail)
	}
}
