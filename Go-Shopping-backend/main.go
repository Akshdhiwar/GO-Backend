package main

import (
	// "fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Define a route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, Gin!",
		})
	})

	router.GET("/products", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"product": "all",
		})
	})

	// Run the server on port 3000
	router.Run(":3000")
}
