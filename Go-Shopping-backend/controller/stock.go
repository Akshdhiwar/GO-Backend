package controller

import (
	"Go-Shopping-backend/database"
	"Go-Shopping-backend/initializers"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UpdateStocks(ctx *gin.Context) {
	id := ctx.Param("id")
	log.Println(id)

	productId, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error parsing UUID",
			"type":    "error",
		})
		return
	}

	// body
	var body struct {
		Unit int `json:"unit"`
	}

	err = ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while binding body data",
			"type":    "error",
		})
		return
	}

	_, err = initializers.DB.Exec(context.Background(), database.UpdateStock, body.Unit, productId)
	if err != nil {
		log.Fatalf("Error creating new product stock: %v", err)
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "stocks added to product",
	})

}

func GetStocks(ctx *gin.Context) {
	id := ctx.Param("id")
	log.Println(id)

	productId, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error parsing UUID",
			"type":    "error",
		})
		return
	}

	var unit int
	err = initializers.DB.QueryRow(context.Background(), database.GetUnits, productId).Scan(&unit)
	if err != nil {
		log.Fatalf("Error getting product stock: %v", err)
	}

	ctx.JSON(http.StatusOK, unit)

}
