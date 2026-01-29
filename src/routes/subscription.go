package routes

import (
	"foglio/v2/src/handlers"

	"github.com/gin-gonic/gin"
)

func SubscriptionRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	handler := handlers.NewSubscriptionHandler()

	// Subscription tiers (admin routes)
	subscriptions := router.Group("/subscriptions")
	subscriptions.POST("", handler.CreateSubscription())
	subscriptions.GET("", handler.GetSubscriptions())
	subscriptions.GET("/:id", handler.GetSubscription())
	subscriptions.PUT("/:id", handler.UpdateSubscription())
	subscriptions.DELETE("/:id", handler.DeleteSubscription())

	// User subscription actions
	userSubs := router.Group("/user/subscriptions")
	userSubs.GET("", handler.GetUserSubscriptions())
	userSubs.GET("/:id", handler.GetUserSubscription())
	userSubs.POST("/:tierId/subscribe", handler.Subscribe())
	userSubs.PUT("/:tierId/upgrade", handler.Upgrade())
	userSubs.PUT("/:tierId/downgrade", handler.Downgrade())
	userSubs.DELETE("/unsubscribe", handler.Unsubscribe())

	return subscriptions
}
