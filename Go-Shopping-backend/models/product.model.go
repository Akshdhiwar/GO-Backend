package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Title       string         `gorm:"unique" json:"title"`
	Price       float64        `json:"price"`
	Description string         `json:"description"`
	Category    string         `json:"category"`
	Image       string         `json:"image"`
	Rating      float32        `json:"rate"`
	Count       int            `json:"count"`
}
