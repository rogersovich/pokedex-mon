package handler

import (
	"context"
	"net/http"
	"pokedex/internal/ability/service"
	"pokedex/utils"
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
