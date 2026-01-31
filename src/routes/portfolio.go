package routes

import (
	"foglio/v2/src/handlers"

	"github.com/gin-gonic/gin"
)

func PortfolioRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	portfolio := router.Group("/portfolio")
	handler := handlers.NewPortfolioHandler()

	// User's own portfolio management (authenticated)
	portfolio.POST("", handler.CreatePortfolio())
	portfolio.GET("", handler.GetPortfolio())
	portfolio.PUT("", handler.UpdatePortfolio())
	portfolio.DELETE("", handler.DeletePortfolio())

	// Publish/Unpublish
	portfolio.POST("/publish", handler.PublishPortfolio())
	portfolio.POST("/unpublish", handler.UnpublishPortfolio())

	// Sections management
	portfolio.POST("/sections", handler.CreateSection())
	portfolio.PUT("/sections/:sectionId", handler.UpdateSection())
	portfolio.DELETE("/sections/:sectionId", handler.DeleteSection())
	portfolio.POST("/sections/reorder", handler.ReorderSections())

	// Public portfolio access (by slug)
	portfolios := router.Group("/portfolios")
	portfolios.GET("/:slug", handler.GetPortfolioBySlug())

	return portfolio
}
