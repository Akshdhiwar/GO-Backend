package controller

import (
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Signup(context *gin.Context) {
	// getting body

	var body struct {
		Email    string
		Password string
	}
	context.Bind(&body)

	// is any of field is empty condition check

	if body.Email == "" || body.Password == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get Email and Password",
		})
		return
	}

	// hashing the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to create hash",
		})
		return
	}

	// storing the hash password and email in user object and saving it to DB
	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to save user to db",
		})
		return
	}

	// sending success message

	context.JSON(http.StatusOK, gin.H{
		"message": "User created successfully.",
	})

}
