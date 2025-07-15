package handler

import (
	"context"
	"net/http"
	"pokedex/internal/evolution/service"
	"pokedex/utils"
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

	data, err := h.evolutionService.GetEvolution(ctx, identifier)
	if err != nil {
		utils.BaseResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.BaseResponseDetailSuccess(c, "success get data", data, nil, nil)

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

	data, err := h.evolutionService.GetEvolutionPokemonType(ctx, pokemon_id)
	if err != nil {
		utils.BaseResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.BaseResponseDetailSuccess(c, "success get data", data, nil, nil)
}
