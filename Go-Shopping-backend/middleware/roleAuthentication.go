package middleware

import (
	"Go-Shopping-backend/initializers"
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func RoleBasedAuthorization(ctx *gin.Context) {
	userEmail, _ := ctx.Get("userEmail")
	if userEmail == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User email not found in context"})
		ctx.Abort()
		return
	}

	// Retrieve user role from database based on email
	userRole, err := GetUserRoleByEmail(userEmail.(string)) // Assuming you have a function to fetch user role from the database based on email
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user role"})
		ctx.Abort()
		return
	}

	// Check if user has the required role
	if userRole != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Call the next handler
	ctx.Next()
}

// GetUserRoleByEmail retrieves the role of a user by email from the database
func GetUserRoleByEmail(email string) (int, error) {
	var role int
	err := initializers.DB.QueryRow(context.Background(), "SELECT role FROM users WHERE email = $1", email).Scan(&role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errors.New("user not found")
		}
		return 0, err
	}
	return role, nil
}
