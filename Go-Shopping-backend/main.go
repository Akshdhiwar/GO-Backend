package main

import (
	"Go-Shopping-backend/api"
	"Go-Shopping-backend/initializers"
	"Go-Shopping-backend/middleware"
	"Go-Shopping-backend/utils"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

func init() {
	// only load the .env file when running locally
	// check for a RAILWAY_ENVIRONMENT, if not found, code is running locally
	if _, exists := os.LookupEnv("RAILWAY_ENVIRONMENT"); !exists {
		initializers.LoadEnvVariables()
	}
	initializers.ConnectToDB()
	initializers.ConnectToRedis()
	initializers.LoadProductsToRedis()
	stripe.Key = os.Getenv("STRIPE_SECREZT_KEY")
}

func main() {
	// Create a new Gin router
	router := gin.Default()

	router.Use(utils.Cors())

	router.Use(middleware.RateLimitMiddleware())

	baseRoute := "api/v1"

	//default route
	api.Default(router.Group(baseRoute))

	// api route for Signup and Login
	// api.AccountRoutes(router.Group(baseRoute + "/account"))
	// have not using this api because he have implemented supabase for authentication

	// api route for Products like add , get, update , delete
	api.ProductRoutes(router.Group(baseRoute + "/products"))

	// api route for Cart
	api.CartRouter(router.Group(baseRoute + "/cart"))

	router.POST(baseRoute+"/create-checkout-session", createCheckoutSession)

	// Run the server on port 3000
	router.Run()
}

func createCheckoutSession(ctx *gin.Context) {
	domain := "http://localhost:5173"
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				// Provide the exact Price ID (for example, pr_1234) of the product you want to sell
				Price:    stripe.String("price_1P9ALnP5EFXn0qOIuuy6AMRR"),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(domain + "?success=true"),
		CancelURL:  stripe.String(domain + "?canceled=true"),
	}

	s, err := session.New(params)

	if err != nil {
		log.Printf("session.New: %v", err)
	}

	ctx.JSON(http.StatusOK, s.URL)
}
