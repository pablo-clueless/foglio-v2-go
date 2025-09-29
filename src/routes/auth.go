package routes

import (
	"foglio/v2/src/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	auth := router.Group("/auth")
	handler := handlers.NewAuthHandler()

	auth.POST("/signup", handler.CreateUser())
	auth.POST("/signin", handler.Signin())
	auth.POST("/verification", handler.Verification())
	auth.POST("/update-password", handler.ChangePassword())
	auth.POST("/forgot-password", handler.ForgotPassword())
	auth.POST("/reset-password", handler.ResetPassword())

	return auth
}
