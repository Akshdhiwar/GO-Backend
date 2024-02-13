package api

import (
	"Go-Shopping-backend/controller"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(router *gin.RouterGroup) {

	// Post route for adding products
	router.POST("/add", controller.AddProducts)

	//GET route for getting all products
	router.GET("/", controller.GetProducts)
}
