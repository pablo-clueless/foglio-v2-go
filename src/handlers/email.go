package handlers

import (
	"foglio/v2/src/config"
	"foglio/v2/src/lib"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EmailHandler struct {
	service *lib.EmailService
}

func NewEmailHandler() *EmailHandler {
	return &EmailHandler{service: lib.GetEmailService()}
}

func (h *EmailHandler) TestEmail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload := lib.EmailDto{
			To:       []string{"oklavender@virgilian.com"},
			Subject:  "Test Email",
			Template: "test",
			Data: map[string]interface{}{
				"Name":             "Okla Vender",
				"Email":            "oklavender@virgilian.com",
				"MessageID":        uuid.New(),
				"Environment":      config.AppConfig.Environment,
				"SentAt":           time.Now().Format("2006-01-02 15:04:05"),
				"ConfirmURL":       config.AppConfig.ClientUrl + "/test-email-delivery",
				"SupportURL":       config.AppConfig.ClientUrl + "/help-center",
				"TrackingPixelURL": config.AppConfig.ClientUrl + "/",
			},
		}

		err := h.service.SendEmail(payload)

		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Test email sent successfully", payload)
	}
}
