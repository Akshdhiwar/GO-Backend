package initializers

import (
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadEnvVariables() {
	godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
}
