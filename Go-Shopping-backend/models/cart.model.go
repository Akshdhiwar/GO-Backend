package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	UserID   int
	Products []int32 `gorm:"type:integer[]" json:"-"`
}
