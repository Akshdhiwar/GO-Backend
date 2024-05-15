package api

import (
	"Go-Shopping-backend/controller"

	"github.com/gin-gonic/gin"
)

func AccountRoutes(router *gin.RouterGroup) {
	//Post signup for creating the new user
	// router.POST("/signup", controller.Signup)

	//Post Login for User Authentication
	// router.POST("/login", controller.Login)

	//GET Api for getting user data
	router.GET("/:id", controller.GetUserData)
}
