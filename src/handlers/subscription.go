package handlers

import (
	"foglio/v2/src/config"
	"foglio/v2/src/database"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/services"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	service *services.SubscriptionService
}

func NewSubscriptionHandler() *SubscriptionHandler {
	return &SubscriptionHandler{
		service: services.NewSubscriptionService(database.GetDatabase()),
	}
}

func (h *SubscriptionHandler) CreateSubscription() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.CreateSubscriptionDto

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		subscription, err := h.service.CreateSubscriptionTier(payload)
		if err != nil {
			lib.InternalServerError(ctx, "Failed to create subscription: "+err.Error())
			return
		}

		lib.Created(ctx, "Subscription created successfully", subscription)
	}
}

func (h *SubscriptionHandler) GetSubscriptions() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var query dto.Pagination

		if err := ctx.ShouldBindQuery(&query); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		subscriptions, err := h.service.GetSubscriptions(&query)
		if err != nil {
			lib.InternalServerError(ctx, "Failed to fetch subscriptions: "+err.Error())
			return
		}

		lib.Success(ctx, "Subscriptions fetched successfully", subscriptions)
	}
}

func (h *SubscriptionHandler) GetSubscription() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		subscription, err := h.service.GetSubscriptionById(id)
		if err != nil {
			lib.NotFound(ctx, "Subscription not found", "")
			return
		}

		lib.Success(ctx, "Subscription fetched successfully", subscription)
	}
}

func (h *SubscriptionHandler) UpdateSubscription() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		var payload dto.UpdateSubscriptionDto

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		subscription, err := h.service.UpdateSubscriptionTier(id, payload)
		if err != nil {
			lib.InternalServerError(ctx, "Failed to update subscription: "+err.Error())
			return
		}

		lib.Success(ctx, "Subscription updated successfully", subscription)
	}
}

func (h *SubscriptionHandler) DeleteSubscription() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		if err := h.service.DeleteSubscriptionTier(id); err != nil {
			lib.InternalServerError(ctx, "Failed to delete subscription: "+err.Error())
			return
		}

		lib.Success(ctx, "Subscription deleted successfully", nil)
	}
}

func (h *SubscriptionHandler) GetUserSubscriptions() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var query dto.Pagination
		if err := ctx.ShouldBindQuery(&query); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		subscriptions, err := h.service.GetUserSubscriptions(userId, &query)
		if err != nil {
			lib.InternalServerError(ctx, "Failed to fetch user subscriptions: "+err.Error())
			return
		}

		lib.Success(ctx, "User subscriptions fetched successfully", subscriptions)
	}
}

func (h *SubscriptionHandler) GetUserSubscription() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		subscription, err := h.service.GetUserSubscriptionById(id)
		if err != nil {
			lib.NotFound(ctx, "User subscription not found", "")
			return
		}

		lib.Success(ctx, "User subscription fetched successfully", subscription)
	}
}

func (h *SubscriptionHandler) Subscribe() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		tierId := ctx.Param("tierId")

		if err := h.service.SubscribeUser(userId, tierId); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		lib.Created(ctx, "Subscribed successfully", nil)
	}
}

func (h *SubscriptionHandler) Upgrade() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		newTierId := ctx.Param("tierId")

		if err := h.service.UpgradeUserSubscription(userId, newTierId); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		lib.Success(ctx, "Subscription upgraded successfully", nil)
	}
}

func (h *SubscriptionHandler) Downgrade() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		newTierId := ctx.Param("tierId")

		if err := h.service.DowngradeUserSubscription(userId, newTierId); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		lib.Success(ctx, "Subscription downgraded successfully", nil)
	}
}

func (h *SubscriptionHandler) Unsubscribe() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		if err := h.service.UnsubscribeUser(userId); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		lib.Success(ctx, "Unsubscribed successfully", nil)
	}
}
