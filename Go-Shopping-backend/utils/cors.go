package utils

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	// Cors config
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "https://dumbles.vercel.app"} // specify the exact origin you want to allow
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowCredentials = true           // if you want to allow cookies to be sent to the server
	config.AllowHeaders = []string{"Origin"} // specify the allowed headers

	return cors.New(config)
}
