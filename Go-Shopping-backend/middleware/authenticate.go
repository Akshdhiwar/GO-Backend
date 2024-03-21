package middleware

import (
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

func Authenticate(ctx *gin.Context) {
	tokenString, err := ctx.Cookie("Authorization")

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		ctx.AbortWithStatus(http.StatusUnauthorized)
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "token expired",
			})
		}
		var user models.User

		err := initializers.DB.QueryRow(context.Background(), "SELECT * FROM users WHERE id = $1", claims["sub"]).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CartID)
		if err != nil {
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
				return
			}
			// Handle other errors if needed
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}

		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "user not found",
			})
		}

		ctx.Set("user", user)
		ctx.Set("userEmail", user.Email)

		ctx.Request.Header.Set("X-User-Email", user.Email)

		ctx.Next()
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
	}

}
