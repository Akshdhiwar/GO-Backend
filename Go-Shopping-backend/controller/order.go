package controller

import (
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/models"
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetOrder(ctx *gin.Context) {
	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil || offset <= 0 {
		offset = 0
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil || limit <= 0 {
		limit = 16
	}

	if limit > 20 {
		limit = 20
	}

	var ordersList []models.Order

	offset = offset * limit

	rows, err := initializers.DB.Query(context.Background(), "SELECT * FROM orders LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error fetching orders from DB: " + err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.CreatedAt, &order.Email, &order.Products, &order.Name, &order.TotalAmount, &order.Status); err != nil {
			log.Printf("Error scanning order row: %v", err)
			continue
		}
		ordersList = append(ordersList, order)
	}

	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error iterating order rows: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ordersList)

}
