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

	// DELETE route to delete the product from cart from user
	router.DELETE("/delete/:id", middleware.Authenticate, controller.DeleteProductFromCart)

}
