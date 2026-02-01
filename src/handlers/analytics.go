package handlers

import (
	"foglio/v2/src/config"
	"foglio/v2/src/database"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/models"
	"foglio/v2/src/services"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	service *services.AnalyticsService
}

func NewAnalyticsHandler() *AnalyticsHandler {
	return &AnalyticsHandler{
		service: services.NewAnalyticsService(database.GetDatabase()),
	}
}

func (h *AnalyticsHandler) TrackPageView() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.TrackPageViewDto
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		var userIDPtr *string
		if userID != "" {
			userIDPtr = &userID
		}

		ipAddress := ctx.ClientIP()
		userAgent := ctx.GetHeader("User-Agent")

		if err := h.service.TrackPageView(userIDPtr, payload, ipAddress, userAgent); err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Page view tracked", nil)
	}
}

func (h *AnalyticsHandler) TrackJobView() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.TrackJobViewDto
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		var userIDPtr *string
		if userID != "" {
			userIDPtr = &userID
		}

		ipAddress := ctx.ClientIP()
		userAgent := ctx.GetHeader("User-Agent")

		if err := h.service.TrackJobView(userIDPtr, payload, ipAddress, userAgent); err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Job view tracked", nil)
	}
}

func (h *AnalyticsHandler) TrackProfileView() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.TrackProfileViewDto
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		var userIDPtr *string
		isRecruiter := false

		if userID != "" {
			userIDPtr = &userID

			if user, exists := ctx.Get("current_user"); exists {
				isRecruiter = user.(*models.User).IsRecruiter
			}
		}

		ipAddress := ctx.ClientIP()
		userAgent := ctx.GetHeader("User-Agent")

		if err := h.service.TrackProfileView(userIDPtr, payload, ipAddress, userAgent, isRecruiter); err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Profile view tracked", nil)
	}
}

func (h *AnalyticsHandler) TrackPortfolioView() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.TrackPortfolioViewDto
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		var userIDPtr *string
		if userID != "" {
			userIDPtr = &userID
		}

		ipAddress := ctx.ClientIP()
		userAgent := ctx.GetHeader("User-Agent")

		if err := h.service.TrackPortfolioView(userIDPtr, payload, ipAddress, userAgent); err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Portfolio view tracked", nil)
	}
}

func (h *AnalyticsHandler) TrackEvent() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.TrackEventDto
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		var userIDPtr *string
		if userID != "" {
			userIDPtr = &userID
		}

		if err := h.service.TrackEvent(userIDPtr, payload); err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Event tracked", nil)
	}
}

func (h *AnalyticsHandler) GetAdminDashboard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get("current_user")
		if !exists {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		currentUser := user.(*models.User)
		if !currentUser.IsAdmin {
			lib.Forbidden(ctx, "Admin access required")
			return
		}

		var params dto.AnalyticsQueryParams
		if err := ctx.ShouldBindQuery(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		analytics, err := h.service.GetAdminDashboardAnalytics(params)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Admin dashboard analytics retrieved", analytics)
	}
}

func (h *AnalyticsHandler) GetPlatformOverview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get("current_user")
		if !exists {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		currentUser := user.(*models.User)
		if !currentUser.IsAdmin {
			lib.Forbidden(ctx, "Admin access required")
			return
		}

		var params dto.AnalyticsQueryParams
		if err := ctx.ShouldBindQuery(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		analytics, err := h.service.GetAdminDashboardAnalytics(params)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Platform overview retrieved", analytics.Overview)
	}
}

func (h *AnalyticsHandler) GetUserAnalytics() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get("current_user")
		if !exists {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		currentUser := user.(*models.User)
		if !currentUser.IsAdmin {
			lib.Forbidden(ctx, "Admin access required")
			return
		}

		var params dto.AnalyticsQueryParams
		if err := ctx.ShouldBindQuery(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		analytics, err := h.service.GetAdminDashboardAnalytics(params)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "User analytics retrieved", analytics.UserStats)
	}
}

func (h *AnalyticsHandler) GetJobAnalytics() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get("current_user")
		if !exists {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		currentUser := user.(*models.User)
		if !currentUser.IsAdmin {
			lib.Forbidden(ctx, "Admin access required")
			return
		}

		var params dto.AnalyticsQueryParams
		if err := ctx.ShouldBindQuery(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		analytics, err := h.service.GetAdminDashboardAnalytics(params)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Job analytics retrieved", analytics.JobStats)
	}
}

func (h *AnalyticsHandler) GetApplicationAnalytics() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get("current_user")
		if !exists {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		currentUser := user.(*models.User)
		if !currentUser.IsAdmin {
			lib.Forbidden(ctx, "Admin access required")
			return
		}

		var params dto.AnalyticsQueryParams
		if err := ctx.ShouldBindQuery(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		analytics, err := h.service.GetAdminDashboardAnalytics(params)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Application analytics retrieved", analytics.ApplicationStats)
	}
}

func (h *AnalyticsHandler) GetRevenueAnalytics() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get("current_user")
		if !exists {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		currentUser := user.(*models.User)
		if !currentUser.IsAdmin {
			lib.Forbidden(ctx, "Admin access required")
			return
		}

		var params dto.AnalyticsQueryParams
		if err := ctx.ShouldBindQuery(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		analytics, err := h.service.GetAdminDashboardAnalytics(params)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Revenue analytics retrieved", analytics.RevenueStats)
	}
}

func (h *AnalyticsHandler) GetRecruiterDashboard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		user, exists := ctx.Get("current_user")
		if !exists || userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		currentUser := user.(*models.User)
		if !currentUser.IsRecruiter && !currentUser.IsAdmin {
			lib.Forbidden(ctx, "Recruiter access required")
			return
		}

		var params dto.AnalyticsQueryParams
		if err := ctx.ShouldBindQuery(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		analytics, err := h.service.GetRecruiterAnalytics(userID, params)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Recruiter dashboard analytics retrieved", analytics)
	}
}

func (h *AnalyticsHandler) GetRecruiterJobPerformance() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		user, exists := ctx.Get("current_user")
		if !exists || userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		currentUser := user.(*models.User)
		if !currentUser.IsRecruiter && !currentUser.IsAdmin {
			lib.Forbidden(ctx, "Recruiter access required")
			return
		}

		var params dto.AnalyticsQueryParams
		if err := ctx.ShouldBindQuery(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		analytics, err := h.service.GetRecruiterAnalytics(userID, params)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Recruiter job performance retrieved", analytics.JobPerformance)
	}
}

func (h *AnalyticsHandler) GetRecruiterApplicationStats() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		user, exists := ctx.Get("current_user")
		if !exists || userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		currentUser := user.(*models.User)
		if !currentUser.IsRecruiter && !currentUser.IsAdmin {
			lib.Forbidden(ctx, "Recruiter access required")
			return
		}

		var params dto.AnalyticsQueryParams
		if err := ctx.ShouldBindQuery(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		analytics, err := h.service.GetRecruiterAnalytics(userID, params)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Recruiter application stats retrieved", analytics.ApplicationStats)
	}
}

func (h *AnalyticsHandler) GetTalentDashboard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var params dto.AnalyticsQueryParams
		if err := ctx.ShouldBindQuery(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		analytics, err := h.service.GetTalentAnalytics(userID, params)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Talent dashboard analytics retrieved", analytics)
	}
}

func (h *AnalyticsHandler) GetProfileViewsAnalytics() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var params dto.AnalyticsQueryParams
		if err := ctx.ShouldBindQuery(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		analytics, err := h.service.GetTalentAnalytics(userID, params)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Profile views analytics retrieved", analytics.ProfileViews)
	}
}

func (h *AnalyticsHandler) GetPortfolioAnalytics() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var params dto.AnalyticsQueryParams
		if err := ctx.ShouldBindQuery(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		analytics, err := h.service.GetTalentAnalytics(userID, params)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Portfolio analytics retrieved", analytics.PortfolioStats)
	}
}

func (h *AnalyticsHandler) GetTalentApplicationStats() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var params dto.AnalyticsQueryParams
		if err := ctx.ShouldBindQuery(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		analytics, err := h.service.GetTalentAnalytics(userID, params)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Talent application stats retrieved", analytics.ApplicationStats)
	}
}

func (h *AnalyticsHandler) GetViewerInsights() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var params dto.AnalyticsQueryParams
		if err := ctx.ShouldBindQuery(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		analytics, err := h.service.GetTalentAnalytics(userID, params)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Viewer insights retrieved", analytics.ViewerInsights)
	}
}
