package main

import (
	"Go-Shopping-backend/controller"
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/middleware"
	"Go-Shopping-backend/models"
	"fmt"
	"log"
	"net/http"
	"os"

	// "github.com/gin-contrib/cors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {

	// only load the .env file when running locally
	// check for a RAILWAY_ENVIRONMENT, if not found, code is running locally
	if _, exists := os.LookupEnv("RAILWAY_ENVIRONMENT"); exists == false {
		initializers.LoadEnvVariables()
	}

	dbUser := os.Getenv("RAILS_DATABASE_USER")
	dbPassword := os.Getenv("RAILS_DATABASE_PASSWORD")
	dbName := os.Getenv("RAILS_DATABASE_NAME")
	dbHost := os.Getenv("RAILS_DATABASE_HOST")
	dbPort := os.Getenv("RAILS_DATABASE_PORT")

	// Construct DSN
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbUser, dbPassword, dbName, dbHost, dbPort)

	fmt.Println(dsn)
	if dsn == "" {
		log.Fatal("RAILS_DB environment variable is empty")
	}

	var err error

	initializers.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	initializers.DB.AutoMigrate(&models.User{})

}

func main() {
	// Create a new Gin router
	router := gin.Default()

	// // Cors config
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"} // specify the origins you want to allow
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	router.Use(cors.New(config))

	router.LoadHTMLGlob("views/*")
	// Define a route
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Post route for adding products
	router.POST("/add", controller.AddProducts)

	//GET route for getting all products
	router.GET("/products", controller.GetProducts)

	//Post signup for creating the new user
	router.POST("/signup", controller.Signup)

	//Post Login for User Authentication
	router.POST("/login", controller.Login)

	router.GET("/data", middleware.Authenticate, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"name": "Akash",
		})
	})

	// Run the server on port 3000
	router.Run()
}
