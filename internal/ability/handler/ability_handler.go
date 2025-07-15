package handler

import (
	"context"
	"net/http"
	"pokedex/internal/ability/service"
	"pokedex/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AbilityHandler struct {
	abilityService service.AbilityService
}

func NewAbilityHandler(svc service.AbilityService) *AbilityHandler {
	return &AbilityHandler{
		abilityService: svc,
	}
}

func (h *AbilityHandler) GetAbilityDetail(c *gin.Context) {
	identifier := c.Param("identifier")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	data, err := h.abilityService.GetAbility(ctx, identifier)
	if err != nil {
		utils.BaseResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.BaseResponseDetailSuccess(c, "success get data", data, nil, nil)
}

func (h *AbilityHandler) GetAbilityList(c *gin.Context) {
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

	baseURLAbility := utils.GetBaseURLDynamic(c, "ability")

	data, err := h.abilityService.GetAbilityList(ctx, limit, offset, baseURLAbility)
	if err != nil {
		utils.BaseResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.BaseResponsePaginateSuccess(c, "success get all data", data.Results, data.Count, data.Next, data.Previous)
}
