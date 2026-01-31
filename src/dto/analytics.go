package dto

import "time"

// Request DTOs

type TrackPageViewDto struct {
	Path      string  `json:"path" binding:"required"`
	SessionID string  `json:"session_id" binding:"required"`
	Referrer  *string `json:"referrer,omitempty"`
	Duration  *int    `json:"duration,omitempty"`
}

type TrackJobViewDto struct {
	JobID     string  `json:"job_id" binding:"required,uuid"`
	SessionID string  `json:"session_id" binding:"required"`
	Referrer  *string `json:"referrer,omitempty"`
}

type TrackProfileViewDto struct {
	ProfileUserID string  `json:"profile_user_id" binding:"required,uuid"`
	SessionID     string  `json:"session_id" binding:"required"`
	Referrer      *string `json:"referrer,omitempty"`
}

type TrackPortfolioViewDto struct {
	PortfolioID string  `json:"portfolio_id" binding:"required,uuid"`
	SessionID   string  `json:"session_id" binding:"required"`
	Referrer    *string `json:"referrer,omitempty"`
	Duration    *int    `json:"duration,omitempty"`
}

type TrackEventDto struct {
	EventType  string                 `json:"event_type" binding:"required"`
	EntityID   *string                `json:"entity_id,omitempty"`
	EntityType *string                `json:"entity_type,omitempty"`
	SessionID  string                 `json:"session_id" binding:"required"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type AnalyticsQueryParams struct {
	StartDate string `form:"start_date"` // YYYY-MM-DD
	EndDate   string `form:"end_date"`   // YYYY-MM-DD
	GroupBy   string `form:"group_by"`   // day, week, month
	Limit     int    `form:"limit"`
}

// Response DTOs

// Admin Dashboard Analytics
type AdminDashboardAnalytics struct {
	Overview        PlatformOverview        `json:"overview"`
	UserStats       UserStatsResponse       `json:"user_stats"`
	JobStats        JobStatsResponse        `json:"job_stats"`
	ApplicationStats ApplicationStatsResponse `json:"application_stats"`
	RevenueStats    RevenueStatsResponse    `json:"revenue_stats"`
	TrendData       []TrendDataPoint        `json:"trend_data"`
	TopPerformers   TopPerformersResponse   `json:"top_performers"`
}

type PlatformOverview struct {
	TotalUsers          int64   `json:"total_users"`
	TotalRecruiters     int64   `json:"total_recruiters"`
	TotalTalents        int64   `json:"total_talents"`
	TotalJobs           int64   `json:"total_jobs"`
	TotalApplications   int64   `json:"total_applications"`
	TotalPageViews      int64   `json:"total_page_views"`
	UniqueVisitors      int64   `json:"unique_visitors"`
	ActiveSubscriptions int64   `json:"active_subscriptions"`
	GrowthRate          float64 `json:"growth_rate"` // percentage
}

type UserStatsResponse struct {
	NewUsersToday     int64            `json:"new_users_today"`
	NewUsersThisWeek  int64            `json:"new_users_this_week"`
	NewUsersThisMonth int64            `json:"new_users_this_month"`
	UsersByProvider   []ProviderCount  `json:"users_by_provider"`
	UsersByLocation   []LocationCount  `json:"users_by_location"`
	VerificationRate  float64          `json:"verification_rate"`
}

type JobStatsResponse struct {
	TotalActiveJobs      int64               `json:"total_active_jobs"`
	NewJobsToday         int64               `json:"new_jobs_today"`
	NewJobsThisWeek      int64               `json:"new_jobs_this_week"`
	NewJobsThisMonth     int64               `json:"new_jobs_this_month"`
	JobsByType           []EmploymentTypeCount `json:"jobs_by_type"`
	JobsByLocation       []LocationCount     `json:"jobs_by_location"`
	AverageApplications  float64             `json:"average_applications_per_job"`
	MostViewedJobs       []JobViewCount      `json:"most_viewed_jobs"`
}

type ApplicationStatsResponse struct {
	TotalApplications      int64                 `json:"total_applications"`
	ApplicationsToday      int64                 `json:"applications_today"`
	ApplicationsThisWeek   int64                 `json:"applications_this_week"`
	ApplicationsThisMonth  int64                 `json:"applications_this_month"`
	ApplicationsByStatus   []StatusCount         `json:"applications_by_status"`
	AverageResponseTime    float64               `json:"average_response_time_hours"`
	AcceptanceRate         float64               `json:"acceptance_rate"`
	HireRate               float64               `json:"hire_rate"`
}

type RevenueStatsResponse struct {
	TotalRevenue          float64           `json:"total_revenue"`
	RevenueThisMonth      float64           `json:"revenue_this_month"`
	RevenueLastMonth      float64           `json:"revenue_last_month"`
	MonthlyGrowth         float64           `json:"monthly_growth"`
	SubscriptionsByTier   []TierCount       `json:"subscriptions_by_tier"`
	AverageRevenuePerUser float64           `json:"average_revenue_per_user"`
}

type TopPerformersResponse struct {
	TopRecruiters []RecruiterPerformance `json:"top_recruiters"`
	TopJobs       []JobPerformance       `json:"top_jobs"`
	TopTalents    []TalentPerformance    `json:"top_talents"`
}

// Recruiter Analytics
type RecruiterAnalytics struct {
	Overview       RecruiterOverview        `json:"overview"`
	JobPerformance []JobPerformanceDetail   `json:"job_performance"`
	ApplicationStats RecruiterApplicationStats `json:"application_stats"`
	TrendData      []TrendDataPoint         `json:"trend_data"`
	TopApplicants  []ApplicantInfo          `json:"top_applicants"`
}

type RecruiterOverview struct {
	TotalJobs           int64   `json:"total_jobs"`
	ActiveJobs          int64   `json:"active_jobs"`
	TotalApplications   int64   `json:"total_applications"`
	PendingApplications int64   `json:"pending_applications"`
	TotalJobViews       int64   `json:"total_job_views"`
	TotalHires          int64   `json:"total_hires"`
	ResponseRate        float64 `json:"response_rate"`
	AverageTimeToHire   float64 `json:"average_time_to_hire_days"`
}

type JobPerformanceDetail struct {
	JobID            string    `json:"job_id"`
	Title            string    `json:"title"`
	Views            int64     `json:"views"`
	Applications     int64     `json:"applications"`
	ConversionRate   float64   `json:"conversion_rate"`
	Status           string    `json:"status"`
	PostedDate       time.Time `json:"posted_date"`
}

type RecruiterApplicationStats struct {
	TotalReceived  int64         `json:"total_received"`
	Pending        int64         `json:"pending"`
	Reviewed       int64         `json:"reviewed"`
	Accepted       int64         `json:"accepted"`
	Rejected       int64         `json:"rejected"`
	Hired          int64         `json:"hired"`
	ByStatus       []StatusCount `json:"by_status"`
}

type ApplicantInfo struct {
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	AppliedDate  time.Time `json:"applied_date"`
	JobTitle     string    `json:"job_title"`
	Status       string    `json:"status"`
}

// Talent Analytics
type TalentAnalytics struct {
	Overview          TalentOverview          `json:"overview"`
	ProfileViews      ProfileViewsStats       `json:"profile_views"`
	PortfolioStats    PortfolioStatsResponse  `json:"portfolio_stats"`
	ApplicationStats  TalentApplicationStats  `json:"application_stats"`
	TrendData         []TrendDataPoint        `json:"trend_data"`
	ViewerInsights    ViewerInsightsResponse  `json:"viewer_insights"`
}

type TalentOverview struct {
	TotalProfileViews   int64   `json:"total_profile_views"`
	TotalPortfolioViews int64   `json:"total_portfolio_views"`
	TotalApplications   int64   `json:"total_applications"`
	PendingApplications int64   `json:"pending_applications"`
	AcceptedApplications int64  `json:"accepted_applications"`
	RecruiterViews      int64   `json:"recruiter_views"`
	ProfileCompleteness float64 `json:"profile_completeness"`
}

type ProfileViewsStats struct {
	TotalViews       int64            `json:"total_views"`
	ViewsToday       int64            `json:"views_today"`
	ViewsThisWeek    int64            `json:"views_this_week"`
	ViewsThisMonth   int64            `json:"views_this_month"`
	UniqueViewers    int64            `json:"unique_viewers"`
	RecruiterViewers int64            `json:"recruiter_viewers"`
	ViewsByCountry   []LocationCount  `json:"views_by_country"`
	ViewsByDevice    []DeviceCount    `json:"views_by_device"`
}

type PortfolioStatsResponse struct {
	TotalViews       int64           `json:"total_views"`
	UniqueVisitors   int64           `json:"unique_visitors"`
	AverageDuration  float64         `json:"average_duration_seconds"`
	BounceRate       float64         `json:"bounce_rate"`
	TopReferrers     []ReferrerCount `json:"top_referrers"`
	ViewsByPage      []PageViewCount `json:"views_by_page"`
}

type TalentApplicationStats struct {
	TotalSent    int64         `json:"total_sent"`
	Pending      int64         `json:"pending"`
	Reviewed     int64         `json:"reviewed"`
	Accepted     int64         `json:"accepted"`
	Rejected     int64         `json:"rejected"`
	Hired        int64         `json:"hired"`
	ResponseRate float64       `json:"response_rate"`
	ByStatus     []StatusCount `json:"by_status"`
}

type ViewerInsightsResponse struct {
	RecruiterViews    int64           `json:"recruiter_views"`
	TalentViews       int64           `json:"talent_views"`
	AnonymousViews    int64           `json:"anonymous_views"`
	TopViewerCompanies []CompanyViewCount `json:"top_viewer_companies"`
}

// Common types
type TrendDataPoint struct {
	Date       string `json:"date"`
	Value      int64  `json:"value"`
	Label      string `json:"label,omitempty"`
}

type ProviderCount struct {
	Provider string `json:"provider"`
	Count    int64  `json:"count"`
}

type LocationCount struct {
	Location string `json:"location"`
	Count    int64  `json:"count"`
}

type EmploymentTypeCount struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type TierCount struct {
	Tier  string `json:"tier"`
	Count int64  `json:"count"`
}

type JobViewCount struct {
	JobID string `json:"job_id"`
	Title string `json:"title"`
	Views int64  `json:"views"`
}

type DeviceCount struct {
	Device string `json:"device"`
	Count  int64  `json:"count"`
}

type ReferrerCount struct {
	Referrer string `json:"referrer"`
	Count    int64  `json:"count"`
}

type PageViewCount struct {
	Path  string `json:"path"`
	Views int64  `json:"views"`
}

type CompanyViewCount struct {
	Company string `json:"company"`
	Views   int64  `json:"views"`
}

type RecruiterPerformance struct {
	UserID         string  `json:"user_id"`
	Name           string  `json:"name"`
	TotalJobs      int64   `json:"total_jobs"`
	TotalHires     int64   `json:"total_hires"`
	ResponseRate   float64 `json:"response_rate"`
}

type JobPerformance struct {
	JobID        string  `json:"job_id"`
	Title        string  `json:"title"`
	Company      string  `json:"company"`
	Views        int64   `json:"views"`
	Applications int64   `json:"applications"`
}

type TalentPerformance struct {
	UserID       string `json:"user_id"`
	Name         string `json:"name"`
	ProfileViews int64  `json:"profile_views"`
	Applications int64  `json:"applications"`
}
