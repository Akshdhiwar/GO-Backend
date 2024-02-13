package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Default(router *gin.RouterGroup) {
	// just to check is server running or not
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
}
