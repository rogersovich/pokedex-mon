package handler

import (
	"context"
	"net/http"
	"pokedex/internal/generation/service"
	"pokedex/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type GenerationHandler struct {
	generationService service.GenerationService
}

func NewGenerationHandler(svc service.GenerationService) *GenerationHandler {
	return &GenerationHandler{
		generationService: svc,
	}
}

func (h *GenerationHandler) GetGenerationDetail(c *gin.Context) {
	identifier := c.Param("identifier")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	data, err := h.generationService.GetGenerationDetail(ctx, identifier)
	if err != nil {
		utils.BaseResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.BaseResponseDetailSuccess(c, "success get data", data, nil, nil)
}

func (h *GenerationHandler) GetGenerationList(c *gin.Context) {
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

	baseURLGeneration := utils.GetBaseURLDynamic(c, "generation")

	data, err := h.generationService.GetGenerationList(ctx, limit, offset, baseURLGeneration)
	if err != nil {
		utils.BaseResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.BaseResponsePaginateSuccess(c, "success get all data", data.Results, data.Count, data.Next, data.Previous)
}
