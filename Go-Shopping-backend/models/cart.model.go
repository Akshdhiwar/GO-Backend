package models

import (
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	UserID   int
	Products []int `gorm:"type:integer[]" json:"-"`
}

// Scan implements the sql.Scanner interface for the Products field
func (c *Cart) Scan(value interface{}) error {
	if value == nil {
		c.Products = nil
		return nil
	}
	// Convert the value to string
	strValue, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("unexpected type for products: %T", value)
	}
	// Unmarshal the string value into a slice of integers
	return json.Unmarshal(strValue, &c.Products)
}
