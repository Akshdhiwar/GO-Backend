package models

import "github.com/google/uuid"

type Order struct {
	ID    uint64
	Items []OrderedProduct
	Email string
}

type OrderedProduct struct {
	ProductID uuid.UUID
	Quantity  int
}
