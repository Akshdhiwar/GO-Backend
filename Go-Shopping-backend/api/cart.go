package api

import (
	"Go-Shopping-backend/controller"
	"Go-Shopping-backend/middleware"

	"github.com/gin-gonic/gin"
)

func CartRouter(router *gin.RouterGroup) {

	// GET route to get the cart of specific user
	router.GET("/:id", middleware.Authenticate, controller.GetCart)

	// POST route to add the product to cart from user
	router.POST("/add", middleware.Authenticate, controller.AddProductToCart)

}
