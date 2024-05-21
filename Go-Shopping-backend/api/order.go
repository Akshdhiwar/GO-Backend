package api

import (
	"Go-Shopping-backend/controller"
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/middleware"
	"Go-Shopping-backend/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/product"
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
	var productsArray []models.OrderedProduct

	for _, lineProduct := range lineItems {

		var orderProduct models.OrderedProduct

		params := &stripe.ProductParams{}
		result, err := product.Get(lineProduct.Price.Product.ID, params)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "Error retriving product from stripe")
			return
		}

		orderProduct.ProductName = result.Name
		orderProduct.Quantity = lineProduct.Quantity
		amount := lineProduct.AmountTotal / 100
		orderProduct.Price = amount
		totalAmount = totalAmount + amount
		productsArray = append(productsArray, orderProduct)

	}

	_, err = initializers.DB.Exec(context.Background(), "INSERT INTO orders (email , products , name , total_amount , status) VALUES ($1, $2 , $3 , $4 , $5)", email, productsArray, userName, totalAmount, models.OrderStatus.Processing)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error creating order",
		})
	}

}

func OrderRoute(router *gin.RouterGroup) {
	router.GET("", middleware.Authenticate, middleware.RoleBasedAuthorization, controller.GetOrder)
}
