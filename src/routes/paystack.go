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

	return payments
}
