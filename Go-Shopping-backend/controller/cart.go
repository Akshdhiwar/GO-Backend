package controller

import (
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/models"
	"encoding/json"
	"log"
	"time"

	"github.com/lib/pq"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCart(context *gin.Context) {
	id := context.Param("id")

	var cart models.Cart

	result := initializers.DB.Find(&cart, "user_id =?", id)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Error querying cart from database",
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

	for _, ProductID := range cart.Products {
		value, err := initializers.RedisClient.Get("product:" + strconv.Itoa(ProductID)).Result()
		if err != nil {
			log.Printf("Error retrieving product with key %s: %v", ProductID, err.Error())
			continue
		}

		var product models.Product
		if err := json.Unmarshal([]byte(value), &product); err != nil {
			log.Printf("Error decoding product with key %s: %v", ProductID, err.Error())
			continue
		}

		products = append(products, product)
	}

	context.JSON(http.StatusOK, products)

}

func AddProductToCart(context *gin.Context) {

	var body struct {
		UserID    int `json:"user_id"`
		ProductID int `json:"product_id"`
	}

	err := context.Bind(&body)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to Bind data",
		})
		return
	}

	if body.ProductID == 0 || body.UserID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Product ID or User ID missing",
		})
		return
	}

	var user models.User

	result := initializers.DB.First(&user, body.UserID)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Error querying user from database",
		})
		return
	}

	var cart models.Cart

	cartResult := initializers.DB.First(&cart, "user_id = ?", body.UserID)

	if cartResult.RowsAffected == 0 {
		var cart models.Cart

		cart.UserID = body.UserID
		cart.Products = []int{body.ProductID}

		result := initializers.DB.Exec("INSERT INTO carts (user_id, products, created_at, updated_at) VALUES ($1, $2 , $3 , $4)", body.UserID, pq.Array(cart.Products), time.Now(), time.Now())

		if result.Error != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to create cart",
			})
			return
		}

		context.JSON(
			http.StatusCreated,
			gin.H{
				"message": "Product added to cart",
			},
		)
	} else {

		var cart models.Cart

		// Retrieve cart from the database
		if err := initializers.DB.First(&cart, "user_id = ?", body.UserID).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Error querying cart from database"})
			return
		}

		log.Println(cart.Products)

		// Check if the product is already in the cart
		found := false
		for _, productID := range cart.Products {
			if productID == body.ProductID {
				found = true
				break
			}
		}

		// If product is already in the cart, return an error
		if found {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Product already in cart"})
			return
		}

		// Add the product to the cart
		cart.Products = append(cart.Products, body.ProductID)

		log.Println(cart.Products)

		// Update the cart in the database
		if err := initializers.DB.Model(&cart).Where("user_id = ?", body.UserID).Updates(map[string]interface{}{"products": cart.Products, "updated_at": time.Now()}).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update cart"})
			return
		}

		// Return success message
		context.JSON(http.StatusCreated, gin.H{"message": "Product added to cart"})

	}

}
