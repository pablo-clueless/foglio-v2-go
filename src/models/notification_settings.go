package models

import (
	"time"

	"github.com/google/uuid"
)

type EmailNotificationSettings struct {
	AppUpdates         bool `json:"app_updates"`
	NewMessages        bool `json:"new_messages"`
	JobRecommendations bool `json:"job_recommendations"`
	Newsletter         bool `json:"newsletter"`
	MarketingEmails    bool `json:"marketing_emails"`
}

type PushNotificationSettings struct {
	AppUpdates  bool `json:"app_updates"`
	NewMessages bool `json:"new_messages"`
	Reminders   bool `json:"reminders"`
}

type InAppNotificationSettings struct {
	ActivityUpdates bool `json:"activity_updates"`
	Mentions        bool `json:"mentions"`
	Announcements   bool `json:"announcements"`
}

type NotificationSettings struct {
	ID        uuid.UUID                  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    uuid.UUID                  `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	User      User                       `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Email     *EmailNotificationSettings `gorm:"type:jsonb;serializer:json" json:"email"`
	Push      *PushNotificationSettings  `gorm:"type:jsonb;serializer:json" json:"push"`
	InApp     *InAppNotificationSettings `gorm:"type:jsonb;serializer:json" json:"in_app"`
	CreatedAt time.Time                  `json:"created_at"`
	UpdatedAt time.Time                  `json:"updated_at"`
}

func DefaultNotificationSettings() *NotificationSettings {
	return &NotificationSettings{
		Email: &EmailNotificationSettings{
			AppUpdates:         true,
			NewMessages:        true,
			JobRecommendations: true,
			Newsletter:         false,
			MarketingEmails:    false,
		},
		Push: &PushNotificationSettings{
			AppUpdates:  true,
			NewMessages: true,
			Reminders:   true,
		},
		InApp: &InAppNotificationSettings{
			ActivityUpdates: true,
			Mentions:        true,
			Announcements:   true,
		},
	}
}
