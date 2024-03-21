package controller

import (
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/models"
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(ctx *gin.Context) {
	// getting body

	var body struct {
		Email    string
		Password string
	}
	ctx.Bind(&body)

	// is any of field is empty condition check

	if body.Email == "" || body.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get Email and Password",
		})
		return
	}

	var users models.User
	initializers.DB.QueryRow(context.Background(), "SELECT id FROM users WHERE email=$1", body.Email).Scan(&users.ID)
	log.Println(users.ID)
	if users.ID != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "User already present please login",
		})
		return
	}
	// hashing the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to create hash",
		})
		return
	}

	// storing the hash password and email in user object and saving it to DB
	user := models.User{Email: body.Email, Password: string(hash)}
	// Construct the SQL query for inserting a new user
	query := `
        INSERT INTO users (email, password)
        VALUES ($1, $2)
    `
	// Execute the SQL query
	_, err = initializers.DB.Exec(ctx, query, user.Email, user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to save user to db",
		})
		return
	}

	// sending success message

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User created successfully.",
	})

}

func Login(ctx *gin.Context) {
	// geting body

	var body struct {
		Email    string
		Password string
	}

	ctx.Bind(&body)

	// validating Email and Password is empty or not
	if body.Email == "" || body.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Email or Password missing",
		})

		return
	}

	var users models.User
	err := initializers.DB.QueryRow(context.Background(), "SELECT id , password FROM users WHERE email=$1", body.Email).Scan(&users.ID, &users.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Email or Password incorrect",
		})
		return
	}

	log.Println(users.ID, users.Email, users.Password)
	if users.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(body.Password))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error creating token",
		})

		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", tokenString, 3600*8, "", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login Successfull",
	})

}
