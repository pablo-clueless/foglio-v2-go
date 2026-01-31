package dto

import (
	"foglio/v2/src/models"
	"time"
)

// CreateAnnouncementDto for creating a new announcement
type CreateAnnouncementDto struct {
	Title          string                       `json:"title" binding:"required,min=1,max=200"`
	Content        string                       `json:"content" binding:"required,min=1"`
	Type           models.AnnouncementType      `json:"type" binding:"required,oneof=ANNOUNCEMENT APP_UPDATE SYSTEM_ALERT MAINTENANCE"`
	TargetAudience models.TargetAudience        `json:"target_audience" binding:"required,oneof=ALL_USERS ADMINS_ONLY RECRUITERS_ONLY PREMIUM_ONLY"`
	Priority       models.AnnouncementPriority  `json:"priority" binding:"omitempty,oneof=LOW NORMAL HIGH CRITICAL"`
	ShowAsBanner   bool                         `json:"show_as_banner"`
	BannerColor    *string                      `json:"banner_color,omitempty"`
	ActionURL      *string                      `json:"action_url,omitempty"`
	ActionText     *string                      `json:"action_text,omitempty"`
	ScheduledAt    *time.Time                   `json:"scheduled_at,omitempty"`
	ExpiresAt      *time.Time                   `json:"expires_at,omitempty"`
	PublishNow     bool                         `json:"publish_now"`
	Metadata       map[string]interface{}       `json:"metadata,omitempty"`
}

// UpdateAnnouncementDto for updating an existing announcement
type UpdateAnnouncementDto struct {
	Title          *string                       `json:"title,omitempty" binding:"omitempty,min=1,max=200"`
	Content        *string                       `json:"content,omitempty"`
	Type           *models.AnnouncementType      `json:"type,omitempty" binding:"omitempty,oneof=ANNOUNCEMENT APP_UPDATE SYSTEM_ALERT MAINTENANCE"`
	TargetAudience *models.TargetAudience        `json:"target_audience,omitempty" binding:"omitempty,oneof=ALL_USERS ADMINS_ONLY RECRUITERS_ONLY PREMIUM_ONLY"`
	Priority       *models.AnnouncementPriority  `json:"priority,omitempty" binding:"omitempty,oneof=LOW NORMAL HIGH CRITICAL"`
	ShowAsBanner   *bool                         `json:"show_as_banner,omitempty"`
	BannerColor    *string                       `json:"banner_color,omitempty"`
	ActionURL      *string                       `json:"action_url,omitempty"`
	ActionText     *string                       `json:"action_text,omitempty"`
	ScheduledAt    *time.Time                    `json:"scheduled_at,omitempty"`
	ExpiresAt      *time.Time                    `json:"expires_at,omitempty"`
	Metadata       map[string]interface{}        `json:"metadata,omitempty"`
}

// AnnouncementQueryParams for filtering announcements
type AnnouncementQueryParams struct {
	Page           int     `json:"page" form:"page"`
	Limit          int     `json:"limit" form:"limit"`
	Type           *string `json:"type" form:"type"`
	TargetAudience *string `json:"target_audience" form:"target_audience"`
	IsPublished    *bool   `json:"is_published" form:"is_published"`
	IncludeExpired bool    `json:"include_expired" form:"include_expired"`
}

// AnnouncementResponse includes user-specific status
type AnnouncementResponse struct {
	models.Announcement
	UserStatus *UserAnnouncementStatusResponse `json:"user_status,omitempty"`
}

// UserAnnouncementStatusResponse for user-facing status
type UserAnnouncementStatusResponse struct {
	IsRead      bool       `json:"is_read"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
	IsDismissed bool       `json:"is_dismissed"`
	DismissedAt *time.Time `json:"dismissed_at,omitempty"`
}

// BannerAnnouncementResponse for active banner announcements
type BannerAnnouncementResponse struct {
	ID          string                       `json:"id"`
	Title       string                       `json:"title"`
	Content     string                       `json:"content"`
	Type        models.AnnouncementType      `json:"type"`
	Priority    models.AnnouncementPriority  `json:"priority"`
	BannerColor *string                      `json:"banner_color,omitempty"`
	ActionURL   *string                      `json:"action_url,omitempty"`
	ActionText  *string                      `json:"action_text,omitempty"`
	IsDismissed bool                         `json:"is_dismissed"`
}

// AnnouncementListResponse is a paginated list of announcements
type AnnouncementListResponse struct {
	Data       []AnnouncementResponse `json:"data"`
	TotalItems int                    `json:"total_items"`
	TotalPages int                    `json:"total_pages"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
}

// AdminAnnouncementListResponse is a paginated list for admin
type AdminAnnouncementListResponse struct {
	Data       []models.Announcement `json:"data"`
	TotalItems int                   `json:"total_items"`
	TotalPages int                   `json:"total_pages"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
}
