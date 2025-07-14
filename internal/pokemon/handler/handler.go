package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"pokedex/internal/pokemon/service"
	"pokedex/utils"

	"github.com/gin-gonic/gin"
)

// PokemonHandler handles HTTP requests related to Pokemons.
type PokemonHandler struct {
	pokemonService service.PokemonService
}

// NewPokemonHandler creates a new instance of PokemonHandler.
func NewPokemonHandler(svc service.PokemonService) *PokemonHandler {
	return &PokemonHandler{
		pokemonService: svc,
	}
}

func (h *PokemonHandler) GetPokemonList(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	searchQuery := c.DefaultQuery("q", "")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	baseURLPokemon := utils.GetBaseURLDynamic(c, "pokemon")

	data, err := h.pokemonService.GetPokemonList(ctx, limit, offset, baseURLPokemon, searchQuery)
	if err != nil {
		utils.BaseResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.BaseResponsePaginateSuccess(c, "success get all data", data.Results, data.Count, data.Next, data.Previous)
}

func (h *PokemonHandler) GetPokemonDetail(c *gin.Context) {
	identifier := c.Param("identifier")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	baseURL := utils.GetBaseURL(c)

	data, next, prev, err := h.pokemonService.GetPokemon(ctx, identifier, baseURL)
	if err != nil {
		utils.BaseResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.BaseResponseDetailSuccess(c, "success get data", data, next, prev)
}
