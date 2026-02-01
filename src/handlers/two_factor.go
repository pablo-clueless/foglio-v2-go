package handlers

import (
	"foglio/v2/src/config"
	"foglio/v2/src/database"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/services"

	"github.com/gin-gonic/gin"
)

type TwoFactorHandler struct {
	service     *services.TwoFactorService
	authService *services.AuthService
}

func NewTwoFactorHandler() *TwoFactorHandler {
	db := database.GetDatabase()
	return &TwoFactorHandler{
		service:     services.NewTwoFactorService(db),
		authService: services.NewAuthService(db),
	}
}

// Setup2FA godoc
// @Summary Start 2FA setup
// @Description Generate a TOTP secret and QR code URL for setting up 2FA
// @Tags 2FA
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Enable2FAResponse
// @Failure 400 {object} lib.Response
// @Failure 401 {object} lib.Response
// @Router /auth/2fa/setup [post]
func (h *TwoFactorHandler) Setup2FA() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		currentUser, err := h.authService.FindUserById(userID)
		if err != nil {
			lib.Unauthorized(ctx, "User not found")
			return
		}

		if currentUser.IsTwoFactorEnabled {
			lib.BadRequest(ctx, "2FA is already enabled", "2FA_ALREADY_ENABLED")
			return
		}

		response, err := h.service.GenerateSecret(currentUser)
		if err != nil {
			lib.InternalServerError(ctx, "Failed to generate 2FA secret: "+err.Error())
			return
		}

		lib.Success(ctx, "2FA setup initiated. Scan the QR code with your authenticator app.", response)
	}
}

// VerifySetup2FA godoc
// @Summary Verify and enable 2FA
// @Description Verify the TOTP code from the authenticator app and enable 2FA
// @Tags 2FA
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.Verify2FASetupRequest true "Verification code"
// @Success 200 {object} lib.Response
// @Failure 400 {object} lib.Response
// @Failure 401 {object} lib.Response
// @Router /auth/2fa/verify-setup [post]
func (h *TwoFactorHandler) VerifySetup2FA() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var payload dto.Verify2FASetupRequest
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, "Invalid request body", "INVALID_PAYLOAD")
			return
		}

		if err := h.service.VerifyAndEnable(userID, payload.Code); err != nil {
			lib.BadRequest(ctx, err.Error(), "VERIFICATION_FAILED")
			return
		}

		lib.Success(ctx, "2FA has been enabled successfully", nil)
	}
}

// Verify2FALogin godoc
// @Summary Verify 2FA during login
// @Description Verify the TOTP code or backup code during login to get the auth token
// @Tags 2FA
// @Accept json
// @Produce json
// @Param request body dto.Verify2FALoginRequest true "User ID and verification code"
// @Success 200 {object} dto.SigninResponse
// @Failure 400 {object} lib.Response
// @Failure 401 {object} lib.Response
// @Router /auth/2fa/verify [post]
func (h *TwoFactorHandler) Verify2FALogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.Verify2FALoginRequest
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, "Invalid request body", "INVALID_PAYLOAD")
			return
		}

		user, err := h.service.VerifyCode(payload.UserID, payload.Code)
		if err != nil {
			lib.Unauthorized(ctx, err.Error())
			return
		}

		// Generate token
		token, err := lib.GenerateToken(user.ID)
		if err != nil {
			lib.InternalServerError(ctx, "Failed to generate token")
			return
		}

		// Clear sensitive fields
		user.Password = ""
		user.Otp = ""

		lib.Success(ctx, "2FA verification successful", services.SigninResponse{
			User:  *user,
			Token: token,
		})
	}
}

// Disable2FA godoc
// @Summary Disable 2FA
// @Description Disable two-factor authentication (requires password confirmation)
// @Tags 2FA
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.Disable2FARequest true "Password confirmation"
// @Success 200 {object} lib.Response
// @Failure 400 {object} lib.Response
// @Failure 401 {object} lib.Response
// @Router /auth/2fa/disable [post]
func (h *TwoFactorHandler) Disable2FA() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		currentUser, err := h.authService.FindUserById(userID)
		if err != nil {
			lib.Unauthorized(ctx, "User not found")
			return
		}

		if !currentUser.IsTwoFactorEnabled {
			lib.BadRequest(ctx, "2FA is not enabled", "2FA_NOT_ENABLED")
			return
		}

		var payload dto.Disable2FARequest
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, "Invalid request body", "INVALID_PAYLOAD")
			return
		}

		if err := h.service.Disable2FA(userID, payload.Password); err != nil {
			lib.BadRequest(ctx, err.Error(), "DISABLE_FAILED")
			return
		}

		lib.Success(ctx, "2FA has been disabled successfully", nil)
	}
}

// RegenerateBackupCodes godoc
// @Summary Regenerate backup codes
// @Description Generate new backup codes (invalidates old ones, requires password)
// @Tags 2FA
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.Disable2FARequest true "Password confirmation"
// @Success 200 {object} dto.BackupCodesResponse
// @Failure 400 {object} lib.Response
// @Failure 401 {object} lib.Response
// @Router /auth/2fa/backup-codes [post]
func (h *TwoFactorHandler) RegenerateBackupCodes() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var payload dto.Disable2FARequest
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, "Invalid request body", "INVALID_PAYLOAD")
			return
		}

		codes, err := h.service.RegenerateBackupCodes(userID, payload.Password)
		if err != nil {
			lib.BadRequest(ctx, err.Error(), "REGENERATE_FAILED")
			return
		}

		lib.Success(ctx, "New backup codes generated. Store them securely.", dto.BackupCodesResponse{
			BackupCodes: codes,
			Message:     "These are your new backup codes. Each code can only be used once.",
		})
	}
}

// GetStatus godoc
// @Summary Get 2FA status
// @Description Get the current 2FA status for the authenticated user
// @Tags 2FA
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.TwoFactorStatusResponse
// @Failure 401 {object} lib.Response
// @Router /auth/2fa/status [get]
func (h *TwoFactorHandler) GetStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		status, err := h.service.GetStatus(userID)
		if err != nil {
			lib.InternalServerError(ctx, "Failed to get 2FA status")
			return
		}

		lib.Success(ctx, "2FA status retrieved", status)
	}
}
