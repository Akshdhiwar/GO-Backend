package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Title       string  `gorm:"unique" json:"title"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Image       string  `json:"image"`
	Rating      float32 `json:"rate"`
	Count       int     `json:"count"`
}
