package api

import (
	"Go-Shopping-backend/controller"
	"Go-Shopping-backend/middleware"

	"github.com/gin-gonic/gin"
)

func StockRoute(router *gin.RouterGroup) {
	router.POST("/update/:id", middleware.Authenticate, middleware.RoleBasedAuthorization, controller.UpdateStocks)
	router.GET("/:id", middleware.Authenticate, middleware.RoleBasedAuthorization, controller.GetStocks)
}
