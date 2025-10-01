package handlers

import (
	"foglio/v2/src/config"
	"foglio/v2/src/database"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/services"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service *services.NotificationService
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{
		service: services.NewNotificationService(database.GetDatabase(), lib.NewHub()),
	}
}

func (h *NotificationHandler) GetNotifications() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var query dto.Pagination
		id := ctx.GetString(config.AppConfig.CurrentUserId)

		if err := ctx.ShouldBindQuery(&query); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		notifications, err := h.service.GetNotifications(id, query)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error, "+err.Error())
			return
		}

		lib.Success(ctx, "Notifications fetched successfully", notifications)
	}
}

func (h *NotificationHandler) GetNotification() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		notification, err := h.service.GetNotification(id)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error, "+err.Error())
			return
		}

		lib.Success(ctx, "Notification fetched successfully", notification)
	}
}

func (h *NotificationHandler) DeleteNotification() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		err := h.service.DeleteNotification(id)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error, "+err.Error())
			return
		}

		lib.Success(ctx, "Notification deleted successfully", nil)
	}
}

func (h *NotificationHandler) ReadNotification() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		err := h.service.ReadNotification(id)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error, "+err.Error())
			return
		}

		lib.Success(ctx, "Notification marked as read successfully", nil)
	}
}
