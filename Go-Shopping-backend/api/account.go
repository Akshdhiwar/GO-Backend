package api

import (
	"Go-Shopping-backend/controller"
	"Go-Shopping-backend/middleware"

	"github.com/gin-gonic/gin"
)

func AccountRoutes(router *gin.RouterGroup) {
	//Post signup for creating the new user
	// router.POST("/signup", controller.Signup)

	//Post Login for User Authentication
	// router.POST("/login", controller.Login)

	//GET Api for getting user data
	router.GET("/:id", controller.GetUserRole)

	router.GET("/data/:id", middleware.Authenticate, controller.GetUserData)

	router.POST("/:id", middleware.Authenticate, controller.UpdateUserData)
}
