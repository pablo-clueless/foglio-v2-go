package handlers

import (
	"foglio/v2/src/config"
	"foglio/v2/src/database"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/models"
	"foglio/v2/src/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AnnouncementHandler struct {
	service *services.AnnouncementService
}

func NewAnnouncementHandler(hub *lib.Hub) *AnnouncementHandler {
	return &AnnouncementHandler{
		service: services.NewAnnouncementService(database.GetDatabase(), hub),
	}
}

// ==================== ADMIN ENDPOINTS ====================
func (h *AnnouncementHandler) CreateAnnouncement() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get("current_user")
		if !exists {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		currentUser := user.(*models.User)
		if !currentUser.IsAdmin {
			lib.Forbidden(ctx, "Admin access required")
			return
		}

		var payload dto.CreateAnnouncementDto
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		announcement, err := h.service.CreateAnnouncement(currentUser.ID.String(), payload)
		if err != nil {
			lib.InternalServerError(ctx, "Failed to create announcement: "+err.Error())
			return
		}

		lib.Created(ctx, "Announcement created successfully", announcement)
	}
}

func (h *AnnouncementHandler) UpdateAnnouncement() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get("current_user")
		if !exists {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		currentUser := user.(*models.User)
		if !currentUser.IsAdmin {
			lib.Forbidden(ctx, "Admin access required")
			return
		}

		id := ctx.Param("id")
		var payload dto.UpdateAnnouncementDto

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		announcement, err := h.service.UpdateAnnouncement(id, payload)
		if err != nil {
			if err.Error() == "announcement not found" {
				lib.NotFound(ctx, err.Error(), "")
				return
			}
			lib.InternalServerError(ctx, "Failed to update announcement: "+err.Error())
			return
		}

		lib.Success(ctx, "Announcement updated successfully", announcement)
	}
}

func (h *AnnouncementHandler) PublishAnnouncement() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get("current_user")
		if !exists {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		currentUser := user.(*models.User)
		if !currentUser.IsAdmin {
			lib.Forbidden(ctx, "Admin access required")
			return
		}

		id := ctx.Param("id")

		announcement, err := h.service.PublishAnnouncement(id)
		if err != nil {
			if err.Error() == "announcement not found" {
				lib.NotFound(ctx, err.Error(), "")
				return
			}
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		lib.Success(ctx, "Announcement published successfully", announcement)
	}
}

func (h *AnnouncementHandler) UnpublishAnnouncement() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get("current_user")
		if !exists {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		currentUser := user.(*models.User)
		if !currentUser.IsAdmin {
			lib.Forbidden(ctx, "Admin access required")
			return
		}

		id := ctx.Param("id")

		announcement, err := h.service.UnpublishAnnouncement(id)
		if err != nil {
			if err.Error() == "announcement not found" {
				lib.NotFound(ctx, err.Error(), "")
				return
			}
			lib.InternalServerError(ctx, "Failed to unpublish announcement: "+err.Error())
			return
		}

		lib.Success(ctx, "Announcement unpublished successfully", announcement)
	}
}

func (h *AnnouncementHandler) DeleteAnnouncement() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get("current_user")
		if !exists {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		currentUser := user.(*models.User)
		if !currentUser.IsAdmin {
			lib.Forbidden(ctx, "Admin access required")
			return
		}

		id := ctx.Param("id")

		if err := h.service.DeleteAnnouncement(id); err != nil {
			if err.Error() == "announcement not found" {
				lib.NotFound(ctx, err.Error(), "")
				return
			}
			lib.InternalServerError(ctx, "Failed to delete announcement: "+err.Error())
			return
		}

		lib.Success(ctx, "Announcement deleted successfully", nil)
	}
}

func (h *AnnouncementHandler) GetAnnouncementAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get("current_user")
		if !exists {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		currentUser := user.(*models.User)
		if !currentUser.IsAdmin {
			lib.Forbidden(ctx, "Admin access required")
			return
		}

		id := ctx.Param("id")

		announcement, err := h.service.GetAnnouncement(id)
		if err != nil {
			if err.Error() == "announcement not found" {
				lib.NotFound(ctx, err.Error(), "")
				return
			}
			lib.InternalServerError(ctx, "Failed to retrieve announcement: "+err.Error())
			return
		}

		lib.Success(ctx, "Announcement retrieved successfully", announcement)
	}
}

func (h *AnnouncementHandler) GetAnnouncementsAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get("current_user")
		if !exists {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		currentUser := user.(*models.User)
		if !currentUser.IsAdmin {
			lib.Forbidden(ctx, "Admin access required")
			return
		}

		var params dto.AnnouncementQueryParams
		if err := ctx.ShouldBindQuery(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		announcements, err := h.service.GetAnnouncementsForAdmin(params)
		if err != nil {
			lib.InternalServerError(ctx, "Failed to retrieve announcements: "+err.Error())
			return
		}

		lib.Success(ctx, "Announcements retrieved successfully", announcements)
	}
}

// ==================== USER ENDPOINTS ====================
func (h *AnnouncementHandler) GetAnnouncements() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get("current_user")
		if !exists {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		currentUser := user.(*models.User)

		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

		announcements, err := h.service.GetAnnouncementsForUser(currentUser, page, limit)
		if err != nil {
			lib.InternalServerError(ctx, "Failed to retrieve announcements: "+err.Error())
			return
		}

		lib.Success(ctx, "Announcements retrieved successfully", announcements)
	}
}

func (h *AnnouncementHandler) GetActiveBanners() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get("current_user")
		if !exists {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		currentUser := user.(*models.User)

		banners, err := h.service.GetActiveBanners(currentUser)
		if err != nil {
			lib.InternalServerError(ctx, "Failed to retrieve banners: "+err.Error())
			return
		}

		lib.Success(ctx, "Banners retrieved successfully", banners)
	}
}

func (h *AnnouncementHandler) MarkAsRead() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		announcementID := ctx.Param("id")

		err := h.service.MarkAnnouncementAsRead(userID, announcementID)
		if err != nil {
			lib.InternalServerError(ctx, "Failed to mark announcement as read: "+err.Error())
			return
		}

		lib.Success(ctx, "Announcement marked as read", nil)
	}
}

func (h *AnnouncementHandler) DismissAnnouncement() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		announcementID := ctx.Param("id")

		err := h.service.DismissAnnouncement(userID, announcementID)
		if err != nil {
			lib.InternalServerError(ctx, "Failed to dismiss announcement: "+err.Error())
			return
		}

		lib.Success(ctx, "Announcement dismissed", nil)
	}
}
