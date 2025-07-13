package handler

import (
	"context"
	"fmt"
	"net/http"
	"pokedex/internal/pokedexes/service"
	"pokedex/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PokedexesHandler struct {
	pokedexesService service.PokedexesService
}

func NewPokedexesHandler(svc service.PokedexesService) *PokedexesHandler {
	return &PokedexesHandler{
		pokedexesService: svc,
	}
}

func (h *PokedexesHandler) GetPokedexDetail(c *gin.Context) {
	identifier := c.Param("identifier")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	data, err := h.pokedexesService.GetPokedexDetail(ctx, identifier)
	if err != nil {
		utils.BaseResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.BaseResponseDetailSuccess(c, "success get all data", data, nil, nil)
}

func (h *PokedexesHandler) GetPokemonTypeList(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
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

	// --- Tambahkan baseUrl di sini ---
	// Mendapatkan skema (http/https), host, dan path dasar dari request
	scheme := "http"
	if c.Request.TLS != nil { // Cek apakah koneksi menggunakan HTTPS
		scheme = "https"
	}
	baseUrl := fmt.Sprintf("%s://%s/api/v1/pokedex", scheme, c.Request.Host)

	data, err := h.pokedexesService.GetPokedexList(ctx, limit, offset, baseUrl)
	if err != nil {
		utils.BaseResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.BaseResponsePaginateSuccess(c, "success get all data", data.Results, data.Count, data.Next, data.Previous)
}
