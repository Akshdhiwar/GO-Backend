package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID   uint  // Foreign key for the user
	Products []int `gorm:"type:integer[]"`
}
