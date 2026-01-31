package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventType string

const (
	EventPageView        EventType = "page_view"
	EventJobView         EventType = "job_view"
	EventProfileView     EventType = "profile_view"
	EventPortfolioView   EventType = "portfolio_view"
	EventApplicationSent EventType = "application_sent"
	EventJobCreated      EventType = "job_created"
	EventUserSignup      EventType = "user_signup"
	EventUserLogin       EventType = "user_login"
)

// PageView tracks individual page views
type PageView struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Path         string         `gorm:"not null;index" json:"path"`
	UserID       *uuid.UUID     `gorm:"type:uuid;index" json:"user_id,omitempty"`
	SessionID    string         `gorm:"index" json:"session_id"`
	IPAddress    string         `json:"ip_address"`
	UserAgent    string         `json:"user_agent"`
	Referrer     *string        `json:"referrer,omitempty"`
	Country      *string        `gorm:"index" json:"country,omitempty"`
	City         *string        `json:"city,omitempty"`
	DeviceType   *string        `gorm:"index" json:"device_type,omitempty"` // desktop, mobile, tablet
	Browser      *string        `json:"browser,omitempty"`
	OS           *string        `json:"os,omitempty"`
	Duration     *int           `json:"duration,omitempty"` // time spent in seconds
	CreatedAt    time.Time      `gorm:"index" json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// JobView tracks job listing views
type JobView struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	JobID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"job_id"`
	Job       *Job           `gorm:"foreignKey:JobID" json:"job,omitempty"`
	UserID    *uuid.UUID     `gorm:"type:uuid;index" json:"user_id,omitempty"`
	SessionID string         `gorm:"index" json:"session_id"`
	IPAddress string         `json:"ip_address"`
	UserAgent string         `json:"user_agent"`
	Referrer  *string        `json:"referrer,omitempty"`
	Country   *string        `gorm:"index" json:"country,omitempty"`
	CreatedAt time.Time      `gorm:"index" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ProfileView tracks user profile views
type ProfileView struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ProfileUserID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"profile_user_id"` // whose profile was viewed
	ViewerUserID    *uuid.UUID     `gorm:"type:uuid;index" json:"viewer_user_id,omitempty"` // who viewed (if logged in)
	SessionID       string         `gorm:"index" json:"session_id"`
	IPAddress       string         `json:"ip_address"`
	UserAgent       string         `json:"user_agent"`
	Referrer        *string        `json:"referrer,omitempty"`
	Country         *string        `gorm:"index" json:"country,omitempty"`
	ViewerIsRecruiter bool         `json:"viewer_is_recruiter"`
	CreatedAt       time.Time      `gorm:"index" json:"created_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// PortfolioView tracks portfolio page views
type PortfolioView struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PortfolioID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"portfolio_id"`
	Portfolio    *Portfolio     `gorm:"foreignKey:PortfolioID" json:"portfolio,omitempty"`
	ViewerUserID *uuid.UUID     `gorm:"type:uuid;index" json:"viewer_user_id,omitempty"`
	SessionID    string         `gorm:"index" json:"session_id"`
	IPAddress    string         `json:"ip_address"`
	UserAgent    string         `json:"user_agent"`
	Referrer     *string        `json:"referrer,omitempty"`
	Country      *string        `gorm:"index" json:"country,omitempty"`
	DeviceType   *string        `gorm:"index" json:"device_type,omitempty"`
	Duration     *int           `json:"duration,omitempty"`
	CreatedAt    time.Time      `gorm:"index" json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// AnalyticsEvent for tracking custom events
type AnalyticsEvent struct {
	ID         uuid.UUID              `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EventType  EventType              `gorm:"not null;index" json:"event_type"`
	UserID     *uuid.UUID             `gorm:"type:uuid;index" json:"user_id,omitempty"`
	EntityID   *uuid.UUID             `gorm:"type:uuid;index" json:"entity_id,omitempty"` // job_id, profile_id, etc.
	EntityType *string                `gorm:"index" json:"entity_type,omitempty"`        // job, profile, portfolio
	SessionID  string                 `gorm:"index" json:"session_id"`
	Properties map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"properties,omitempty"`
	CreatedAt  time.Time              `gorm:"index" json:"created_at"`
	DeletedAt  gorm.DeletedAt         `gorm:"index" json:"-"`
}

// DailyStats for pre-aggregated daily statistics
type DailyStats struct {
	ID                uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Date              time.Time      `gorm:"type:date;not null;uniqueIndex:idx_daily_stats_date_type" json:"date"`
	StatType          string         `gorm:"not null;uniqueIndex:idx_daily_stats_date_type" json:"stat_type"` // platform, recruiter, talent
	EntityID          *uuid.UUID     `gorm:"type:uuid;index" json:"entity_id,omitempty"`                      // user_id for recruiter/talent stats
	TotalPageViews    int            `gorm:"default:0" json:"total_page_views"`
	UniqueVisitors    int            `gorm:"default:0" json:"unique_visitors"`
	TotalJobViews     int            `gorm:"default:0" json:"total_job_views"`
	TotalProfileViews int            `gorm:"default:0" json:"total_profile_views"`
	TotalApplications int            `gorm:"default:0" json:"total_applications"`
	NewUsers          int            `gorm:"default:0" json:"new_users"`
	NewJobs           int            `gorm:"default:0" json:"new_jobs"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (pv *PageView) BeforeCreate(tx *gorm.DB) error {
	pv.CreatedAt = time.Now()
	return nil
}

func (jv *JobView) BeforeCreate(tx *gorm.DB) error {
	jv.CreatedAt = time.Now()
	return nil
}

func (pv *ProfileView) BeforeCreate(tx *gorm.DB) error {
	pv.CreatedAt = time.Now()
	return nil
}

func (pv *PortfolioView) BeforeCreate(tx *gorm.DB) error {
	pv.CreatedAt = time.Now()
	return nil
}

func (ae *AnalyticsEvent) BeforeCreate(tx *gorm.DB) error {
	ae.CreatedAt = time.Now()
	return nil
}

func (ds *DailyStats) BeforeCreate(tx *gorm.DB) error {
	ds.CreatedAt = time.Now()
	ds.UpdatedAt = time.Now()
	return nil
}

func (ds *DailyStats) BeforeUpdate(tx *gorm.DB) error {
	ds.UpdatedAt = time.Now()
	return nil
}
