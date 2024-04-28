package initializers

import (
	"Go-Shopping-backend/database"
	"Go-Shopping-backend/models"
	"context"
	"encoding/json"
	"log"
)

func LoadProductsToRedis() {
	var products []models.Product

	rows, err := DB.Query(context.Background(), database.SelectAllProducts)
	if err != nil {
		log.Fatalf("Error fetching products from database: %v", err)
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
	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating product rows: %v", err)
		return
	}

	// Delete existing products in Redis
	redisProducts, err := RedisClient.Keys("product:*").Result()
	if err != nil {
		log.Fatalf("Error retrieving products from Redis: %v", err)
		return
	}
	for _, key := range redisProducts {
		err := RedisClient.Del(key).Err()
		if err != nil {
			log.Fatalf("Error deleting key %s: %v", key, err)
			return
		}
	}

	// Store products in Redis
	for _, product := range products {
		// Serialize product to JSON
		serializedProduct, err := json.Marshal(product)
		if err != nil {
			log.Printf("Error serializing product: %v", err)
			continue
		}
		key := "product:" + product.ID.String()
		err = RedisClient.Set(key, serializedProduct, 0).Err()
		if err != nil {
			log.Printf("Error storing product in Redis: %v", err)
			continue
		}
	}
}
