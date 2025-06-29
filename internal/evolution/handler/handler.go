package handler

import (
	"context"
	"fmt"
	"net/http"
	"pokedex/internal/evolution/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type EvolutionHandler struct {
	evolutionService service.EvolutionService
}

func NewEvolutionHandler(svc service.EvolutionService) *EvolutionHandler {
	return &EvolutionHandler{
		evolutionService: svc,
	}
}

func (h *EvolutionHandler) GetEvolutionDetail(c *gin.Context) {
	identifier := c.Param("identifier")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	evolution, err := h.evolutionService.GetEvolution(ctx, identifier)
	if err != nil {
		// More robust error checking for "not found"
		if err.Error() == fmt.Sprintf("evolution not found: %s", identifier) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve evolution detail"})
		return
	}

	c.JSON(http.StatusOK, evolution)
}

func (h *EvolutionHandler) GetEvolutionPokemonType(c *gin.Context) {
	idStr := c.Param("pokemon-id")
	pokemon_id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pokemon-id parameter"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	evolution, err := h.evolutionService.GetEvolutionPokemonType(ctx, pokemon_id)
	if err != nil {
		// More robust error checking for "not found"
		if err.Error() == fmt.Sprintf("pokemon not found: %d", pokemon_id) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}

	c.JSON(http.StatusOK, evolution)
}
