package handler

import (
	"github.com/gin-gonic/gin"
)

func RegisterAbilityRoutes(rg *gin.RouterGroup, handler *AbilityHandler) {
	abilityGroup := rg.Group("/ability")
	{
		abilityGroup.GET("/:identifier", handler.GetAbilityDetail)
	}
}
