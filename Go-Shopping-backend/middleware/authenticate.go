package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var emailCtxKey = "email"

func Authenticate(c *gin.Context) {
	hmacSecret := "WWfHYuaouEZyeedr+hOAVvnyM9/Lu1aCKmZh4F7IEe6Mb4zo6nkwCK4vd2ajNwmiOud4R5sr9tfTP57gA/0Z9g=="
	// Read the Authorization header
	token := c.GetHeader("Authorization")
	log.Println(token)
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Validate token
	// convert strign to a byte array
	email, err := parseJWTToken(token, []byte(hmacSecret))

	if err != nil {
		log.Printf("Error parsing token: %s", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	log.Printf("Received request from %s", email)

	// Save the email in the context to use later in the handler
	ctx := context.WithValue(c, emailCtxKey, email)
	c.Request = c.Request.WithContext(ctx)

	// Authenticated. Continue (call next handler)
	c.Next()
}

// List of claims that we want to parse from the JWT token.
// The RegisteredClaims struct contains the standard claims.
// See https://pkg.go.dev/github.com/golang-jwt/jwt/v5#RegisteredClaims
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// This function parses the JWT token and returns the email claim
func parseJWTToken(token string, hmacSecret []byte) (email string, err error) {
	// Parse the token and validate the signature
	t, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	// Check if the token is valid
	if err != nil {
		return "", fmt.Errorf("error validating token: %v", err)
	} else if claims, ok := t.Claims.(*Claims); ok {
		return claims.Email, nil
	}

	return "", fmt.Errorf("error parsing token: %v", err)
}
