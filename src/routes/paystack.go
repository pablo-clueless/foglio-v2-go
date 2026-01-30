package routes

import (
	"foglio/v2/src/handlers"

	"github.com/gin-gonic/gin"
)

func PaystackRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	handler := handlers.NewPaystackHandler()

	payments := router.Group("/payments")

	payments.POST("/initialize", handler.InitiatePayment())
	payments.GET("/verify", handler.VerifyPayment())
	payments.DELETE("/cancel", handler.CancelSubscription())
	payments.POST("/webhook", handler.Webhook())

	// Payment Methods
	payments.GET("/methods", handler.GetPaymentMethods())
	payments.POST("/methods", handler.AddPaymentMethod())
	payments.DELETE("/methods/:authCode", handler.RemovePaymentMethod())

	// Invoices
	payments.GET("/invoices", handler.GetInvoices())
	payments.GET("/invoices/:id", handler.GetInvoice())

	return payments
}
