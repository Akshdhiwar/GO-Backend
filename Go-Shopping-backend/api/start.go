package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Default(router *gin.RouterGroup) {
	// just to check is server running or not
	router.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "Server is running",
		})
	})
}
