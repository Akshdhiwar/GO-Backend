package routes

import (
	"Go-Shopping-backend/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProducts(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"products": store.Product,
	})
}

func AddProducts(context *gin.Context) {
	if err := context.BindJSON(&store.Product); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Products successfully received",
	})
}
