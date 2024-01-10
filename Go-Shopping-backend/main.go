package main

import (
	"Go-Shopping-backend/controller"
	"Go-Shopping-backend/initializers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Cors config
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"} // specify the origins you want to allow
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	router.Use(cors.New(config))

	// Define a route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, Gin!",
		})
	})

	// Post route for adding products
	router.POST("/add", controller.AddProducts)

	//GET route for getting all products
	router.GET("/products", controller.GetProducts)

	//Post signup for creating the new user
	router.POST("/signup", controller.Signup)

	// Run the server on port 3000
	router.Run()
}
