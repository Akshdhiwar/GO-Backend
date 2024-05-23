package api

import (
	"Go-Shopping-backend/controller"
	"Go-Shopping-backend/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoute(router *gin.RouterGroup) {
	router.GET("/product/:id", middleware.Authenticate, middleware.RoleBasedAuthorization, controller.GetSingleProductAdmin)
}
