package handler

import (
	"context"
	"fmt"
	"net/http"
	"pokedex/internal/ability/service"
	"time"

	"github.com/gin-gonic/gin"
)

type AbilityHandler struct {
	abilityService service.AbilityService
}

func NewAbilityHandler(svc service.AbilityService) *AbilityHandler {
	return &AbilityHandler{
		abilityService: svc,
	}
}

func (h *AbilityHandler) GetAbilityDetail(c *gin.Context) {
	identifier := c.Param("identifier")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	ability, err := h.abilityService.GetAbility(ctx, identifier)
	if err != nil {
		// More robust error checking for "not found"
		if err.Error() == fmt.Sprintf("ability not found: %s", identifier) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ability detail"})
		return
	}

	c.JSON(http.StatusOK, ability)
}
