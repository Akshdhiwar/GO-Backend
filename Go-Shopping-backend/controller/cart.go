package controller

import (
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func AddProductToCart(ctx *gin.Context) {

	var body struct {
		UserID    int    `json:"user_id"`
		ProductID string `json:"product_id"`
	}

	err := ctx.Bind(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to Bind data",
		})
		return
	}

	if body.ProductID == "" || body.UserID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Product ID or User ID missing",
		})
		return
	}

	id, err := uuid.Parse(body.ProductID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid product ID",
		})
		return
	}

	var user models.User

	err = initializers.DB.QueryRow(context.Background(), "SELECT id FROM users WHERE id=$1", body.UserID).Scan(&user.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error querying user from database",
		})
		return
	}

	var product models.Product

	row := initializers.DB.QueryRow(context.Background(), "SELECT id FROM products WHERE id = $1", id)

	err = row.Scan(&product.ID)

	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Product not found",
			})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Error querying product from database",
			})
		}
		return
	}

	var cart models.Cart

	row = initializers.DB.QueryRow(context.Background(), "SELECT id FROM carts WHERE user_id=$1", body.UserID)

	err = row.Scan(&cart.ID)
	if err == pgx.ErrNoRows {
		fmt.Println("in IF statement")
		var cart models.Cart

		cart.UserID = body.UserID
		cart.Products = []string{id.String()}

		_, err := initializers.DB.Exec(context.Background(), "INSERT INTO carts (user_id, products, created_at, updated_at) VALUES ($1, $2 , $3 , $4)", body.UserID, cart.Products, time.Now(), time.Now())

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to create cart",
			})
			return
		}

		var cartID uuid.UUID
		err = initializers.DB.QueryRow(context.Background(), "SELECT id FROM carts WHERE user_id=$1", body.UserID).Scan(&cartID)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to get cart id of saved cart",
			})
			return
		}

		_, err = initializers.DB.Exec(context.Background(), "UPDATE users SET cart_id=$1 WHERE id=$2", cartID, body.UserID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to update user cart id",
			})
			return
		}

		ctx.JSON(
			http.StatusCreated,
			gin.H{
				"message": "Product added to cart",
			},
		)
	} else {
		fmt.Println("in ELSE statement")

		var newcart models.Cart

		err := initializers.DB.QueryRow(context.Background(), "SELECT products , user_id FROM carts WHERE user_id = $1", body.UserID).Scan(&newcart.Products, &newcart.UserID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to retrieve cart"})
			return
		}

		log.Println(newcart.Products)

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
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Product already in cart"})
			return
		}

		// Add the product to the cart
		newProducts := make([]string, len(newcart.Products))
		copy(newProducts, newcart.Products)
		newProducts = append(newProducts, body.ProductID)
		newcart.Products = newProducts

		// Update the cart in the database
		query := `
        UPDATE carts 
        SET products = $1, updated_at = $2 
        WHERE user_id = $3
    `
		// Execute the SQL query
		_, err = initializers.DB.Exec(context.Background(), query, newcart.Products, time.Now(), newcart.UserID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update cart"})
			return
		}

		// Return success message
		ctx.JSON(http.StatusCreated, gin.H{"message": "Product added to cart"})

	}
}

func GetCart(ctx *gin.Context) {
	id := ctx.Param("id")
	log.Println(id)

	var isCartIDPresent sql.NullString

	err := initializers.DB.QueryRow(context.Background(), "SELECT cart_id FROM users WHERE id = $1", id).Scan(&isCartIDPresent)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error querying user from database",
		})
		return
	}
	log.Println(isCartIDPresent)

	if !isCartIDPresent.Valid {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "No products in cart",
		})
		return
	}

	var cart models.Cart

	err = initializers.DB.QueryRow(context.Background(), "SELECT products FROM carts WHERE user_id = $1", id).Scan(&cart.Products)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to retrieve cart"})
		return
	}

	// if length of cart is 0 then return
	if len(cart.Products) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "No products in cart",
		})
		return
	}

	var products []models.Product

	for _, productID := range cart.Products {
		value, err := initializers.RedisClient.Get("product:" + productID).Result()
		if err != nil {
			var product models.Product

			err := initializers.DB.QueryRow(context.Background(), "SELECT * FROM products WHERE id=$1", productID).Scan(&product)
			if err != nil {
				log.Printf("Error querying product from database: %v", err)
				continue
			}

			products = append(products, product)
			continue
		}

		var product models.Product
		if err := json.Unmarshal([]byte(value), &product); err != nil {
			log.Printf("Error decoding product with key in redis: %v", err)
			continue
		}

		products = append(products, product)
	}
	ctx.JSON(http.StatusOK, products)
}

func DeleteProductFromCart(ctx *gin.Context) {

	var productID = ctx.Param("id")

	var body struct {
		UserID int `json:"user_id"`
	}

	err := ctx.ShouldBind(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to Bind data",
		})
	}

	if body.UserID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "User ID missing or invalid",
		})
	}

	_, err = uuid.Parse(productID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid product ID",
		})
	}

	var products []uuid.UUID

	err = initializers.DB.QueryRow(context.Background(), "SELECT products FROM carts WHERE user_id = $1", body.UserID).Scan(&products)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to retrieve cart",
		})
		return
	}

	for i, product := range products {
		if product.String() == productID {
			products = append(products[:i], products[i+1:]...)
			break
		}
	}

	_, err = initializers.DB.Exec(context.Background(), "UPDATE carts SET products = $1 WHERE user_id = $2", products, body.UserID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to update cart",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Product removed from cart",
	})

}
