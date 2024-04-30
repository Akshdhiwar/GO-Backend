package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var emailCtxKey = "email"

func extractToken(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		// Invalid authorization header format
		return ""
	}
	return parts[1]
}

func Authenticate(c *gin.Context) {
	// Read the Authorization header
	token := c.GetHeader("Authorization")
	token = extractToken(token)

	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var email string
	var err error

	// Validate token
	// convert strign to a byte array
	if os.Getenv("ENVIRONMENT") == "LOCAL" {
		log.Println("in Local")
		email, err = parseJWTToken(token, []byte(os.Getenv("RAILS_JWTSECRET")))
	} else {
		log.Println("in PROD")
		email, err = parseJWTToken(token, []byte(os.Getenv("RAILS_JWTSECRET_PROD")))
	}

	if err != nil {
		log.Printf("Error parsing token: %s", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unautorized"})
		return
	}

	log.Printf("Received request from %s", email)
	c.Set("userEmail", email)

	c.Request.Header.Set("X-User-Email", email)

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
	log.Println(token)
	// Parse the token and validate the signature
	t, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	// Check if there's an error in parsing or validating the token
	if err != nil {
		return "", fmt.Errorf("error parsing or validating token: %v", err)
	}

	// Check if the token is valid and extract email from claims
	if claims, ok := t.Claims.(*Claims); ok && t.Valid {
		return claims.Email, nil
	}

	return "", fmt.Errorf("invalid token")
}
