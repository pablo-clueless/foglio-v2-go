package services

import (
	"errors"
	"foglio/v2/src/dto"
	"foglio/v2/src/models"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrUnauthorizedAnalytics = errors.New("you are not authorized to view these analytics")
	ErrInvalidDateRange      = errors.New("invalid date range")
)

type AnalyticsService struct {
	database *gorm.DB
}

func NewAnalyticsService(database *gorm.DB) *AnalyticsService {
	return &AnalyticsService{
		database: database,
	}
}

// ==================== TRACKING METHODS ====================

// TrackPageView records a page view
func (s *AnalyticsService) TrackPageView(userID *string, payload dto.TrackPageViewDto, ipAddress, userAgent string) error {
	var uid *uuid.UUID
	if userID != nil && *userID != "" {
		parsed, err := uuid.Parse(*userID)
		if err == nil {
			uid = &parsed
		}
	}

	deviceType := detectDeviceType(userAgent)
	browser := detectBrowser(userAgent)
	os := detectOS(userAgent)

	pageView := &models.PageView{
		Path:       payload.Path,
		UserID:     uid,
		SessionID:  payload.SessionID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Referrer:   payload.Referrer,
		DeviceType: &deviceType,
		Browser:    &browser,
		OS:         &os,
		Duration:   payload.Duration,
	}

	return s.database.Create(pageView).Error
}

// TrackJobView records a job view
func (s *AnalyticsService) TrackJobView(userID *string, payload dto.TrackJobViewDto, ipAddress, userAgent string) error {
	jobID, err := uuid.Parse(payload.JobID)
	if err != nil {
		return errors.New("invalid job ID")
	}

	var uid *uuid.UUID
	if userID != nil && *userID != "" {
		parsed, err := uuid.Parse(*userID)
		if err == nil {
			uid = &parsed
		}
	}

	jobView := &models.JobView{
		JobID:     jobID,
		UserID:    uid,
		SessionID: payload.SessionID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Referrer:  payload.Referrer,
	}

	return s.database.Create(jobView).Error
}

// TrackProfileView records a profile view
func (s *AnalyticsService) TrackProfileView(viewerUserID *string, payload dto.TrackProfileViewDto, ipAddress, userAgent string, isRecruiter bool) error {
	profileUserID, err := uuid.Parse(payload.ProfileUserID)
	if err != nil {
		return errors.New("invalid profile user ID")
	}

	var viewerUID *uuid.UUID
	if viewerUserID != nil && *viewerUserID != "" {
		parsed, err := uuid.Parse(*viewerUserID)
		if err == nil {
			viewerUID = &parsed
		}
	}

	profileView := &models.ProfileView{
		ProfileUserID:     profileUserID,
		ViewerUserID:      viewerUID,
		SessionID:         payload.SessionID,
		IPAddress:         ipAddress,
		UserAgent:         userAgent,
		Referrer:          payload.Referrer,
		ViewerIsRecruiter: isRecruiter,
	}

	return s.database.Create(profileView).Error
}

// TrackPortfolioView records a portfolio view
func (s *AnalyticsService) TrackPortfolioView(viewerUserID *string, payload dto.TrackPortfolioViewDto, ipAddress, userAgent string) error {
	portfolioID, err := uuid.Parse(payload.PortfolioID)
	if err != nil {
		return errors.New("invalid portfolio ID")
	}

	var viewerUID *uuid.UUID
	if viewerUserID != nil && *viewerUserID != "" {
		parsed, err := uuid.Parse(*viewerUserID)
		if err == nil {
			viewerUID = &parsed
		}
	}

	deviceType := detectDeviceType(userAgent)

	portfolioView := &models.PortfolioView{
		PortfolioID:  portfolioID,
		ViewerUserID: viewerUID,
		SessionID:    payload.SessionID,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Referrer:     payload.Referrer,
		DeviceType:   &deviceType,
		Duration:     payload.Duration,
	}

	return s.database.Create(portfolioView).Error
}

// TrackEvent records a custom event
func (s *AnalyticsService) TrackEvent(userID *string, payload dto.TrackEventDto) error {
	var uid *uuid.UUID
	if userID != nil && *userID != "" {
		parsed, err := uuid.Parse(*userID)
		if err == nil {
			uid = &parsed
		}
	}

	var entityID *uuid.UUID
	if payload.EntityID != nil {
		parsed, err := uuid.Parse(*payload.EntityID)
		if err == nil {
			entityID = &parsed
		}
	}

	event := &models.AnalyticsEvent{
		EventType:  models.EventType(payload.EventType),
		UserID:     uid,
		EntityID:   entityID,
		EntityType: payload.EntityType,
		SessionID:  payload.SessionID,
		Properties: payload.Properties,
	}

	return s.database.Create(event).Error
}

// ==================== ADMIN ANALYTICS ====================

// GetAdminDashboardAnalytics returns comprehensive platform analytics for admins
func (s *AnalyticsService) GetAdminDashboardAnalytics(params dto.AnalyticsQueryParams) (*dto.AdminDashboardAnalytics, error) {
	startDate, endDate := s.parseDateRange(params)

	overview := s.getPlatformOverview()
	userStats := s.getUserStats(startDate, endDate)
	jobStats := s.getJobStats(startDate, endDate)
	appStats := s.getApplicationStats(startDate, endDate)
	revenueStats := s.getRevenueStats(startDate, endDate)
	trendData := s.getTrendData(startDate, endDate, params.GroupBy, "platform")
	topPerformers := s.getTopPerformers(startDate, endDate)

	return &dto.AdminDashboardAnalytics{
		Overview:         overview,
		UserStats:        userStats,
		JobStats:         jobStats,
		ApplicationStats: appStats,
		RevenueStats:     revenueStats,
		TrendData:        trendData,
		TopPerformers:    topPerformers,
	}, nil
}

func (s *AnalyticsService) getPlatformOverview() dto.PlatformOverview {
	var totalUsers, totalRecruiters, totalTalents, totalJobs, totalApplications int64
	var totalPageViews, uniqueVisitors, activeSubscriptions int64

	s.database.Model(&models.User{}).Count(&totalUsers)
	s.database.Model(&models.User{}).Where("is_recruiter = ?", true).Count(&totalRecruiters)
	s.database.Model(&models.User{}).Where("is_recruiter = ?", false).Count(&totalTalents)
	s.database.Model(&models.Job{}).Count(&totalJobs)
	s.database.Model(&models.JobApplication{}).Count(&totalApplications)
	s.database.Model(&models.PageView{}).Count(&totalPageViews)
	s.database.Model(&models.PageView{}).Distinct("session_id").Count(&uniqueVisitors)
	s.database.Model(&models.UserSubscription{}).Where("status = ?", "active").Count(&activeSubscriptions)

	// Calculate growth rate (comparing this month to last month)
	var usersThisMonth, usersLastMonth int64
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	startOfLastMonth := startOfMonth.AddDate(0, -1, 0)

	s.database.Model(&models.User{}).Where("created_at >= ?", startOfMonth).Count(&usersThisMonth)
	s.database.Model(&models.User{}).Where("created_at >= ? AND created_at < ?", startOfLastMonth, startOfMonth).Count(&usersLastMonth)

	var growthRate float64
	if usersLastMonth > 0 {
		growthRate = float64(usersThisMonth-usersLastMonth) / float64(usersLastMonth) * 100
	}

	return dto.PlatformOverview{
		TotalUsers:          totalUsers,
		TotalRecruiters:     totalRecruiters,
		TotalTalents:        totalTalents,
		TotalJobs:           totalJobs,
		TotalApplications:   totalApplications,
		TotalPageViews:      totalPageViews,
		UniqueVisitors:      uniqueVisitors,
		ActiveSubscriptions: activeSubscriptions,
		GrowthRate:          growthRate,
	}
}

func (s *AnalyticsService) getUserStats(startDate, endDate time.Time) dto.UserStatsResponse {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	weekAgo := today.AddDate(0, 0, -7)
	monthAgo := today.AddDate(0, -1, 0)

	var newToday, newWeek, newMonth int64
	s.database.Model(&models.User{}).Where("created_at >= ?", today).Count(&newToday)
	s.database.Model(&models.User{}).Where("created_at >= ?", weekAgo).Count(&newWeek)
	s.database.Model(&models.User{}).Where("created_at >= ?", monthAgo).Count(&newMonth)

	var providers []dto.ProviderCount
	s.database.Model(&models.User{}).
		Select("provider, count(*) as count").
		Group("provider").
		Scan(&providers)

	var locations []dto.LocationCount
	s.database.Model(&models.User{}).
		Select("location, count(*) as count").
		Where("location IS NOT NULL AND location != ''").
		Group("location").
		Order("count DESC").
		Limit(10).
		Scan(&locations)

	var totalUsers, verifiedUsers int64
	s.database.Model(&models.User{}).Count(&totalUsers)
	s.database.Model(&models.User{}).Where("verified = ?", true).Count(&verifiedUsers)

	var verificationRate float64
	if totalUsers > 0 {
		verificationRate = float64(verifiedUsers) / float64(totalUsers) * 100
	}

	return dto.UserStatsResponse{
		NewUsersToday:     newToday,
		NewUsersThisWeek:  newWeek,
		NewUsersThisMonth: newMonth,
		UsersByProvider:   providers,
		UsersByLocation:   locations,
		VerificationRate:  verificationRate,
	}
}

func (s *AnalyticsService) getJobStats(startDate, endDate time.Time) dto.JobStatsResponse {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	weekAgo := today.AddDate(0, 0, -7)
	monthAgo := today.AddDate(0, -1, 0)

	var activeJobs, newToday, newWeek, newMonth int64
	s.database.Model(&models.Job{}).Count(&activeJobs)
	s.database.Model(&models.Job{}).Where("created_at >= ?", today).Count(&newToday)
	s.database.Model(&models.Job{}).Where("created_at >= ?", weekAgo).Count(&newWeek)
	s.database.Model(&models.Job{}).Where("created_at >= ?", monthAgo).Count(&newMonth)

	var jobsByType []dto.EmploymentTypeCount
	s.database.Model(&models.Job{}).
		Select("employment_type as type, count(*) as count").
		Group("employment_type").
		Scan(&jobsByType)

	var jobsByLocation []dto.LocationCount
	s.database.Model(&models.Job{}).
		Select("location, count(*) as count").
		Group("location").
		Order("count DESC").
		Limit(10).
		Scan(&jobsByLocation)

	// Average applications per job
	var totalApps int64
	s.database.Model(&models.JobApplication{}).Count(&totalApps)
	var avgApps float64
	if activeJobs > 0 {
		avgApps = float64(totalApps) / float64(activeJobs)
	}

	// Most viewed jobs
	var mostViewed []dto.JobViewCount
	s.database.Model(&models.JobView{}).
		Select("job_views.job_id, jobs.title, count(*) as views").
		Joins("JOIN jobs ON jobs.id = job_views.job_id").
		Where("job_views.created_at BETWEEN ? AND ?", startDate, endDate).
		Group("job_views.job_id, jobs.title").
		Order("views DESC").
		Limit(10).
		Scan(&mostViewed)

	return dto.JobStatsResponse{
		TotalActiveJobs:     activeJobs,
		NewJobsToday:        newToday,
		NewJobsThisWeek:     newWeek,
		NewJobsThisMonth:    newMonth,
		JobsByType:          jobsByType,
		JobsByLocation:      jobsByLocation,
		AverageApplications: avgApps,
		MostViewedJobs:      mostViewed,
	}
}

func (s *AnalyticsService) getApplicationStats(startDate, endDate time.Time) dto.ApplicationStatsResponse {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	weekAgo := today.AddDate(0, 0, -7)
	monthAgo := today.AddDate(0, -1, 0)

	var total, appsToday, appsWeek, appsMonth int64
	s.database.Model(&models.JobApplication{}).Count(&total)
	s.database.Model(&models.JobApplication{}).Where("created_at >= ?", today).Count(&appsToday)
	s.database.Model(&models.JobApplication{}).Where("created_at >= ?", weekAgo).Count(&appsWeek)
	s.database.Model(&models.JobApplication{}).Where("created_at >= ?", monthAgo).Count(&appsMonth)

	var byStatus []dto.StatusCount
	s.database.Model(&models.JobApplication{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&byStatus)

	var accepted, hired int64
	s.database.Model(&models.JobApplication{}).Where("status = ?", "ACCEPTED").Count(&accepted)
	s.database.Model(&models.JobApplication{}).Where("status = ?", "HIRED").Count(&hired)

	var acceptanceRate, hireRate float64
	if total > 0 {
		acceptanceRate = float64(accepted) / float64(total) * 100
		hireRate = float64(hired) / float64(total) * 100
	}

	return dto.ApplicationStatsResponse{
		TotalApplications:     total,
		ApplicationsToday:     appsToday,
		ApplicationsThisWeek:  appsWeek,
		ApplicationsThisMonth: appsMonth,
		ApplicationsByStatus:  byStatus,
		AcceptanceRate:        acceptanceRate,
		HireRate:              hireRate,
	}
}

func (s *AnalyticsService) getRevenueStats(startDate, endDate time.Time) dto.RevenueStatsResponse {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	startOfLastMonth := startOfMonth.AddDate(0, -1, 0)

	var totalRevenue, revenueThisMonth, revenueLastMonth float64

	s.database.Model(&models.SubscriptionInvoice{}).
		Where("status = ?", "paid").
		Select("COALESCE(SUM(amount_paid), 0)").
		Scan(&totalRevenue)

	s.database.Model(&models.SubscriptionInvoice{}).
		Where("status = ? AND paid_at >= ?", "paid", startOfMonth).
		Select("COALESCE(SUM(amount_paid), 0)").
		Scan(&revenueThisMonth)

	s.database.Model(&models.SubscriptionInvoice{}).
		Where("status = ? AND paid_at >= ? AND paid_at < ?", "paid", startOfLastMonth, startOfMonth).
		Select("COALESCE(SUM(amount_paid), 0)").
		Scan(&revenueLastMonth)

	var monthlyGrowth float64
	if revenueLastMonth > 0 {
		monthlyGrowth = (revenueThisMonth - revenueLastMonth) / revenueLastMonth * 100
	}

	var tierCounts []dto.TierCount
	s.database.Model(&models.UserSubscription{}).
		Select("subscriptions.tier, count(*) as count").
		Joins("JOIN subscriptions ON subscriptions.id = user_subscriptions.subscription_id").
		Where("user_subscriptions.status = ?", "active").
		Group("subscriptions.tier").
		Scan(&tierCounts)

	var activeUsers int64
	s.database.Model(&models.UserSubscription{}).Where("status = ?", "active").Count(&activeUsers)
	var arpu float64
	if activeUsers > 0 {
		arpu = totalRevenue / float64(activeUsers)
	}

	return dto.RevenueStatsResponse{
		TotalRevenue:          totalRevenue,
		RevenueThisMonth:      revenueThisMonth,
		RevenueLastMonth:      revenueLastMonth,
		MonthlyGrowth:         monthlyGrowth,
		SubscriptionsByTier:   tierCounts,
		AverageRevenuePerUser: arpu,
	}
}

func (s *AnalyticsService) getTopPerformers(startDate, endDate time.Time) dto.TopPerformersResponse {
	var topRecruiters []dto.RecruiterPerformance
	s.database.Raw(`
		SELECT u.id as user_id, u.name,
			COUNT(DISTINCT j.id) as total_jobs,
			COUNT(DISTINCT CASE WHEN ja.status = 'HIRED' THEN ja.id END) as total_hires
		FROM users u
		LEFT JOIN jobs j ON j.created_by = u.id
		LEFT JOIN job_applications ja ON ja.job_id = j.id
		WHERE u.is_recruiter = true
		GROUP BY u.id, u.name
		ORDER BY total_hires DESC, total_jobs DESC
		LIMIT 10
	`).Scan(&topRecruiters)

	var topJobs []dto.JobPerformance
	s.database.Raw(`
		SELECT j.id as job_id, j.title, j.company,
			COUNT(DISTINCT jv.id) as views,
			COUNT(DISTINCT ja.id) as applications
		FROM jobs j
		LEFT JOIN job_views jv ON jv.job_id = j.id
		LEFT JOIN job_applications ja ON ja.job_id = j.id
		GROUP BY j.id, j.title, j.company
		ORDER BY views DESC
		LIMIT 10
	`).Scan(&topJobs)

	var topTalents []dto.TalentPerformance
	s.database.Raw(`
		SELECT u.id as user_id, u.name,
			COUNT(DISTINCT pv.id) as profile_views,
			COUNT(DISTINCT ja.id) as applications
		FROM users u
		LEFT JOIN profile_views pv ON pv.profile_user_id = u.id
		LEFT JOIN job_applications ja ON ja.applicant_id = u.id
		WHERE u.is_recruiter = false
		GROUP BY u.id, u.name
		ORDER BY profile_views DESC
		LIMIT 10
	`).Scan(&topTalents)

	return dto.TopPerformersResponse{
		TopRecruiters: topRecruiters,
		TopJobs:       topJobs,
		TopTalents:    topTalents,
	}
}

// ==================== RECRUITER ANALYTICS ====================

// GetRecruiterAnalytics returns analytics for a recruiter's jobs and applications
func (s *AnalyticsService) GetRecruiterAnalytics(userID string, params dto.AnalyticsQueryParams) (*dto.RecruiterAnalytics, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	startDate, endDate := s.parseDateRange(params)

	overview := s.getRecruiterOverview(uid)
	jobPerformance := s.getRecruiterJobPerformance(uid, startDate, endDate)
	appStats := s.getRecruiterApplicationStats(uid)
	trendData := s.getRecruiterTrendData(uid, startDate, endDate, params.GroupBy)
	topApplicants := s.getTopApplicants(uid, 10)

	return &dto.RecruiterAnalytics{
		Overview:         overview,
		JobPerformance:   jobPerformance,
		ApplicationStats: appStats,
		TrendData:        trendData,
		TopApplicants:    topApplicants,
	}, nil
}

func (s *AnalyticsService) getRecruiterOverview(userID uuid.UUID) dto.RecruiterOverview {
	var totalJobs, activeJobs, totalApps, pendingApps, totalJobViews, totalHires int64

	s.database.Model(&models.Job{}).Where("created_by = ?", userID).Count(&totalJobs)
	s.database.Model(&models.Job{}).Where("created_by = ? AND (deadline IS NULL OR deadline > ?)", userID, time.Now()).Count(&activeJobs)

	s.database.Model(&models.JobApplication{}).
		Joins("JOIN jobs ON jobs.id = job_applications.job_id").
		Where("jobs.created_by = ?", userID).
		Count(&totalApps)

	s.database.Model(&models.JobApplication{}).
		Joins("JOIN jobs ON jobs.id = job_applications.job_id").
		Where("jobs.created_by = ? AND job_applications.status = ?", userID, "PENDING").
		Count(&pendingApps)

	s.database.Model(&models.JobView{}).
		Joins("JOIN jobs ON jobs.id = job_views.job_id").
		Where("jobs.created_by = ?", userID).
		Count(&totalJobViews)

	s.database.Model(&models.JobApplication{}).
		Joins("JOIN jobs ON jobs.id = job_applications.job_id").
		Where("jobs.created_by = ? AND job_applications.status = ?", userID, "HIRED").
		Count(&totalHires)

	var responseRate float64
	if totalApps > 0 {
		var responded int64
		s.database.Model(&models.JobApplication{}).
			Joins("JOIN jobs ON jobs.id = job_applications.job_id").
			Where("jobs.created_by = ? AND job_applications.status != ?", userID, "PENDING").
			Count(&responded)
		responseRate = float64(responded) / float64(totalApps) * 100
	}

	return dto.RecruiterOverview{
		TotalJobs:           totalJobs,
		ActiveJobs:          activeJobs,
		TotalApplications:   totalApps,
		PendingApplications: pendingApps,
		TotalJobViews:       totalJobViews,
		TotalHires:          totalHires,
		ResponseRate:        responseRate,
	}
}

func (s *AnalyticsService) getRecruiterJobPerformance(userID uuid.UUID, startDate, endDate time.Time) []dto.JobPerformanceDetail {
	var performance []dto.JobPerformanceDetail

	s.database.Raw(`
		SELECT j.id as job_id, j.title,
			COUNT(DISTINCT jv.id) as views,
			COUNT(DISTINCT ja.id) as applications,
			j.posted_date,
			CASE WHEN j.deadline IS NULL OR j.deadline > NOW() THEN 'active' ELSE 'closed' END as status
		FROM jobs j
		LEFT JOIN job_views jv ON jv.job_id = j.id AND jv.created_at BETWEEN ? AND ?
		LEFT JOIN job_applications ja ON ja.job_id = j.id
		WHERE j.created_by = ?
		GROUP BY j.id, j.title, j.posted_date, j.deadline
		ORDER BY j.posted_date DESC
	`, startDate, endDate, userID).Scan(&performance)

	// Calculate conversion rates
	for i := range performance {
		if performance[i].Views > 0 {
			performance[i].ConversionRate = float64(performance[i].Applications) / float64(performance[i].Views) * 100
		}
	}

	return performance
}

func (s *AnalyticsService) getRecruiterApplicationStats(userID uuid.UUID) dto.RecruiterApplicationStats {
	var total, pending, reviewed, accepted, rejected, hired int64

	baseQuery := s.database.Model(&models.JobApplication{}).
		Joins("JOIN jobs ON jobs.id = job_applications.job_id").
		Where("jobs.created_by = ?", userID)

	baseQuery.Count(&total)
	baseQuery.Where("job_applications.status = ?", "PENDING").Count(&pending)

	s.database.Model(&models.JobApplication{}).
		Joins("JOIN jobs ON jobs.id = job_applications.job_id").
		Where("jobs.created_by = ? AND job_applications.status = ?", userID, "REVIEWED").Count(&reviewed)

	s.database.Model(&models.JobApplication{}).
		Joins("JOIN jobs ON jobs.id = job_applications.job_id").
		Where("jobs.created_by = ? AND job_applications.status = ?", userID, "ACCEPTED").Count(&accepted)

	s.database.Model(&models.JobApplication{}).
		Joins("JOIN jobs ON jobs.id = job_applications.job_id").
		Where("jobs.created_by = ? AND job_applications.status = ?", userID, "REJECTED").Count(&rejected)

	s.database.Model(&models.JobApplication{}).
		Joins("JOIN jobs ON jobs.id = job_applications.job_id").
		Where("jobs.created_by = ? AND job_applications.status = ?", userID, "HIRED").Count(&hired)

	byStatus := []dto.StatusCount{
		{Status: "PENDING", Count: pending},
		{Status: "REVIEWED", Count: reviewed},
		{Status: "ACCEPTED", Count: accepted},
		{Status: "REJECTED", Count: rejected},
		{Status: "HIRED", Count: hired},
	}

	return dto.RecruiterApplicationStats{
		TotalReceived: total,
		Pending:       pending,
		Reviewed:      reviewed,
		Accepted:      accepted,
		Rejected:      rejected,
		Hired:         hired,
		ByStatus:      byStatus,
	}
}

func (s *AnalyticsService) getRecruiterTrendData(userID uuid.UUID, startDate, endDate time.Time, groupBy string) []dto.TrendDataPoint {
	var trendData []dto.TrendDataPoint

	dateFormat := "YYYY-MM-DD"
	if groupBy == "week" {
		dateFormat = "IYYY-IW"
	} else if groupBy == "month" {
		dateFormat = "YYYY-MM"
	}

	s.database.Raw(`
		SELECT TO_CHAR(ja.created_at, ?) as date, COUNT(*) as value
		FROM job_applications ja
		JOIN jobs j ON j.id = ja.job_id
		WHERE j.created_by = ? AND ja.created_at BETWEEN ? AND ?
		GROUP BY TO_CHAR(ja.created_at, ?)
		ORDER BY date
	`, dateFormat, userID, startDate, endDate, dateFormat).Scan(&trendData)

	return trendData
}

func (s *AnalyticsService) getTopApplicants(userID uuid.UUID, limit int) []dto.ApplicantInfo {
	var applicants []dto.ApplicantInfo

	s.database.Raw(`
		SELECT u.id as user_id, u.name, u.email, ja.submission_date as applied_date,
			j.title as job_title, ja.status
		FROM job_applications ja
		JOIN users u ON u.id = ja.applicant_id
		JOIN jobs j ON j.id = ja.job_id
		WHERE j.created_by = ?
		ORDER BY ja.submission_date DESC
		LIMIT ?
	`, userID, limit).Scan(&applicants)

	return applicants
}

// ==================== TALENT ANALYTICS ====================

// GetTalentAnalytics returns analytics for a talent's profile and applications
func (s *AnalyticsService) GetTalentAnalytics(userID string, params dto.AnalyticsQueryParams) (*dto.TalentAnalytics, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	startDate, endDate := s.parseDateRange(params)

	overview := s.getTalentOverview(uid)
	profileViews := s.getTalentProfileViews(uid, startDate, endDate)
	portfolioStats := s.getTalentPortfolioStats(uid, startDate, endDate)
	appStats := s.getTalentApplicationStats(uid)
	trendData := s.getTalentTrendData(uid, startDate, endDate, params.GroupBy)
	viewerInsights := s.getViewerInsights(uid, startDate, endDate)

	return &dto.TalentAnalytics{
		Overview:         overview,
		ProfileViews:     profileViews,
		PortfolioStats:   portfolioStats,
		ApplicationStats: appStats,
		TrendData:        trendData,
		ViewerInsights:   viewerInsights,
	}, nil
}

func (s *AnalyticsService) getTalentOverview(userID uuid.UUID) dto.TalentOverview {
	var profileViews, portfolioViews, totalApps, pendingApps, acceptedApps, recruiterViews int64

	s.database.Model(&models.ProfileView{}).Where("profile_user_id = ?", userID).Count(&profileViews)
	s.database.Model(&models.PortfolioView{}).
		Joins("JOIN portfolios ON portfolios.id = portfolio_views.portfolio_id").
		Where("portfolios.user_id = ?", userID).Count(&portfolioViews)

	s.database.Model(&models.JobApplication{}).Where("applicant_id = ?", userID).Count(&totalApps)
	s.database.Model(&models.JobApplication{}).Where("applicant_id = ? AND status = ?", userID, "PENDING").Count(&pendingApps)
	s.database.Model(&models.JobApplication{}).Where("applicant_id = ? AND status = ?", userID, "ACCEPTED").Count(&acceptedApps)
	s.database.Model(&models.ProfileView{}).Where("profile_user_id = ? AND viewer_is_recruiter = ?", userID, true).Count(&recruiterViews)

	// Calculate profile completeness
	var user models.User
	s.database.Preload("Projects").Preload("Experiences").Preload("Education").
		Preload("Skills").Preload("Certifications").First(&user, "id = ?", userID)

	completeness := calculateProfileCompleteness(&user)

	return dto.TalentOverview{
		TotalProfileViews:    profileViews,
		TotalPortfolioViews:  portfolioViews,
		TotalApplications:    totalApps,
		PendingApplications:  pendingApps,
		AcceptedApplications: acceptedApps,
		RecruiterViews:       recruiterViews,
		ProfileCompleteness:  completeness,
	}
}

func (s *AnalyticsService) getTalentProfileViews(userID uuid.UUID, startDate, endDate time.Time) dto.ProfileViewsStats {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	weekAgo := today.AddDate(0, 0, -7)
	monthAgo := today.AddDate(0, -1, 0)

	var total, viewsToday, viewsWeek, viewsMonth, uniqueViewers, recruiterViewers int64

	s.database.Model(&models.ProfileView{}).Where("profile_user_id = ?", userID).Count(&total)
	s.database.Model(&models.ProfileView{}).Where("profile_user_id = ? AND created_at >= ?", userID, today).Count(&viewsToday)
	s.database.Model(&models.ProfileView{}).Where("profile_user_id = ? AND created_at >= ?", userID, weekAgo).Count(&viewsWeek)
	s.database.Model(&models.ProfileView{}).Where("profile_user_id = ? AND created_at >= ?", userID, monthAgo).Count(&viewsMonth)
	s.database.Model(&models.ProfileView{}).Where("profile_user_id = ?", userID).Distinct("session_id").Count(&uniqueViewers)
	s.database.Model(&models.ProfileView{}).Where("profile_user_id = ? AND viewer_is_recruiter = ?", userID, true).Count(&recruiterViewers)

	var byCountry []dto.LocationCount
	s.database.Model(&models.ProfileView{}).
		Select("country as location, count(*) as count").
		Where("profile_user_id = ? AND country IS NOT NULL", userID).
		Group("country").
		Order("count DESC").
		Limit(10).
		Scan(&byCountry)

	return dto.ProfileViewsStats{
		TotalViews:       total,
		ViewsToday:       viewsToday,
		ViewsThisWeek:    viewsWeek,
		ViewsThisMonth:   viewsMonth,
		UniqueViewers:    uniqueViewers,
		RecruiterViewers: recruiterViewers,
		ViewsByCountry:   byCountry,
	}
}

func (s *AnalyticsService) getTalentPortfolioStats(userID uuid.UUID, startDate, endDate time.Time) dto.PortfolioStatsResponse {
	var portfolio models.Portfolio
	if err := s.database.Where("user_id = ?", userID).First(&portfolio).Error; err != nil {
		return dto.PortfolioStatsResponse{}
	}

	var total, uniqueVisitors int64
	s.database.Model(&models.PortfolioView{}).Where("portfolio_id = ?", portfolio.ID).Count(&total)
	s.database.Model(&models.PortfolioView{}).Where("portfolio_id = ?", portfolio.ID).Distinct("session_id").Count(&uniqueVisitors)

	var avgDuration float64
	s.database.Model(&models.PortfolioView{}).
		Where("portfolio_id = ? AND duration IS NOT NULL", portfolio.ID).
		Select("COALESCE(AVG(duration), 0)").
		Scan(&avgDuration)

	var topReferrers []dto.ReferrerCount
	s.database.Model(&models.PortfolioView{}).
		Select("referrer, count(*) as count").
		Where("portfolio_id = ? AND referrer IS NOT NULL", portfolio.ID).
		Group("referrer").
		Order("count DESC").
		Limit(10).
		Scan(&topReferrers)

	return dto.PortfolioStatsResponse{
		TotalViews:      total,
		UniqueVisitors:  uniqueVisitors,
		AverageDuration: avgDuration,
		TopReferrers:    topReferrers,
	}
}

func (s *AnalyticsService) getTalentApplicationStats(userID uuid.UUID) dto.TalentApplicationStats {
	var total, pending, reviewed, accepted, rejected, hired int64

	s.database.Model(&models.JobApplication{}).Where("applicant_id = ?", userID).Count(&total)
	s.database.Model(&models.JobApplication{}).Where("applicant_id = ? AND status = ?", userID, "PENDING").Count(&pending)
	s.database.Model(&models.JobApplication{}).Where("applicant_id = ? AND status = ?", userID, "REVIEWED").Count(&reviewed)
	s.database.Model(&models.JobApplication{}).Where("applicant_id = ? AND status = ?", userID, "ACCEPTED").Count(&accepted)
	s.database.Model(&models.JobApplication{}).Where("applicant_id = ? AND status = ?", userID, "REJECTED").Count(&rejected)
	s.database.Model(&models.JobApplication{}).Where("applicant_id = ? AND status = ?", userID, "HIRED").Count(&hired)

	var responseRate float64
	if total > 0 {
		responded := reviewed + accepted + rejected + hired
		responseRate = float64(responded) / float64(total) * 100
	}

	byStatus := []dto.StatusCount{
		{Status: "PENDING", Count: pending},
		{Status: "REVIEWED", Count: reviewed},
		{Status: "ACCEPTED", Count: accepted},
		{Status: "REJECTED", Count: rejected},
		{Status: "HIRED", Count: hired},
	}

	return dto.TalentApplicationStats{
		TotalSent:    total,
		Pending:      pending,
		Reviewed:     reviewed,
		Accepted:     accepted,
		Rejected:     rejected,
		Hired:        hired,
		ResponseRate: responseRate,
		ByStatus:     byStatus,
	}
}

func (s *AnalyticsService) getTalentTrendData(userID uuid.UUID, startDate, endDate time.Time, groupBy string) []dto.TrendDataPoint {
	var trendData []dto.TrendDataPoint

	dateFormat := "YYYY-MM-DD"
	if groupBy == "week" {
		dateFormat = "IYYY-IW"
	} else if groupBy == "month" {
		dateFormat = "YYYY-MM"
	}

	s.database.Raw(`
		SELECT TO_CHAR(created_at, ?) as date, COUNT(*) as value
		FROM profile_views
		WHERE profile_user_id = ? AND created_at BETWEEN ? AND ?
		GROUP BY TO_CHAR(created_at, ?)
		ORDER BY date
	`, dateFormat, userID, startDate, endDate, dateFormat).Scan(&trendData)

	return trendData
}

func (s *AnalyticsService) getViewerInsights(userID uuid.UUID, startDate, endDate time.Time) dto.ViewerInsightsResponse {
	var recruiterViews, talentViews, anonymousViews int64

	s.database.Model(&models.ProfileView{}).
		Where("profile_user_id = ? AND viewer_is_recruiter = ? AND created_at BETWEEN ? AND ?", userID, true, startDate, endDate).
		Count(&recruiterViews)

	s.database.Model(&models.ProfileView{}).
		Where("profile_user_id = ? AND viewer_is_recruiter = ? AND viewer_user_id IS NOT NULL AND created_at BETWEEN ? AND ?", userID, false, startDate, endDate).
		Count(&talentViews)

	s.database.Model(&models.ProfileView{}).
		Where("profile_user_id = ? AND viewer_user_id IS NULL AND created_at BETWEEN ? AND ?", userID, startDate, endDate).
		Count(&anonymousViews)

	return dto.ViewerInsightsResponse{
		RecruiterViews: recruiterViews,
		TalentViews:    talentViews,
		AnonymousViews: anonymousViews,
	}
}

// ==================== HELPER METHODS ====================

func (s *AnalyticsService) getTrendData(startDate, endDate time.Time, groupBy, statType string) []dto.TrendDataPoint {
	var trendData []dto.TrendDataPoint

	dateFormat := "YYYY-MM-DD"
	if groupBy == "week" {
		dateFormat = "IYYY-IW"
	} else if groupBy == "month" {
		dateFormat = "YYYY-MM"
	}

	s.database.Raw(`
		SELECT TO_CHAR(created_at, ?) as date, COUNT(*) as value
		FROM page_views
		WHERE created_at BETWEEN ? AND ?
		GROUP BY TO_CHAR(created_at, ?)
		ORDER BY date
	`, dateFormat, startDate, endDate, dateFormat).Scan(&trendData)

	return trendData
}

func (s *AnalyticsService) parseDateRange(params dto.AnalyticsQueryParams) (time.Time, time.Time) {
	now := time.Now()
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)
	startDate := endDate.AddDate(0, -1, 0) // Default: last 30 days

	if params.StartDate != "" {
		if parsed, err := time.Parse("2006-01-02", params.StartDate); err == nil {
			startDate = parsed
		}
	}

	if params.EndDate != "" {
		if parsed, err := time.Parse("2006-01-02", params.EndDate); err == nil {
			endDate = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 23, 59, 59, 0, time.UTC)
		}
	}

	return startDate, endDate
}

func calculateProfileCompleteness(user *models.User) float64 {
	var score float64
	total := 10.0

	if user.Name != "" {
		score++
	}
	if user.Headline != nil && *user.Headline != "" {
		score++
	}
	if user.Summary != nil && *user.Summary != "" {
		score++
	}
	if user.Location != nil && *user.Location != "" {
		score++
	}
	if user.Image != nil && *user.Image != "" {
		score++
	}
	if len(user.Skills) > 0 {
		score++
	}
	if len(user.Projects) > 0 {
		score++
	}
	if len(user.Experiences) > 0 {
		score++
	}
	if len(user.Education) > 0 {
		score++
	}
	if user.SocialMedia != nil {
		score++
	}

	return (score / total) * 100
}

func detectDeviceType(userAgent string) string {
	ua := strings.ToLower(userAgent)
	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") || strings.Contains(ua, "iphone") {
		return "mobile"
	}
	if strings.Contains(ua, "tablet") || strings.Contains(ua, "ipad") {
		return "tablet"
	}
	return "desktop"
}

func detectBrowser(userAgent string) string {
	ua := strings.ToLower(userAgent)
	if strings.Contains(ua, "chrome") && !strings.Contains(ua, "edge") {
		return "Chrome"
	}
	if strings.Contains(ua, "firefox") {
		return "Firefox"
	}
	if strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome") {
		return "Safari"
	}
	if strings.Contains(ua, "edge") {
		return "Edge"
	}
	if strings.Contains(ua, "opera") {
		return "Opera"
	}
	return "Other"
}

func detectOS(userAgent string) string {
	ua := strings.ToLower(userAgent)
	if strings.Contains(ua, "windows") {
		return "Windows"
	}
	if strings.Contains(ua, "mac os") || strings.Contains(ua, "macos") {
		return "macOS"
	}
	if strings.Contains(ua, "linux") {
		return "Linux"
	}
	if strings.Contains(ua, "android") {
		return "Android"
	}
	if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") || strings.Contains(ua, "ios") {
		return "iOS"
	}
	return "Other"
}
