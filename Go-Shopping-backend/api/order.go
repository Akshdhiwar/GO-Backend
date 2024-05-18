package api

import (
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
)

func CreateOrder(ctx *gin.Context, lineItems []*stripe.LineItem, email string) {

	var userName string

	err := initializers.DB.QueryRow(context.Background(), "SELECT full_name FROM users WHERE email=$1", email).Scan(&userName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error getting username from epecific email",
		})
	}

	var totalAmount int64

	for _, product := range lineItems {
		amount := product.AmountTotal / 100
		totalAmount = totalAmount + amount
	}

	_, err = initializers.DB.Exec(context.Background(), "INSERT INTO orders (email , products , name , total_amount , status) VALUES ($1, $2 , $3 , $4 , $5)", email, lineItems, userName, totalAmount, models.OrderStatus.Processing)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error creating order",
		})
	}

}
