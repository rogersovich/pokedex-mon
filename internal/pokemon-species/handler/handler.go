package handler

import (
	"context"
	"net/http"
	"time"

	"pokedex/internal/pokemon-species/service"
	"pokedex/utils"

	"github.com/gin-gonic/gin"
)

// PokemonSpeciesHandler handles HTTP requests related to Pokemons.
type PokemonSpeciesHandler struct {
	pokemonSpeciesService service.PokemonSpeciesService
}

// NewPokemonSpeciesHandler creates a new instance of PokemonSpeciesHandler.
func NewPokemonSpeciesHandler(svc service.PokemonSpeciesService) *PokemonSpeciesHandler {
	return &PokemonSpeciesHandler{
		pokemonSpeciesService: svc,
	}
}

func (h *PokemonSpeciesHandler) GetPokemonSpeciesDetail(c *gin.Context) {
	identifier := c.Param("identifier")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	data, err := h.pokemonSpeciesService.GetPokemonSpecies(ctx, identifier)
	if err != nil {
		utils.BaseResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.BaseResponseDetailSuccess(c, "success get data", data, nil, nil)
}
