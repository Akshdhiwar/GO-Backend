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

	// POST route to add the quantity of respective product
	router.POST("/inc/:id", middleware.Authenticate, controller.AddQuantity)

	// POST route to remove the quantity of respective product
	router.POST("/dec/:id", middleware.Authenticate, controller.RemoveQuantity)

	// DELETE route to delete the product from cart from user
	router.DELETE("/delete/:id", middleware.Authenticate, controller.DeleteProductFromCart)

}
