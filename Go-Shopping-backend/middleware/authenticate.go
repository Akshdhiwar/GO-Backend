package middleware

import (
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authenticate(context *gin.Context) {
	tokenString, err := context.Cookie("Authorization")

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("JWTSECRET")), nil
	})
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			context.JSON(http.StatusUnauthorized, gin.H{
				"message": "token expired",
			})
		}
		var user models.User

		if err := initializers.DB.First(&user, claims["sub"]).Error; err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
			})
		}

		if user.ID == 0 {
			context.JSON(http.StatusUnauthorized, gin.H{
				"message": "user not found",
			})
		}

		context.Set("user", user)
		context.Set("userEmail", user.Email)

		context.Request.Header.Set("X-User-Email", user.Email)

		context.Next()
	} else {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
	}

}
