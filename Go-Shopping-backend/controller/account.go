package controller

import (
	"Go-Shopping-backend/database"
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/models"
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Signup(ctx *gin.Context) {
	// getting body

	var body struct {
		Email     string
		Password  string
		FirstName string
		LastName  string
	}
	ctx.Bind(&body)

	// is any of field is empty condition check

	if body.Email == "" || body.Password == "" || body.FirstName == "" || body.LastName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get Email and Password",
		})
		return
	}

	var users models.User
	initializers.DB.QueryRow(context.Background(), database.SelectUserIdFromEmail, body.Email).Scan(&users.ID)
	if users.ID != uuid.Nil {
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
	user := models.User{Email: body.Email, Password: string(hash), FirstName: body.FirstName, LastName: body.LastName}

	// Execute the SQL query
	_, err = initializers.DB.Exec(context.Background(), database.SaveUserPassword, user.Email, user.Password, user.FirstName, body.LastName)
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
	err := initializers.DB.QueryRow(context.Background(), database.SelectUserDetailsFromEmail, body.Email).Scan(&users.ID, &users.Password, &users.FirstName, &users.LastName)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Email or Password incorrect",
		})
		return
	}

	if users.ID == uuid.Nil {
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
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWTSECRET")))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error creating token",
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Login Successful",
		"access_token": tokenString,
		"user_Details": users,
	})

}

func GetUserRole(ctx *gin.Context) {

	id := ctx.Param("id")

	var user models.User

	err := initializers.DB.QueryRow(context.Background(), "Select role from users WHERE id = $1", id).Scan(&user.Role)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"type":    "error",
		})
		return
	}

	ctx.JSON(http.StatusOK, user.Role)
}

func GetUserData(ctx *gin.Context) {

	id := ctx.Param("id")

	var user struct {
		name   string
		email  string
		avatar string
	}

	err := initializers.DB.QueryRow(context.Background(), "SELECT full_name , avatar_url , email from users WHERE id = $1", id).Scan(&user.name, &user.avatar, &user.email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"type":    "error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"email":  user.email,
		"name":   user.name,
		"avatar": user.avatar,
	})
}

func UpdateUserData(ctx *gin.Context) {
	id := ctx.Param("id")

	var body struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	err := ctx.ShouldBind(&body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Error while binding data to body")
	}

	if body.Email == "" || body.Name == "" {
		ctx.JSON(http.StatusInternalServerError, "please send all the required fields")
	}

	_, err = initializers.DB.Exec(context.Background(), "UPDATE users SET name=$1 , email=$2 WHERE id=$3", body.Name, body.Email, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Error while Quering to DB")
	}

	ctx.JSON(http.StatusOK, "Updated Data Successfully")
}
