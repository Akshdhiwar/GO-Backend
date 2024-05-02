package controller

import (
	"Go-Shopping-backend/database"
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/models"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/price"
	"github.com/stripe/stripe-go/v78/product"
)

func GetProducts(ctx *gin.Context) {

	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil || offset <= 0 {
		offset = 0
	}

	var originalOffset = offset

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil || limit <= 0 {
		limit = 16
	}

	if limit > 20 {
		limit = 20
	}

	offset = offset * limit

	// defining products array with product model
	var products []models.Product

	// getting the product from redis server
	keys, err := initializers.RedisClient.Keys("product:*").Result()
	if err != nil {
		log.Printf("Failed to retrieve products from Redis")
	}

	start := offset
	end := (limit * originalOffset) + limit
	if end > len(keys) {
		end = len(keys)
	}
	keys = keys[start:end]

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
		ctx.JSON(http.StatusOK, products)
		return
	}

	// Key "products" doesn't exist in Redis
	// Fetch products from DB
	rows, err := initializers.DB.Query(context.Background(), database.SelectAllProductsLimit, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error fetching products from DB: " + err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt, &product.Title, &product.Price, &product.Description, &product.Category, &product.Image, &product.Rating, &product.Count, &product.PriceID); err != nil {
			log.Printf("Error scanning product row: %v", err)
			continue
		}
		products = append(products, product)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error iterating product rows: " + err.Error()})
		return
	}

	// if no products were found than exit
	if len(products) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "No products found",
		})
		return
	}

	// Set products in Redis
	err = SetProductsInRedis(products)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error setting products in Redis: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func SetProductsInRedis(products []models.Product) error {
	// Set products in Redis

	for _, product := range products {
		// Convert product to JSON
		productJSON, err := json.Marshal(product)
		if err != nil {
			log.Printf("Error marshaling product: %v", err)
			continue
		}

		// Set product in Redis with key in the format "product:id"
		key := "product:" + product.ID.String()
		err = initializers.RedisClient.Set(key, productJSON, 0).Err()
		if err != nil {
			log.Printf("Error setting product in Redis: %v", err)
		}
	}

	return nil
}

func AddProducts(ctx *gin.Context) {

	//getting body
	var body struct {
		Title       string  `json:"title"`
		Price       float64 `json:"price"`
		Description string  `json:"description"`
		Category    string  `json:"category"`
		Image       string  `json:"image"`
	}

	// binding the body with variable
	ctx.ShouldBindJSON(&body)

	// condition check for all fields are required
	if body.Title == "" || body.Price == 0 || body.Image == "" || body.Category == "" || body.Description == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "All fields are required",
		})
		return
	}

	var nProduct models.Product

	// checking if title with same name is present or not
	var err error
	err = initializers.DB.QueryRow(context.Background(), database.SelectProductDetailsFromTitle, body.Title).Scan(&nProduct.ID, &nProduct.CreatedAt, &nProduct.UpdatedAt, &nProduct.Title, nProduct.Price, &nProduct.Description, &nProduct.Category, &nProduct.Image, &nProduct.Rating, &nProduct.Count, &nProduct.PriceID)
	if err == pgx.ErrNoRows {
		log.Printf("No product with title '%s' found", body.Title)
	}

	if nProduct.Title == body.Title {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Product already present please check Title. Title is matching with some other product",
		})
		return
	}

	images := []*string{&body.Image}

	params := &stripe.ProductParams{
		Name:        stripe.String(body.Title),
		Description: stripe.String(body.Description),
		Images:      images,
	}

	stripeProduct, err := product.New(params)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	priceParams := &stripe.PriceParams{
		Product:           stripe.String(stripeProduct.ID),
		Currency:          stripe.String(string(stripe.CurrencyUSD)),
		UnitAmountDecimal: stripe.Float64(body.Price*100),
	}

	price, err := price.New(priceParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	newProduct := models.Product{ID: uuid.New(), CreatedAt: time.Now(), Title: body.Title, Price: body.Price, Category: body.Category, Image: body.Image, Description: body.Description, PriceID: price.ID}

	_, err = initializers.DB.Exec(context.Background(), database.SaveNewProduct, newProduct.ID, newProduct.CreatedAt, newProduct.Title, newProduct.Price, newProduct.Category, newProduct.Image, newProduct.Description, 0, 0, newProduct.PriceID)
	if err != nil {
		log.Fatalf("Error creating new product: %v", err)
	}

	// var newlyAddedProduct models.Product
	// err = initializers.DB.QueryRow(context.Background(), database.SelectProductDetailsFromTitle, body.Title).Scan(&newlyAddedProduct.ID, &newlyAddedProduct.CreatedAt, &newlyAddedProduct.UpdatedAt, &newlyAddedProduct.Title, &newlyAddedProduct.Price, &newlyAddedProduct.Description, &newlyAddedProduct.Category, &newlyAddedProduct.Image, &newlyAddedProduct.Rating, &newlyAddedProduct.Count, &newlyAddedProduct.PriceID)
	// if err != nil {
	// 	log.Fatalf("Error fetching product from DB: %v", err)
	// }

	// Convert product to JSON
	productJSON, err := json.Marshal(newProduct)
	if err != nil {
		log.Printf("Error marshaling product: %v", err)
	}

	// Set product in Redis with key in the format "product:id"
	key := "product:" + newProduct.ID.String()
	err = initializers.RedisClient.Set(key, productJSON, 0).Err()
	if err != nil {
		log.Printf("Error setting product in Redis: %v", err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Products successfully received",
	})

}

func GetSingleProduct(ctx *gin.Context) {

	// getting id from url
	id := ctx.Param("id")

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

		ctx.JSON(http.StatusOK, product)
		return
	}

	err = initializers.DB.QueryRow(context.Background(), database.SelectAllFromID, id).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Title, &product.Price, &product.Description, &product.Category, &product.Image, &product.Rating, &product.Count, &product.PriceID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No product found",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Error querying db",
		})
		return
	}

	ctx.JSON(http.StatusOK, product)

}

func DeleteProduct(ctx *gin.Context) {

	// getting id from the url body
	idStr := ctx.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid product id",
		})
		return
	}

	var product models.Product

	err = initializers.DB.QueryRow(context.Background(), database.SelectAllFromID, id).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Title, &product.Price, &product.Description, &product.Category, &product.Image, &product.Rating, &product.Count, &product.PriceID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No product found",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Error querying db",
		})
		return
	}

	key := "product:" + id.String()
	delete := initializers.RedisClient.Del(key).Err()
	if delete != nil {
		log.Fatalf("Error deleting key %s: %v", key, err)
	}

	// Delete the product from the database
	_, err = initializers.DB.Exec(context.Background(), database.DeleteProduct, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error deleting product from database",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Product Deleted successfully",
	})

}

func UpdateProduct(ctx *gin.Context) {
	// Getting body
	var body struct {
		Title       string  `json:"title"`
		Price       float64 `json:"price"`
		Description string  `json:"description"`
		Category    string  `json:"category"`
		Image       string  `json:"image"`
		Rating      float32
		Count       int
	}

	// Binding the body with variable
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
		return
	}

	// Condition check for all fields being required
	if body.Title == "" || body.Price == 0 || body.Image == "" || body.Category == "" || body.Description == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "All fields are required"})
		return
	}

	idStr := ctx.Param("id")

	// Parsing the UUID string from the URL parameter
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID"})
		return
	}

	var product models.Product

	// Find the product by ID
	err = initializers.DB.QueryRow(context.Background(), database.SelectAllFromID, id).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Title, &product.Price, &product.Description, &product.Category, &product.Image, &product.Rating, &product.Count, &product.PriceID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No product found",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Error querying db",
		})
		return
	}

	// Check if a product with the same title exists (excluding the current product being updated)
	var existingProduct models.Product
	err = initializers.DB.QueryRow(context.Background(), database.SelectIdFromProductsMismatch, body.Title, id).Scan(&existingProduct.ID)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Product already present. Please check Title, it matches with another product"})
		return
	} else if err != pgx.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error querying db"})
		return
	}

	// Update the product with data from the request body
	product.Title = body.Title
	product.Price = body.Price
	product.Description = body.Description
	product.Category = body.Category
	product.Image = body.Image
	product.Rating = body.Rating
	product.Count = body.Count

	// Save the updated product to the database
	_, err = initializers.DB.Exec(context.Background(), database.UpdateProduct, product.Title, product.Price, product.Description, product.Category, product.Image, product.Rating, product.Count, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update product"})
		return
	}

	// get updated product and add to redis
	var updatedProduct models.Product
	err = initializers.DB.QueryRow(context.Background(), database.SelectAllFromID, id).Scan(&updatedProduct.ID, &updatedProduct.CreatedAt, &updatedProduct.UpdatedAt, &updatedProduct.Title, &updatedProduct.Price, &updatedProduct.Description, &updatedProduct.Category, &updatedProduct.Image, &updatedProduct.Rating, &updatedProduct.Count, &updatedProduct.PriceID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No product found",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Error querying db",
		})
		return
	}

	// Convert product to JSON
	productJSON, err := json.Marshal(updatedProduct)
	if err != nil {
		log.Printf("Error marshaling product: %v", err)
	}

	// Set product in Redis with key in the format "product:id"
	key := "product:" + updatedProduct.ID.String()
	err = initializers.RedisClient.Set(key, productJSON, 0).Err()
	if err != nil {
		log.Printf("Error setting product in Redis: %v", err)
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}
