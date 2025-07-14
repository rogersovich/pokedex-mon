package handler

import (
	"context"
	"net/http"
	"pokedex/internal/pokemon-type/service"
	"pokedex/utils"
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

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	baseURLPokeType := utils.GetBaseURLDynamic(c, "type")

	data, err := h.pokemonTypeService.GetPokemonTypeList(ctx, limit, offset, baseURLPokeType)
	if err != nil {
		utils.BaseResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.BaseResponsePaginateSuccess(c, "success get all data", data.Results, data.Count, data.Next, data.Previous)
}

func (h *PokemonTypeHandler) GetPokemonTypeDetail(c *gin.Context) {
	identifier := c.Param("identifier")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	data, err := h.pokemonTypeService.GetPokemonType(ctx, identifier)
	if err != nil {
		utils.BaseResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.BaseResponseDetailSuccess(c, "success get data", data, nil, nil)
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

	data, err := h.pokemonTypeService.GetWeaknessPokemonTypes(ctx, pokemonIDInt, pokemonTypes)
	if err != nil {
		utils.BaseResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.BaseResponseDetailSuccess(c, "success get data", data, nil, nil)
}
