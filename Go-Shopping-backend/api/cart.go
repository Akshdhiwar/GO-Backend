package api

import (
	"Go-Shopping-backend/controller"

	"github.com/gin-gonic/gin"
)

func CartRouter(router *gin.RouterGroup) {

	// GET route to get the cart of specific user
	router.GET("/:id", controller.GetCart)

	router.DELETE("/:id")
}
