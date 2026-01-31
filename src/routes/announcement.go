package routes

import (
	"foglio/v2/src/handlers"
	"foglio/v2/src/lib"

	"github.com/gin-gonic/gin"
)

func AnnouncementRoutes(router *gin.RouterGroup, hub *lib.Hub) *gin.RouterGroup {
	handler := handlers.NewAnnouncementHandler(hub)

	// User-facing routes
	announcements := router.Group("/announcements")
	announcements.GET("", handler.GetAnnouncements())
	announcements.GET("/banners", handler.GetActiveBanners())
	announcements.PUT("/:id/read", handler.MarkAsRead())
	announcements.PUT("/:id/dismiss", handler.DismissAnnouncement())

	// Admin routes
	admin := router.Group("/admin/announcements")
	admin.GET("", handler.GetAnnouncementsAdmin())
	admin.GET("/:id", handler.GetAnnouncementAdmin())
	admin.POST("", handler.CreateAnnouncement())
	admin.PUT("/:id", handler.UpdateAnnouncement())
	admin.PUT("/:id/publish", handler.PublishAnnouncement())
	admin.PUT("/:id/unpublish", handler.UnpublishAnnouncement())
	admin.DELETE("/:id", handler.DeleteAnnouncement())

	return announcements
}
