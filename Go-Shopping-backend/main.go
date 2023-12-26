package main

import (
	// "fmt"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Rating struct {
	Rate  float32 `json:"rate"`
	Count int     `json:"count"`
}

type Product struct {
	Id          int     `json:"id"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Image       string  `json:"image"`
	Rating      Rating  `json:"rating"`
}

func main() {
	// Create a new Gin router
	router := gin.Default()
	var products []Product
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

	router.POST("/add", func(c *gin.Context) {
		if err := c.BindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Products successfully received",
		})

		fmt.Printf("%+v\n", products)
	})

	router.GET("/products", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"product": products,
		})
	})

	// Run the server on port 3000
	router.Run(":3000")
}
