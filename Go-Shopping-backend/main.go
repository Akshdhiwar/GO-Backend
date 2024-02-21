package main

import (
	"Go-Shopping-backend/api"
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/utils"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {

	// only load the .env file when running locally
	// check for a RAILWAY_ENVIRONMENT, if not found, code is running locally
	if _, exists := os.LookupEnv("RAILWAY_ENVIRONMENT"); !exists {
		initializers.LoadEnvVariables()
	}

	initializers.ConnectToDB()
}

func main() {
	// Create a new Gin router
	router := gin.Default()

	router.Use(utils.Cors())
	router.LoadHTMLGlob("views/*")

	//default route
	api.Default(router.Group("api/v1"))

	// api route for Signup and Login
	api.AccountRoutes(router.Group("api/v1/account"))

	// api route for Products like add , get, update , delete
	api.ProductRoutes(router.Group("api/v1/products"))

	// Run the server on port 3000
	router.Run()
}
