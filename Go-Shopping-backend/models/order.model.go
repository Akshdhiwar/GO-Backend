package models

import (
	"time"
)

type Order struct {
	ID          uint64           `json:"id"`
	Email       string           `json:"email"`
	Products    []OrderedProduct `json:"products"`
	CreatedAt   time.Time        `json:"created_at"`
	Name        string           `json:"name"`
	TotalAmount float64          `json:"total_amount"`
	Status      string           `json:"status"`
}

type OrderedProduct struct {
	ProductName string `json:"product_name"`
	Quantity    int64  `json:"quantity"`
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
