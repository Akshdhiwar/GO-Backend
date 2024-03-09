package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
	Role     int  `gorm:"default:2"`
	CartID   uint // Foreign key for the Cart
}
