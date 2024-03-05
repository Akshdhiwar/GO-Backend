package api

import (
	"Go-Shopping-backend/controller"
	"Go-Shopping-backend/middleware"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(router *gin.RouterGroup) {

	// Post route for adding products
	router.POST("/add", middleware.Authenticate, middleware.RoleBasedAuthorization, controller.AddProducts)

	//GET route for getting all products
	router.GET("/", controller.GetProducts)

	//Get route for getting single product through ID.
	router.GET("/:id", controller.GetSingleProduct)

	//Delete route for deleting product through ID.
	router.DELETE("/:id", middleware.Authenticate, middleware.RoleBasedAuthorization, controller.DeleteProduct)

	//Update route for updating products
	router.PUT("/:id", middleware.Authenticate, middleware.RoleBasedAuthorization, controller.UpdateProduct)
}
