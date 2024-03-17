package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
	Role     int       `gorm:"default:2"`
	CartID   uuid.UUID `gorm:"default:null"`
	Cart     Cart      `gorm:"foreignKey:CartID"` // Define the relationship
}
