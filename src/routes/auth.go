package routes

import (
	"foglio/v2/src/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	auth := router.Group("/auth")
	handler := handlers.NewAuthHandler()
	twoFactorHandler := handlers.NewTwoFactorHandler()

	auth.POST("/signup", handler.CreateUser())
	auth.POST("/signin", handler.Signin())
	auth.POST("/request-verification", handler.RequestVerification())
	auth.POST("/verification", handler.Verification())
	auth.POST("/update-password", handler.ChangePassword())
	auth.POST("/forgot-password", handler.ForgotPassword())
	auth.POST("/reset-password", handler.ResetPassword())
	auth.GET("/:provider", handler.GetOAuthURL())
	auth.GET("/:provider/callback", handler.HandleOAuthCallback())

	auth.POST("/2fa/verify", twoFactorHandler.Verify2FALogin())
	twoFactor := auth.Group("/2fa")
	{
		twoFactor.GET("/status", twoFactorHandler.GetStatus())
		twoFactor.POST("/setup", twoFactorHandler.Setup2FA())
		twoFactor.POST("/verify-setup", twoFactorHandler.VerifySetup2FA())
		twoFactor.POST("/disable", twoFactorHandler.Disable2FA())
		twoFactor.POST("/backup-codes", twoFactorHandler.RegenerateBackupCodes())
	}

	return auth
}

// https://foglio.onrender.com/api/v2/auth/github/callback
// https://foglio.onrender.com/api/v2/auth/github/webhook
