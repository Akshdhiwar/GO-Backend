package main

import (
	"Go-Shopping-backend/api"
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/middleware"
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
	initializers.ConnectToRedis()
	initializers.LoadProductsToRedis()
}

func main() {
	// Create a new Gin router
	router := gin.Default()

	router.Use(utils.Cors())

	router.Use(middleware.RateLimitMiddleware())

	baseRoute := "api/v1"

	//default route
	api.Default(router.Group(baseRoute))

	// api route for Signup and Login
	// api.AccountRoutes(router.Group(baseRoute + "/account"))
	// have not using this api because he have implemented supabase for authentication

	// api route for Products like add , get, update , delete
	api.ProductRoutes(router.Group(baseRoute + "/products"))

	// api route for Cart
	api.CartRouter(router.Group(baseRoute + "/cart"))

	// Run the server on port 3000
	router.Run()
}
