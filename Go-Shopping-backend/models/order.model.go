package models

import "github.com/google/uuid"

type Order struct {
	ID    uint64
	Items []OrderedProduct
	Email string
	Name  string
}

type OrderedProduct struct {
	ProductID uuid.UUID
	Quantity  int
}

type Status struct {
	Processing     string
	OrderReceived  string
	Processed      string
	OutForDelivery string
	Delivered      string
}

const (
	Processing     = "Processing"
	OrderReceived  = "Order Received"
	Processed      = "Processed"
	OutForDelivery = "Out For Delivery"
	Delivered      = "Delivered"
)

var OrderStatus = Status{
	Processing:     Processing,
	OrderReceived:  OrderReceived,
	Processed:      Processed,
	OutForDelivery: OutForDelivery,
	Delivered:      Delivered,
}
