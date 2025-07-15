package handler

import (
	"github.com/gin-gonic/gin"
)

func RegisterItemRoutes(rg *gin.RouterGroup, handler *ItemHandler) {
	itemGroup := rg.Group("/item")
	{
		itemGroup.GET("", handler.GetItemList)
		itemGroup.GET("/:identifier", handler.GetItemDetail)
	}
}
