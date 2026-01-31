package routes

import (
	"foglio/v2/src/handlers"

	"github.com/gin-gonic/gin"
)

func AnalyticsRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	analytics := router.Group("/analytics")
	handler := handlers.NewAnalyticsHandler()

	// ==================== TRACKING ENDPOINTS (mostly public) ====================
	track := analytics.Group("/track")
	{
		track.POST("/page-view", handler.TrackPageView())
		track.POST("/job-view", handler.TrackJobView())
		track.POST("/profile-view", handler.TrackProfileView())
		track.POST("/portfolio-view", handler.TrackPortfolioView())
		track.POST("/event", handler.TrackEvent())
	}

	// ==================== ADMIN ANALYTICS (admin only) ====================
	admin := analytics.Group("/admin")
	{
		admin.GET("/dashboard", handler.GetAdminDashboard())
		admin.GET("/overview", handler.GetPlatformOverview())
		admin.GET("/users", handler.GetUserAnalytics())
		admin.GET("/jobs", handler.GetJobAnalytics())
		admin.GET("/applications", handler.GetApplicationAnalytics())
		admin.GET("/revenue", handler.GetRevenueAnalytics())
	}

	// ==================== RECRUITER ANALYTICS (recruiter/admin only) ====================
	recruiter := analytics.Group("/recruiter")
	{
		recruiter.GET("/dashboard", handler.GetRecruiterDashboard())
		recruiter.GET("/jobs", handler.GetRecruiterJobPerformance())
		recruiter.GET("/applications", handler.GetRecruiterApplicationStats())
	}

	// ==================== TALENT ANALYTICS (authenticated users) ====================
	talent := analytics.Group("/talent")
	{
		talent.GET("/dashboard", handler.GetTalentDashboard())
		talent.GET("/profile-views", handler.GetProfileViewsAnalytics())
		talent.GET("/portfolio", handler.GetPortfolioAnalytics())
		talent.GET("/applications", handler.GetTalentApplicationStats())
		talent.GET("/viewer-insights", handler.GetViewerInsights())
	}

	return analytics
}
