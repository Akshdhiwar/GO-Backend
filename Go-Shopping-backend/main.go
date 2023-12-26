package main

import (
	// "fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new Gin router
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "https://go-backend-olive.vercel.app/"} // specify the origins you want to allow
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	router.Use(cors.New(config))

	// Define a route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, Gin!",
		})
	})

	router.GET("/products", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"product": "hello",
		})
	})

	// Run the server on port 3000
	router.Run(":3000")
}
