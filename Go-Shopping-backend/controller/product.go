package controller

import (
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProducts(context *gin.Context) {
	var products []models.Product
	result := initializers.DB.Find(&products)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Error querying products from database",
		})
		return
	}

	if len(products) == 0 {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "No products found",
		})
		return
	}

	context.JSON(http.StatusOK, products)
}

func AddProducts(context *gin.Context) {

	//getting body
	var body struct {
		Title       string
		Price       float64
		Description string
		Category    string
		Image       string
	}

	context.ShouldBind(&body)

	if body.Title == "" || body.Price == 0 || body.Image == "" || body.Category == "" || body.Description == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "All fields are required",
		})
		return
	}

	var product models.Product
	initializers.DB.First(&product, "title = ?", body.Title)

	if product.Title == body.Title {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Product already present please check Title. Title is matching with some other product",
		})
		return
	}

	newProduct := models.Product{Title: body.Title, Price: body.Price, Category: body.Category, Image: body.Image, Description: body.Description}
	result := initializers.DB.Create(&newProduct)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to save product to db",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Products successfully received",
	})
}
