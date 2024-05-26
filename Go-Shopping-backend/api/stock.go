package api

import (
	"Go-Shopping-backend/controller"

	"github.com/gin-gonic/gin"
)

func StockRoute(router *gin.RouterGroup) {
	router.POST("/update/:id", controller.UpdateStocks)
	router.GET("/:id", controller.GetStocks)
}
