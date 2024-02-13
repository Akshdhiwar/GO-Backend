package controller

import (
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	var users models.User
	initializers.DB.First(&users, "email = ?", body.Email)

	if users.ID != 0 {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "User already present please login",
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

func Login(context *gin.Context) {
	// geting body

	var body struct {
		Email    string
		Password string
	}

	context.Bind(&body)

	// validating Email and Password is empty or not
	if body.Email == "" || body.Password == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Email or Password missing",
		})

		return
	}

	var users models.User
	initializers.DB.First(&users, "email = ?", body.Email)

	if users.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})

		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(body.Password))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Password incorrect",
		})

		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": users.ID,
		"exp": time.Now().Add(time.Hour * 8).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWTSECRET")))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Error creating token",
		})

		return
	}

	context.SetSameSite(http.SameSiteLaxMode)
	context.SetCookie("Authorization", tokenString, 3600*8, "", "", false, true)

	context.JSON(http.StatusOK, gin.H{
		"message": "Login Successfull",
	})

}
