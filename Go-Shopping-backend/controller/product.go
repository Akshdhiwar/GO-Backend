package controller

import (
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetProducts(context *gin.Context) {
	var products []models.Product

	keys, err := initializers.RedisClient.Keys("product:*").Result()
	if err != nil {
		log.Printf("Failed to retrieve products from Redis")
	}

	for _, key := range keys {
		val, err := initializers.RedisClient.Get(key).Result()
		if err != nil {
			log.Printf("Error retrieving product with key %s: %v", key, err)
			continue
		}

		var product models.Product
		if err := json.Unmarshal([]byte(val), &product); err != nil {
			log.Printf("Error decoding product with key %s: %v", key, err)
			continue
		}
		products = append(products, product)
	}

	if len(products) > 0 {
		context.JSON(http.StatusOK, products)
		fmt.Println("got product from redis")
		return
	}

	// Key "products" doesn't exist in Redis
	// Fetch products from DB

	result := initializers.DB.Find(&products)
	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Error fetching products from DB: " + result.Error.Error(),
		})
		return
	}

	if len(products) == 0 {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "No products found",
		})
		return
	}

	// Set products in Redis
	err = SetProductsInRedis(products)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error setting products in Redis: " + err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, products)
}

func SetProductsInRedis(products []models.Product) error {
	// Set products in Redis

	fmt.Println("adding product to redis")

	for _, product := range products {
		// Convert product to JSON
		productJSON, err := json.Marshal(product)
		if err != nil {
			log.Printf("Error marshaling product: %v", err)
			continue
		}

		// Set product in Redis with key in the format "product:id"
		key := "product:" + strconv.Itoa(int(product.ID))
		err = initializers.RedisClient.Set(key, productJSON, 0).Err()
		if err != nil {
			log.Printf("Error setting product in Redis: %v", err)
		}
	}

	return nil
}

func AddProducts(context *gin.Context) {

	//getting body
	var body struct {
		Title       string
		Price       float64
		Description string
		Category    string
		Image       string
	}

	context.ShouldBind(&body)

	if body.Title == "" || body.Price == 0 || body.Image == "" || body.Category == "" || body.Description == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "All fields are required",
		})
		return
	}

	var product models.Product
	initializers.DB.First(&product, "title = ?", body.Title)

	if product.Title == body.Title {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Product already present please check Title. Title is matching with some other product",
		})
		return
	}

	newProduct := models.Product{Title: body.Title, Price: body.Price, Category: body.Category, Image: body.Image, Description: body.Description}
	result := initializers.DB.Create(&newProduct)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to save product to db",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Products successfully received",
	})

	var newlyAddedProduct models.Product
	initializers.DB.First(&newlyAddedProduct, "title = ?", body.Title)

	// Convert product to JSON
	productJSON, err := json.Marshal(newlyAddedProduct)
	if err != nil {
		log.Printf("Error marshaling product: %v", err)
	}

	// Set product in Redis with key in the format "product:id"
	key := "product:" + strconv.Itoa(int(newlyAddedProduct.ID))
	err = initializers.RedisClient.Set(key, productJSON, 0).Err()
	if err != nil {
		log.Printf("Error setting product in Redis: %v", err)
	}

}

func GetSingleProduct(context *gin.Context) {

	id := context.Param("id")

	var product models.Product

	key := "product:" + id

	exists, err := initializers.RedisClient.Exists(key).Result()

	if err != nil {
		panic(err)
	}

	if exists == 1 {
		val, err := initializers.RedisClient.Get(key).Result()
		if err != nil {
			log.Printf("Error retrieving product with key %s: %v", key, err)
		}

		if err := json.Unmarshal([]byte(val), &product); err != nil {
			log.Printf("Error decoding product with key %s: %v", key, err)
		}

		context.JSON(http.StatusOK, product)
		return
	}

	result := initializers.DB.First(&product, id)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Error querying product from database",
		})
		return
	}

	// If no product found with the given id
	if result.RowsAffected == 0 {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "No product found",
		})
		return
	}

	context.JSON(http.StatusOK, product)

}
