package initializers

import "Go-Shopping-backend/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
