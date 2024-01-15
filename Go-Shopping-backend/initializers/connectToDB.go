package initializers

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
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
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Optional: Set connection pool settings if needed
	// DB.DB().SetMaxIdleConns(10)
	// DB.DB().SetMaxOpenConns(100)
}
