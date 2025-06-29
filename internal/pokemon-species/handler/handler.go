package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"pokedex/internal/pokemon-species/service"

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

	pokemon, err := h.pokemonSpeciesService.GetPokemonSpecies(ctx, identifier)
	if err != nil {
		// More robust error checking for "not found"
		if err.Error() == fmt.Sprintf("pokemon not found: %s", identifier) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve pokemon detail"})
		return
	}

	c.JSON(http.StatusOK, pokemon)
}
