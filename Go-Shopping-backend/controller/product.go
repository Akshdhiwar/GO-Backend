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

func GetProducts(context *gin.Context) {
	// defining products array with product model
	var products []models.Product

	// getting the product from redis server
	keys, err := initializers.RedisClient.Keys("product:*").Result()
	if err != nil {
		log.Printf("Failed to retrieve products from Redis")
	}

	log.Println(keys)

	// looping over keys to get all the products and appending it to products array
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
		log.Println("got product from redis")
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

	// if no products were found than exit
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

	log.Println("adding product to redis")

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

	// binding the body with variable
	context.ShouldBind(&body)

	// condition check for all fields are required
	if body.Title == "" || body.Price == 0 || body.Image == "" || body.Category == "" || body.Description == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "All fields are required",
		})
		return
	}

	var product models.Product

	// checking if title with same name is present or not
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

	// getting id from url
	id := context.Param("id")

	var product models.Product

	key := "product:" + id

	// checking the product in redis server
	exists, err := initializers.RedisClient.Exists(key).Result()

	if err != nil {
		panic(err)
	}

	// if product exists in redis server sendeing the product
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

	// if product is not found in redis the product will be retrieved from result
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

func DeleteProduct(context *gin.Context) {

	// getting id from the url body
	id := context.Param("id")

	var product models.Product

	// checking if product is present in database
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

	key := "product:" + id
	err := initializers.RedisClient.Del(key).Err()
	if err != nil {
		log.Fatalf("Error deleting key %s: %v", key, err)
	}

	// deleting the product from DB
	result = initializers.DB.Delete(&product, id)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Error querying product from database",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Product Deleted successfully",
	})

}

func UpdateProduct(context *gin.Context) {
	// Getting body
	var body struct {
		Title       string  `json:"title"`
		Price       float64 `json:"price"`
		Description string  `json:"description"`
		Category    string  `json:"category"`
		Image       string  `json:"image"`
	}

	// Binding the body with variable
	if err := context.BindJSON(&body); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
		return
	}

	// Condition check for all fields being required
	if body.Title == "" || body.Price == 0 || body.Image == "" || body.Category == "" || body.Description == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "All fields are required"})
		return
	}

	id := context.Param("id")

	var product models.Product

	// Find the product by ID
	if err := initializers.DB.First(&product, id).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Product not found"})
		return
	}

	// Check if a product with the same title exists (excluding the current product being updated)
	var existingProduct models.Product
	if err := initializers.DB.Where("title = ? AND id != ?", body.Title, id).First(&existingProduct).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Product already present. Please check Title, it matches with another product"})
		return
	}

	// Update the product with data from the request body
	product.Title = body.Title
	product.Price = body.Price
	product.Description = body.Description
	product.Category = body.Category
	product.Image = body.Image

	// Save the updated product to the database
	if err := initializers.DB.Save(&product).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update product"})
		return
	}

	// get updated product and add to redis
	var updatedProduct models.Product
	initializers.DB.First(&updatedProduct, id)

	// Convert product to JSON
	productJSON, err := json.Marshal(updatedProduct)
	if err != nil {
		log.Printf("Error marshaling product: %v", err)
	}

	// Set product in Redis with key in the format "product:id"
	key := "product:" + strconv.Itoa(int(updatedProduct.ID))
	err = initializers.RedisClient.Set(key, productJSON, 0).Err()
	if err != nil {
		log.Printf("Error setting product in Redis: %v", err)
	}

	context.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}
