package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnnouncementType string

const (
	AnnouncementTypeAnnouncement AnnouncementType = "ANNOUNCEMENT"
	AnnouncementTypeAppUpdate    AnnouncementType = "APP_UPDATE"
	AnnouncementTypeSystemAlert  AnnouncementType = "SYSTEM_ALERT"
	AnnouncementTypeMaintenance  AnnouncementType = "MAINTENANCE"
)

type TargetAudience string

const (
	TargetAllUsers       TargetAudience = "ALL_USERS"
	TargetAdminsOnly     TargetAudience = "ADMINS_ONLY"
	TargetRecruitersOnly TargetAudience = "RECRUITERS_ONLY"
	TargetPremiumOnly    TargetAudience = "PREMIUM_ONLY"
)

type AnnouncementPriority string

const (
	PriorityLow      AnnouncementPriority = "LOW"
	PriorityNormal   AnnouncementPriority = "NORMAL"
	PriorityHigh     AnnouncementPriority = "HIGH"
	PriorityCritical AnnouncementPriority = "CRITICAL"
)

// Announcement represents admin announcements
type Announcement struct {
	ID             uuid.UUID                `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Title          string                   `gorm:"not null" json:"title"`
	Content        string                   `gorm:"not null" json:"content"`
	Type           AnnouncementType         `gorm:"not null" json:"type"`
	TargetAudience TargetAudience           `gorm:"not null" json:"target_audience"`
	Priority       AnnouncementPriority     `gorm:"not null;default:'NORMAL'" json:"priority"`
	ShowAsBanner   bool                     `gorm:"default:false" json:"show_as_banner"`
	BannerColor    *string                  `json:"banner_color,omitempty"`
	ActionURL      *string                  `json:"action_url,omitempty"`
	ActionText     *string                  `json:"action_text,omitempty"`
	ScheduledAt    *time.Time               `json:"scheduled_at,omitempty"`
	ExpiresAt      *time.Time               `json:"expires_at,omitempty"`
	IsPublished    bool                     `gorm:"default:false" json:"is_published"`
	PublishedAt    *time.Time               `json:"published_at,omitempty"`
	Metadata       map[string]interface{}   `gorm:"type:jsonb;serializer:json" json:"metadata,omitempty"`
	CreatedBy      uuid.UUID                `gorm:"type:uuid;not null;index" json:"created_by"`
	CreatedByUser  User                     `gorm:"foreignKey:CreatedBy;references:ID" json:"created_by_user,omitempty"`
	UserStatuses   []UserAnnouncementStatus `gorm:"foreignKey:AnnouncementID;constraint:OnDelete:CASCADE" json:"user_statuses,omitempty"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
	DeletedAt      gorm.DeletedAt           `gorm:"index" json:"-"`
}

// UserAnnouncementStatus tracks user interaction with announcements
type UserAnnouncementStatus struct {
	ID             uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID         uuid.UUID    `gorm:"type:uuid;not null;index" json:"user_id"`
	User           User         `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	AnnouncementID uuid.UUID    `gorm:"type:uuid;not null;index" json:"announcement_id"`
	Announcement   Announcement `gorm:"foreignKey:AnnouncementID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	IsRead         bool         `gorm:"default:false" json:"is_read"`
	ReadAt         *time.Time   `json:"read_at,omitempty"`
	IsDismissed    bool         `gorm:"default:false" json:"is_dismissed"`
	DismissedAt    *time.Time   `json:"dismissed_at,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

// IsActive returns true if the announcement should be visible
func (a *Announcement) IsActive() bool {
	now := time.Now()

	if !a.IsPublished {
		return false
	}

	if a.ScheduledAt != nil && a.ScheduledAt.After(now) {
		return false
	}

	if a.ExpiresAt != nil && a.ExpiresAt.Before(now) {
		return false
	}

	return true
}

// IsVisibleToUser checks if the announcement should be shown to a specific user
func (a *Announcement) IsVisibleToUser(user *User) bool {
	if !a.IsActive() {
		return false
	}

	switch a.TargetAudience {
	case TargetAllUsers:
		return true
	case TargetAdminsOnly:
		return user.IsAdmin
	case TargetRecruitersOnly:
		return user.IsRecruiter
	case TargetPremiumOnly:
		return user.IsPremium
	default:
		return false
	}
}
