package handler

import (
	"context"
	"net/http"
	"pokedex/internal/items/service"
	"pokedex/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ItemHandler struct {
	itemService service.ItemsService
}

func NewItemHandler(svc service.ItemsService) *ItemHandler {
	return &ItemHandler{
		itemService: svc,
	}
}

func (h *ItemHandler) GetItemDetail(c *gin.Context) {
	identifier := c.Param("identifier")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	data, err := h.itemService.GetItemDetail(ctx, identifier)
	if err != nil {
		utils.BaseResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.BaseResponseDetailSuccess(c, "success get all data", data, nil, nil)
}

func (h *ItemHandler) GetItemList(c *gin.Context) {
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

	baseURLItem := utils.GetBaseURLDynamic(c, "item")

	data, err := h.itemService.GetItemList(ctx, limit, offset, baseURLItem)
	if err != nil {
		utils.BaseResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.BaseResponsePaginateSuccess(c, "success get all data", data.Results, data.Count, data.Next, data.Previous)
}
