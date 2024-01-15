package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	dsn := os.Getenv("RAILS_DB")

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
