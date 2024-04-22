package main

import (
	"Go-Shopping-backend/api"
	"Go-Shopping-backend/initializers"
	//"Go-Shopping-backend/middleware"
	"Go-Shopping-backend/utils"
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	// only load the .env file when running locally
	// check for a RAILWAY_ENVIRONMENT, if not found, code is running locally
	if _, exists := os.LookupEnv("RAILWAY_ENVIRONMENT"); !exists {
		initializers.LoadEnvVariables()
	}
	initializers.ConnectToDB()
	initializers.ConnectToRedis()
	initializers.LoadProductsToRedis()
}

func main() {
	// Create a new Gin router
	router := gin.Default()

	router.Use(utils.Cors())

	//router.Use(middleware.RateLimitMiddleware())

	baseRoute := "api/v1"

	//default route
	api.Default(router.Group(baseRoute))

	// api route for Signup and Login
	// api.AccountRoutes(router.Group(baseRoute + "/account"))
	// have not using this api because he have implemented supabase for authentication

	// api route for Products like add , get, update , delete
	api.ProductRoutes(router.Group(baseRoute + "/products"))

	// api route for Cart
	api.CartRouter(router.Group(baseRoute + "/cart"))

	rows, err := initializers.DB.Query(context.Background(), "SELECT id , username , email FROM users")
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}

	for rows.Next() {
		var id uuid.UUID
		var username any
		var email string
		// Add more variables as per your user schema
		if err := rows.Scan(&id, &username, &email); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		fmt.Printf("User ID: %d, Username: %s, Email: %s\n", id, username, email)
		// Print more user data as needed
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating over rows:", err)
		return
	}

	// Run the server on port 3000
	router.Run()
}
