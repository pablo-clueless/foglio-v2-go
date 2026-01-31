package routes

import (
	"foglio/v2/src/handlers"

	"github.com/gin-gonic/gin"
)

func NotificationSettingsRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	settings := router.Group("/notification-settings")
	handler := handlers.NewNotificationSettingsHandler()

	settings.GET("", handler.GetSettings())
	settings.PUT("", handler.UpdateSettings())
	settings.PUT("/email", handler.UpdateEmailSettings())
	settings.PUT("/push", handler.UpdatePushSettings())
	settings.PUT("/in-app", handler.UpdateInAppSettings())

	return settings
}
