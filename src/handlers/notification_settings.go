package handlers

import (
	"foglio/v2/src/config"
	"foglio/v2/src/database"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/services"

	"github.com/gin-gonic/gin"
)

type NotificationSettingsHandler struct {
	service *services.NotificationSettingsService
}

func NewNotificationSettingsHandler() *NotificationSettingsHandler {
	return &NotificationSettingsHandler{
		service: services.NewNotificationSettingsService(database.GetDatabase()),
	}
}

func (h *NotificationSettingsHandler) GetSettings() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		settings, err := h.service.GetOrCreateSettings(userID)
		if err != nil {
			lib.InternalServerError(ctx, "Failed to retrieve notification settings: "+err.Error())
			return
		}

		lib.Success(ctx, "Notification settings retrieved successfully", settings)
	}
}

func (h *NotificationSettingsHandler) UpdateSettings() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var payload dto.UpdateNotificationSettingsDto

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		settings, err := h.service.UpdateSettings(userID, payload)
		if err != nil {
			lib.InternalServerError(ctx, "Failed to update notification settings: "+err.Error())
			return
		}

		lib.Success(ctx, "Notification settings updated successfully", settings)
	}
}

func (h *NotificationSettingsHandler) UpdateEmailSettings() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var payload dto.UpdateEmailSettingsDto

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		settings, err := h.service.UpdateSettings(userID, dto.UpdateNotificationSettingsDto{
			Email: &payload,
		})
		if err != nil {
			lib.InternalServerError(ctx, "Failed to update email settings: "+err.Error())
			return
		}

		lib.Success(ctx, "Email notification settings updated successfully", settings)
	}
}

func (h *NotificationSettingsHandler) UpdatePushSettings() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var payload dto.UpdatePushSettingsDto

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		settings, err := h.service.UpdateSettings(userID, dto.UpdateNotificationSettingsDto{
			Push: &payload,
		})
		if err != nil {
			lib.InternalServerError(ctx, "Failed to update push settings: "+err.Error())
			return
		}

		lib.Success(ctx, "Push notification settings updated successfully", settings)
	}
}

func (h *NotificationSettingsHandler) UpdateInAppSettings() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var payload dto.UpdateInAppSettingsDto

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		settings, err := h.service.UpdateSettings(userID, dto.UpdateNotificationSettingsDto{
			InApp: &payload,
		})
		if err != nil {
			lib.InternalServerError(ctx, "Failed to update in-app settings: "+err.Error())
			return
		}

		lib.Success(ctx, "In-app notification settings updated successfully", settings)
	}
}
