package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"pokedex/internal/pokemon-type/service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PokemonTypeHandler struct {
	pokemonTypeService service.PokemonTypeService
}

func NewPokemonTypeHandler(svc service.PokemonTypeService) *PokemonTypeHandler {
	return &PokemonTypeHandler{
		pokemonTypeService: svc,
	}
}

func (h *PokemonTypeHandler) GetPokemonTypeList(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "30")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	pokemonTypes := c.QueryArray("types")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	log.Printf("Pokemon Types: %v", pokemonTypes)

	// --- Tambahkan baseUrl di sini ---
	// Mendapatkan skema (http/https), host, dan path dasar dari request
	scheme := "http"
	if c.Request.TLS != nil { // Cek apakah koneksi menggunakan HTTPS
		scheme = "https"
	}
	baseUrl := fmt.Sprintf("%s://%s/api/v1/type", scheme, c.Request.Host)

	listResponse, err := h.pokemonTypeService.GetPokemonTypeList(ctx, limit, offset, baseUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve type list"})
		return
	}

	c.JSON(http.StatusOK, listResponse)
}

func (h *PokemonTypeHandler) GetPokemonTypeDetail(c *gin.Context) {
	identifier := c.Param("identifier")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	res, err := h.pokemonTypeService.GetPokemonType(ctx, identifier)
	if err != nil {
		// More robust error checking for "not found"
		if err.Error() == fmt.Sprintf("data not found: %s", identifier) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *PokemonTypeHandler) GetWeaknessPokemonTypes(c *gin.Context) {
	pokemonID := c.Param("pokemon-id")
	pokemonTypesStr := c.Query("types")

	pokemonTypes := strings.Split(pokemonTypesStr, ",")

	if pokemonID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pokemon-id parameter"})
		return
	}

	pokemonIDInt, err := strconv.Atoi(pokemonID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Format pokemon-id to number"})
		return
	}

	if len(pokemonTypes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid types parameter"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	res, err := h.pokemonTypeService.GetWeaknessPokemonTypes(ctx, pokemonIDInt, pokemonTypes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
