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

	var jwtSecret string

	if os.Getenv("ENVIRONMENT") == "LOCAL" {
		jwtSecret = os.Getenv("JWTSECRET")
	} else {
		jwtSecret = os.Getenv("JWTSECRET_PROD")
	}

	log.Println(os.Getenv("ENVIRONMENT"))

	// Validate token
	// convert strign to a byte array
	email, err := parseJWTToken(token, []byte(jwtSecret))

	if err != nil {
		log.Printf("Error parsing token: %s", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
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
