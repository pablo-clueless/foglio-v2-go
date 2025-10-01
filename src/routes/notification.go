package routes

import (
	"foglio/v2/src/handlers"

	"github.com/gin-gonic/gin"
)

func NotificationRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	notifications := router.Group("/notifications")
	handler := handlers.NewNotificationHandler()

	notifications.GET("", handler.GetNotifications())
	notifications.GET("/:id", handler.GetNotification())
	notifications.PUT("/:id", handler.ReadNotification())
	notifications.DELETE("/:id", handler.DeleteNotification())

	return notifications
}
