package controller

import (
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCart(context *gin.Context) {
	id := context.Param("id")

	var cart models.Cart

	result := initializers.DB.First(&cart, "user_id = ?", id)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Error querying cart from database",
		})
		return
	}

	if result.RowsAffected == 0 {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "No cart found",
		})
		return
	}

	// if length of cart is 0 then return
	if len(cart.Products) == 0 {
		context.JSON(http.StatusOK, gin.H{
			"message": "No products in cart",
		})
		return
	}

	var products []models.Product

	for _, productID := range cart.Products {
		value, err := initializers.RedisClient.Get("product:" + strconv.Itoa(productID)).Result()
		if err != nil {
			log.Printf("Error retrieving product with key %s: %v", productID, err)
			continue
		}

		var product models.Product
		if err := json.Unmarshal([]byte(value), &product); err != nil {
			log.Printf("Error decoding product with key %s: %v", productID, err)
			continue
		}

		products = append(products, product)
	}

	context.JSON(http.StatusOK, products)

}
