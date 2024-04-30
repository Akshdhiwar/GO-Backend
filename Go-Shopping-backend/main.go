package main

import (
	"Go-Shopping-backend/api"
	"Go-Shopping-backend/initializers"

	// "Go-Shopping-backend/middleware"
	"Go-Shopping-backend/utils"
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
	stripe.Key = os.Getenv("RAILS_STRIPE_SECREZT_KEY")
}

func main() {
	// Create a new Gin router
	router := gin.Default()

	router.Use(utils.Cors())

	//router.Use(middleware.RateLimitMiddleware())

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

	type Product struct {
		PriceID  string `json:"price_id"`
		Quantity int64  `json:"quantity"`
	}

	var body struct {
		Products []Product `json:"products"`
	}

	err := ctx.ShouldBind(&body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	lineItems := []*stripe.CheckoutSessionLineItemParams{}

	for _, item := range body.Products {
		stripeProduct := &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(item.PriceID),
			Quantity: stripe.Int64(item.Quantity),
		}

		lineItems = append(lineItems, stripeProduct)
	}

	params := &stripe.CheckoutSessionParams{
		LineItems:  lineItems,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(domain + "?success=true"),
		CancelURL:  stripe.String(domain + "?canceled=true"),
	}

	s, err := session.New(params)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create checkout session"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"url": s.URL})
}
