package initializers

import (
	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
}
