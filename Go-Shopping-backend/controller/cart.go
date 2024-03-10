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

	rows, err := initializers.DB.Raw("SELECT products FROM carts WHERE user_id = ?", id).Rows()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to retrieve cart"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		// Scan the array directly into the slice of int32
		var productsArray pq.Int32Array
		if err := rows.Scan(&productsArray); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to scan cart data"})
			return
		}

		// Convert pq.Int32Array to []int32
		cart.Products = []int32(productsArray)
	}
	if err := rows.Err(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Error occurred while retrieving cart data"})
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
		value, err := initializers.RedisClient.Get("product:" + strconv.Itoa(int(ProductID))).Result()
		if err != nil {
			var product models.Product

			initializers.DB.First(&product, ProductID)
			products = append(products, product)
			continue
		}

		var product models.Product
		if err := json.Unmarshal([]byte(value), &product); err != nil {
			log.Printf("Error decoding product with key in redis")
			continue
		}

		products = append(products, product)
	}

	context.JSON(http.StatusOK, products)

}

func AddProductToCart(context *gin.Context) {

	var body struct {
		UserID    int   `json:"user_id"`
		ProductID int32 `json:"product_id"`
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
		cart.Products = []int32{body.ProductID}

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

		var newcart models.Cart

		// Use custom scan function to retrieve cart data
		rows, err := initializers.DB.Raw("SELECT products FROM carts WHERE user_id = ?", body.UserID).Rows()
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to retrieve cart"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			// Scan the array directly into the slice of int32
			var productsArray pq.Int32Array
			if err := rows.Scan(&productsArray); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to scan cart data"})
				return
			}

			// Convert pq.Int32Array to []int32
			newcart.Products = []int32(productsArray)
		}
		if err := rows.Err(); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Error occurred while retrieving cart data"})
			return
		}

		// Check if the product is already in the cart
		found := false
		for _, productID := range newcart.Products {
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
		newProducts := make([]int32, len(newcart.Products))
		copy(newProducts, newcart.Products)
		newProducts = append(newProducts, body.ProductID)
		newcart.Products = newProducts

		// Update the cart in the database
		if err := initializers.DB.Model(&newcart).Where("user_id = ?", body.UserID).Updates(map[string]interface{}{"products": pq.Array(newcart.Products), "updated_at": time.Now(), "user_id": body.UserID}).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update cart"})
			return
		}

		// Return success message
		context.JSON(http.StatusCreated, gin.H{"message": "Product added to cart"})

	}

}
