package main

import (
	"Go-Shopping-backend/api"
	"Go-Shopping-backend/initializers"
	"encoding/json"
	"io"
	"log"

	// "Go-Shopping-backend/middleware"
	"Go-Shopping-backend/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/webhook"
)

var domain string

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
	if os.Getenv("RAILS_ENVIRONMENT") == "LOCAL" {
		domain = "http://localhost:5173/order-status"
	} else {
		domain = "https://dumbles.vercel.app/order-status"
	}
}

type WebhookData struct {
	Type string `json:"type"`
	// Add other fields if needed
}

const MaxBodyBytes = int64(65536)

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

	router.POST("/webhook", WebhookController)

	// Run the server on port 3000
	router.Run()
}

func WebhookController(ctx *gin.Context) {
	const MaxBodyBytes = int64(65536)
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, MaxBodyBytes)

	body, err := io.ReadAll(ctx.Request.Body)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "ERROR WHILE BINDING BODY")
		log.Println("ERROR WHILE BINDING BODY")
		return
	}

	endpointSecret := "whsec_40e53ab232abf2f63fb1e0f7d8d61c195b6532ca7776082bf8c223331cb1c44e"
	event, err := webhook.ConstructEvent(body, ctx.GetHeader("Stripe-Signature"), endpointSecret)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error verifing webhook signature",
		})
		log.Println("Error verifing webhook signature")
		return
	}

	if event.Type == "checkout.session.completed" {
		var CheckoutSession stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &CheckoutSession)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error Unmashling the event data",
			})
			return
		}

		params := &stripe.CheckoutSessionParams{}
		params.AddExpand("line_items")

		// Retrieve the session. If you require line items in the response, you may include them by expanding line_items.
		sessionWithLineItems, _ := session.Get(CheckoutSession.ID, params)
		lineItems := sessionWithLineItems.LineItems

		log.Println("LineItems", lineItems.Data)
	}

	log.Println(event.Type)
}

func createCheckoutSession(ctx *gin.Context) {

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

// router.GET("/webhook", func(c *gin.Context) {
// 	// Limit request body size
// 	const MaxBodyBytes = int64(65536)
// 	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

// 	body, err := io.ReadAll(c.Request.Body)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
// 		c.AbortWithStatus(http.StatusServiceUnavailable)
// 		return
// 	}

// 	// Pass the request body and Stripe-Signature header to ConstructEvent, along with the webhook signing key
// 	// You can find your endpoint's secret in your webhook settings
// 	endpointSecret := "whsec_40e53ab232abf2f63fb1e0f7d8d61c195b6532ca7776082bf8c223331cb1c44e"
// 	event, err := webhook.ConstructEvent(body, c.GetHeader("Stripe-Signature"), endpointSecret)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
// 		c.String(http.StatusBadRequest, "Error verifying webhook signature")
// 		return
// 	}

// 	// Handle the checkout.session.completed event
// 	if event.Type == "checkout.session.completed" {
// 		// var session stripe.CheckoutSession

// 		// jsonData, err := json.Marshal(session)
// 		// if err != nil {
// 		// 	c.JSON(http.StatusInternalServerError, gin.H{
// 		// 		"message": "Error marshlisng session",
// 		// 	})
// 		// 	return
// 		// }

// 		// err = event.Data.UnmarshalJSON(jsonData)
// 		// if err != nil {
// 		// 	fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
// 		// 	c.String(http.StatusBadRequest, "Error parsing webhook JSON")
// 		// 	return
// 		// }

// 		// params := &stripe.CheckoutSessionParams{}
// 		// params.AddExpand("line_items")

// 		// fmt.Println("in checkout session")

// 		// // Retrieve the session. If you require line items in the response, you may include them by expanding line_items.
// 		// sessionWithLineItems, _ := session.Get(session.ID, params)
// 		// lineItems := sessionWithLineItems.LineItems
// 		// // Fulfill the purchase...
// 		// log.Println(lineItems)

// 		log.Println("in checkout session")
// 	}

// 	c.String(http.StatusOK, "Webhook received")
// })
