package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"pokedex/internal/pokemon/service"

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

	// --- Tambahkan baseUrl di sini ---
	// Mendapatkan skema (http/https), host, dan path dasar dari request
	scheme := "http"
	if c.Request.TLS != nil { // Cek apakah koneksi menggunakan HTTPS
		scheme = "https"
	}
	baseUrl := fmt.Sprintf("%s://%s/api/v1/pokemon", scheme, c.Request.Host)

	listResponse, err := h.pokemonService.GetPokemonList(ctx, limit, offset, baseUrl, searchQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve pokemon list"})
		return
	}

	c.JSON(http.StatusOK, listResponse)
}

func (h *PokemonHandler) GetPokemonDetail(c *gin.Context) {
	identifier := c.Param("identifier")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	pokemon, err := h.pokemonService.GetPokemon(ctx, identifier)
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
