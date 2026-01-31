package routes

import (
	"foglio/v2/src/handlers"

	"github.com/gin-gonic/gin"
)

func DomainRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	domain := router.Group("/domain")
	handler := handlers.NewDomainHandler()

	// Get user's domain configuration
	domain.GET("", handler.GetDomain())

	// Check subdomain availability (can be accessed without claiming)
	domain.GET("/check/:subdomain", handler.CheckSubdomainAvailability())

	// Subdomain operations (available to all users)
	domain.POST("/subdomain", handler.ClaimSubdomain())
	domain.PUT("/subdomain", handler.UpdateSubdomain())

	// Custom domain operations (paid users only)
	domain.POST("/custom", handler.SetCustomDomain())
	domain.POST("/custom/verify", handler.VerifyCustomDomain())
	domain.DELETE("/custom", handler.RemoveCustomDomain())

	return domain
}
