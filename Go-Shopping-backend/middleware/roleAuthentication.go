package middleware

import (
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoleBasedAuthorization(context *gin.Context) {
	userEmail, _ := context.Get("userEmail")
	if userEmail == nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "User email not found in context"})
		context.Abort()
		return
	}

	// Retrieve user role from database based on email
	userRole, err := GetUserRoleByEmail(userEmail.(string)) // Assuming you have a function to fetch user role from the database based on email
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user role"})
		context.Abort()
		return
	}

	// Check if user has the required role
	if userRole != 1 {
		context.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Call the next handler
	context.Next()
}

// GetUserRoleByEmail retrieves the role of a user by email from the database
func GetUserRoleByEmail(email string) (int, error) {
	var user models.User
	if err := initializers.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("user not found")
		}
		return 0, err
	}
	return user.Role, nil
}
