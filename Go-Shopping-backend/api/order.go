package api

import (
	"Go-Shopping-backend/initializers"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
)

func CreateOrder(ctx *gin.Context, lineItems []*stripe.LineItem, email string) {

	log.Println(email)

	_, err := initializers.DB.Exec(context.Background(), "INSERT INTO orders (email , products) VALUES ($1, $2)", email, lineItems)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error creating order",
		})
	}

}
