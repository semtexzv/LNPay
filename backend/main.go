package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go"
	"lnpay/kube"

	"lnpay/api"
	"os"
)

func init() {
	godotenv.Load()
	apiKey := os.Getenv("STRIPE_SECRET_KEY")
	stripe.Key = apiKey
}

// ChargeJSON incoming data for Stripe API
type ChargeJSON struct {
	Amount       int64  `json:"amount"`
	ReceiptEmail string `json:"receiptEmail"`
}

func main() {
	// set up server
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"*"},

	}))
	kube.GetSvcEndpoint("lnd")

	r.GET("/api/v1/metadata" , api.GetMetadata)
	r.POST("/api/v1/payment", api.CreatePayment)
	r.POST("/api/v1/payment/:payment_id/charge", api.FinishPayment)
	r.Run(":8080")
}
