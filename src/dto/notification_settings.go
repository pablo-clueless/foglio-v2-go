package dto

import "foglio/v2/src/models"

type UpdateEmailSettingsDto struct {
	AppUpdates         *bool `json:"app_updates,omitempty"`
	NewMessages        *bool `json:"new_messages,omitempty"`
	JobRecommendations *bool `json:"job_recommendations,omitempty"`
	Newsletter         *bool `json:"newsletter,omitempty"`
	MarketingEmails    *bool `json:"marketing_emails,omitempty"`
}

type UpdatePushSettingsDto struct {
	AppUpdates  *bool `json:"app_updates,omitempty"`
	NewMessages *bool `json:"new_messages,omitempty"`
	Reminders   *bool `json:"reminders,omitempty"`
}

type UpdateInAppSettingsDto struct {
	ActivityUpdates *bool `json:"activity_updates,omitempty"`
	Mentions        *bool `json:"mentions,omitempty"`
	Announcements   *bool `json:"announcements,omitempty"`
}

type UpdateNotificationSettingsDto struct {
	Email *UpdateEmailSettingsDto `json:"email,omitempty"`
	Push  *UpdatePushSettingsDto  `json:"push,omitempty"`
	InApp *UpdateInAppSettingsDto `json:"in_app,omitempty"`
}

type NotificationSettingsResponse struct {
	ID    string                            `json:"id"`
	Email *models.EmailNotificationSettings `json:"email"`
	Push  *models.PushNotificationSettings  `json:"push"`
	InApp *models.InAppNotificationSettings `json:"in_app"`
}
