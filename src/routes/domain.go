package routes

import (
	"foglio/v2/src/handlers"

	"github.com/gin-gonic/gin"
)

func DomainRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	domain := router.Group("/domain")
	handler := handlers.NewDomainHandler()

	domain.GET("", handler.GetDomain())
	domain.GET("/check/:subdomain", handler.CheckSubdomainAvailability())
	domain.POST("/subdomain", handler.ClaimSubdomain())
	domain.PUT("/subdomain", handler.UpdateSubdomain())
	domain.POST("/custom", handler.SetCustomDomain())
	domain.POST("/custom/verify", handler.VerifyCustomDomain())
	domain.DELETE("/custom", handler.RemoveCustomDomain())

	return domain
}
