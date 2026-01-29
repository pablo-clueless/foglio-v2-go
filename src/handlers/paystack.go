package handlers

import (
	"encoding/json"
	"io"

	"foglio/v2/src/config"
	"foglio/v2/src/database"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/services"

	"github.com/gin-gonic/gin"
)

type PaystackHandler struct {
	service *services.PaystackService
}

func NewPaystackHandler() *PaystackHandler {
	return &PaystackHandler{
		service: services.NewPaystackService(database.GetDatabase()),
	}
}

func (h *PaystackHandler) InitiatePayment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var payload dto.InitiatePaymentDto
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		callbackURL := payload.CallbackURL
		if callbackURL == "" {
			callbackURL = config.AppConfig.ClientUrl + "/subscription/callback"
		}

		response, err := h.service.InitializeTransaction(userID, payload.SubscriptionTierID, callbackURL)
		if err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		lib.Success(ctx, "Payment initialized successfully", response)
	}
}

func (h *PaystackHandler) VerifyPayment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		reference := ctx.Query("reference")
		if reference == "" {
			lib.BadRequest(ctx, "Reference is required", "")
			return
		}

		txData, err := h.service.VerifyTransaction(reference)
		if err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		if txData.Status != "success" {
			lib.BadRequest(ctx, "Payment not successful: "+txData.GatewayResponse, "")
			return
		}

		if err := h.service.ProcessSuccessfulPayment(txData); err != nil {
			lib.InternalServerError(ctx, "Failed to activate subscription: "+err.Error())
			return
		}

		lib.Success(ctx, "Payment verified and subscription activated", map[string]interface{}{
			"status":    txData.Status,
			"reference": txData.Reference,
			"amount":    float64(txData.Amount) / 100,
			"currency":  txData.Currency,
		})
	}
}

func (h *PaystackHandler) Webhook() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Failed to read request body"})
			return
		}

		signature := ctx.GetHeader("x-paystack-signature")
		if signature == "" {
			ctx.JSON(400, gin.H{"error": "Missing signature"})
			return
		}

		if !h.service.VerifyWebhookSignature(body, signature) {
			ctx.JSON(401, gin.H{"error": "Invalid signature"})
			return
		}

		var event dto.PaystackWebhookEvent
		if err := json.Unmarshal(body, &event); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid JSON"})
			return
		}

		if err := h.service.HandleWebhook(&event); err != nil {
			ctx.JSON(200, gin.H{"status": "error", "message": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"status": "success"})
	}
}

func (h *PaystackHandler) CancelSubscription() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		if err := h.service.CancelUserSubscription(userID); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		lib.Success(ctx, "Subscription cancelled successfully", nil)
	}
}
