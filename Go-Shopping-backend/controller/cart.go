package controller

import (
	"Go-Shopping-backend/database"
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
		UserID    string `json:"user_id"`
		ProductID string `json:"product_id"`
	}

	err := ctx.Bind(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to Bind data",
			"type":    "error",
		})
		return
	}

	if body.ProductID == "" || body.UserID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Product ID or User ID missing",
			"type":    "error",
		})
		return
	}

	id, err := uuid.Parse(body.ProductID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid product ID",
			"type":    "error",
		})
		return
	}

	var user models.User

	userId, err := uuid.Parse(body.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error parsing UUID",
			"type":    "error",
		})
	}

	fmt.Println(userId)

	err = initializers.DB.QueryRow(context.Background(), database.SelectUserIdFromID, userId).Scan(&user.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error querying user from database",
			"type":    "error",
		})
		return
	}

	var product models.Product

	row := initializers.DB.QueryRow(context.Background(), database.SelectProductIdFromId, id)

	err = row.Scan(&product.ID)

	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Product not found",
				"type":    "success",
			})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Error querying product from database",
				"type":    "error",
			})
		}
		return
	}

	var cart models.Cart

	row = initializers.DB.QueryRow(context.Background(), database.SelectCartIdFromUserId, userId)

	err = row.Scan(&cart.ID)
	if err == pgx.ErrNoRows {
		fmt.Println("in IF statement")
		var cart models.Cart

		var cartProducts models.CartProduct
		cartProducts.ProductID = id
		cartProducts.Quantity = 1

		id, err := uuid.Parse(body.UserID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error parsing UUID",
				"type":    "error",
			})
		}

		cart.UserID = id
		cart.Products = []models.CartProduct{cartProducts}

		_, err = initializers.DB.Exec(context.Background(), database.SaveCart, userId, cart.Products, time.Now(), time.Now())

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to create cart",
				"type":    "error",
			})
			return
		}

		var cartID uuid.UUID
		err = initializers.DB.QueryRow(context.Background(), database.SelectCartIdFromUserId, userId).Scan(&cartID)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to get cart id of saved cart",
				"type":    "error",
			})
			return
		}

		_, err = initializers.DB.Exec(context.Background(), database.UpdateCartId, cartID, userId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to update user cart id",
				"type":    "error",
			})
			return
		}

		ctx.JSON(
			http.StatusCreated,
			gin.H{
				"message": "Product added to cart",
				"type":    "success",
			},
		)
	} else {
		fmt.Println("in ELSE statement")

		var newcart models.Cart

		err := initializers.DB.QueryRow(context.Background(), database.SelectCartDetailsFromUserId, userId).Scan(&newcart.Products, &newcart.UserID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to retrieve cart", "type": "error"})
			return
		}

		// Check if the product is already in the cart
		found := false
		for _, product := range newcart.Products {
			if product.ProductID.String() == body.ProductID {
				found = true
				break
			}
		}

		// If product is already in the cart, return an error
		if found {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Product already in cart",
				"type":    "error",
			})
			return
		}

		// Add the product to the cart
		newProducts := make([]models.CartProduct, len(newcart.Products))
		copy(newProducts, newcart.Products)

		var newProductJSON models.CartProduct

		productID, err := uuid.Parse(body.ProductID)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to parse UUID",
				"type":    "error",
			})
		}

		newProductJSON.ProductID = productID
		newProductJSON.Quantity = 1

		newProducts = append(newProducts, newProductJSON)
		newcart.Products = newProducts

		// Execute the SQL query
		_, err = initializers.DB.Exec(context.Background(), database.UpdateCart, newcart.Products, time.Now(), newcart.UserID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to update cart",
				"type":    "error",
			})
			return
		}

		// Return success message
		ctx.JSON(http.StatusCreated, gin.H{"message": "Product added to cart", "type": "success"})

	}
}

func GetCart(ctx *gin.Context) {
	id := ctx.Param("id")
	var isCartIDPresent sql.NullString

	err := initializers.DB.QueryRow(context.Background(), database.SelectCartIdFromId, id).Scan(&isCartIDPresent)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error querying user from database",
			"type":    "error",
		})
		return
	}

	if !isCartIDPresent.Valid {
		products := make([]int, 0)
		ctx.JSON(http.StatusOK, products)
		return
	}

	var cart models.Cart

	err = initializers.DB.QueryRow(context.Background(), database.SelectProductsFromUserID, id).Scan(&cart.Products)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to retrieve cart", "type": "error"})
		return
	}

	// if length of cart is 0 then return
	if len(cart.Products) == 0 {
		products := make([]int, 0)
		ctx.JSON(http.StatusOK, products)
		return
	}

	var products []any

	for _, cartProduct := range cart.Products {
		id := cartProduct.ProductID
		value, err := initializers.RedisClient.Get("product:" + cartProduct.ProductID.String()).Result()
		if err != nil {
			var product models.Product

			err := initializers.DB.QueryRow(context.Background(), database.SelectAllFromID, id).Scan(&product)
			if err != nil {
				log.Printf("Error querying product from database: %v", err)
				continue
			}

			var productJson struct {
				Product  models.Product
				Quantity int
			}

			productJson.Product = product
			productJson.Quantity = cartProduct.Quantity

			products = append(products, productJson)
			continue
		}

		var product models.Product
		if err := json.Unmarshal([]byte(value), &product); err != nil {
			log.Printf("Error decoding product with key in redis: %v", err)
			continue
		}

		var productJson struct {
			Product  models.Product
			Quantity int
		}

		productJson.Product = product
		productJson.Quantity = cartProduct.Quantity

		products = append(products, productJson)
	}

	ctx.JSON(http.StatusOK, products)

}

func DeleteProductFromCart(ctx *gin.Context) {

	var productID = ctx.Param("id")

	var body struct {
		UserID string `json:"user_id"`
	}

	err := ctx.ShouldBind(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to Bind data",
			"type":    "error",
		})
		return
	}

	if body.UserID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "User ID missing or invalid",
			"type":    "error",
		})
		return
	}

	id, err := uuid.Parse(body.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid product ID",
			"type":    "error",
		})
		return
	}

	_, err = uuid.Parse(productID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid product ID",
			"type":    "error",
		})
		return
	}

	var products []models.CartProduct

	err = initializers.DB.QueryRow(context.Background(), database.SelectProductsFromUserID, id).Scan(&products)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to retrieve cart",
			"type":    "error",
		})
		return
	}

	for i, product := range products {
		if (product.ProductID).String() == productID {
			products = append(products[:i], products[i+1:]...)
			break
		}
	}

	_, err = initializers.DB.Exec(context.Background(), database.UpdateCartProductWhereUserId, products, body.UserID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to update cart",
			"type":    "error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Product removed from cart",
		"type":    "success",
	})

}

func AddQuantity(ctx *gin.Context) {

	var body struct {
		UserID string
	}

	err := ctx.ShouldBind(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error Binding data",
			"type":    "error",
		})
		return
	}

	userId, err := uuid.Parse(body.UserID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error pearsing uuid",
		})
	}

	if userId == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Please provide UserId",
			"type":    "error",
		})
		return
	}

	productID := ctx.Param("id")

	if productID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Please provide product id",
			"type":    "error",
		})
		return
	}

	var cartProducts []models.CartProduct

	err = initializers.DB.QueryRow(context.Background(), database.SelectProductsFromUserID, body.UserID).Scan(&cartProducts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error Quering calls to DB",
			"type":    "error",
		})
		return
	}

	pID, err := uuid.Parse(productID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "unable to parse product id",
			"type":    "error",
		})
	}

	var updatedCart []models.CartProduct

	for _, product := range cartProducts {
		if product.ProductID == pID {
			product.Quantity++
			updatedCart = append(updatedCart, product)
			continue
		}

		updatedCart = append(updatedCart, product)
	}

	_, err = initializers.DB.Exec(context.Background(), database.UpdateCartProductWhereUserId, updatedCart, body.UserID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to update cart",
			"type":    "error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Added Quantity for Product",
		"type":    "success",
	})

}

func RemoveQuantity(ctx *gin.Context) {

	var body struct {
		UserID string
	}

	err := ctx.ShouldBind(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error Binding data",
			"type":    "error",
		})
		return
	}

	userId, err := uuid.Parse(body.UserID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error pearsing uuid",
		})
	}

	if userId == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Please provide UserId",
			"type":    "error",
		})
		return
	}

	productID := ctx.Param("id")

	if productID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Please provide product id",
			"type":    "error",
		})
		return
	}

	var cartProducts []models.CartProduct

	err = initializers.DB.QueryRow(context.Background(), database.SelectProductsFromUserID, body.UserID).Scan(&cartProducts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error Quering calls to DB",
			"type":    "error",
		})
		return
	}

	pID, err := uuid.Parse(productID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "unable to parse product id",
			"type":    "error",
		})
	}

	var updatedCart []models.CartProduct

	for _, product := range cartProducts {
		if product.ProductID == pID {
			product.Quantity--
			updatedCart = append(updatedCart, product)
			continue
		}

		updatedCart = append(updatedCart, product)
	}

	_, err = initializers.DB.Exec(context.Background(), database.UpdateCartProductWhereUserId, updatedCart, body.UserID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to update cart",
			"type":    "error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Removed Quantity for Product",
		"type":    "success",
	})

}
