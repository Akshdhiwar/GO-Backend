package controller

import (
	"Go-Shopping-backend/database"
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetSingleProductAdmin(ctx *gin.Context) {

	// getting id from url
	id := ctx.Param("id")

	var product models.Product

	err := initializers.DB.QueryRow(context.Background(), database.SelectAllFromID, id).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Title, &product.Price, &product.Description, &product.Category, &product.Image, &product.Rating, &product.Count, &product.PriceID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No product found",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Error querying db",
		})
		return
	}

	ctx.JSON(http.StatusOK, product)

}
