package services

import (
	"errors"
	"foglio/v2/src/dto"
	"foglio/v2/src/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationSettingsService struct {
	database *gorm.DB
}

func NewNotificationSettingsService(database *gorm.DB) *NotificationSettingsService {
	return &NotificationSettingsService{database: database}
}

// GetOrCreateSettings retrieves or creates default notification settings for a user
func (s *NotificationSettingsService) GetOrCreateSettings(userID string) (*models.NotificationSettings, error) {
	var settings models.NotificationSettings

	err := s.database.Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create default settings
			defaults := models.DefaultNotificationSettings()
			defaults.UserID = uuid.Must(uuid.Parse(userID))

			if err := s.database.Create(defaults).Error; err != nil {
				return nil, err
			}
			return defaults, nil
		}
		return nil, err
	}

	return &settings, nil
}

// UpdateSettings updates notification settings for a user
func (s *NotificationSettingsService) UpdateSettings(userID string, payload dto.UpdateNotificationSettingsDto) (*models.NotificationSettings, error) {
	settings, err := s.GetOrCreateSettings(userID)
	if err != nil {
		return nil, err
	}

	// Update email settings
	if payload.Email != nil {
		if settings.Email == nil {
			settings.Email = &models.EmailNotificationSettings{}
		}
		if payload.Email.AppUpdates != nil {
			settings.Email.AppUpdates = *payload.Email.AppUpdates
		}
		if payload.Email.NewMessages != nil {
			settings.Email.NewMessages = *payload.Email.NewMessages
		}
		if payload.Email.JobRecommendations != nil {
			settings.Email.JobRecommendations = *payload.Email.JobRecommendations
		}
		if payload.Email.Newsletter != nil {
			settings.Email.Newsletter = *payload.Email.Newsletter
		}
		if payload.Email.MarketingEmails != nil {
			settings.Email.MarketingEmails = *payload.Email.MarketingEmails
		}
	}

	// Update push settings
	if payload.Push != nil {
		if settings.Push == nil {
			settings.Push = &models.PushNotificationSettings{}
		}
		if payload.Push.AppUpdates != nil {
			settings.Push.AppUpdates = *payload.Push.AppUpdates
		}
		if payload.Push.NewMessages != nil {
			settings.Push.NewMessages = *payload.Push.NewMessages
		}
		if payload.Push.Reminders != nil {
			settings.Push.Reminders = *payload.Push.Reminders
		}
	}

	// Update in-app settings
	if payload.InApp != nil {
		if settings.InApp == nil {
			settings.InApp = &models.InAppNotificationSettings{}
		}
		if payload.InApp.ActivityUpdates != nil {
			settings.InApp.ActivityUpdates = *payload.InApp.ActivityUpdates
		}
		if payload.InApp.Mentions != nil {
			settings.InApp.Mentions = *payload.InApp.Mentions
		}
		if payload.InApp.Announcements != nil {
			settings.InApp.Announcements = *payload.InApp.Announcements
		}
	}

	if err := s.database.Save(settings).Error; err != nil {
		return nil, err
	}

	return settings, nil
}

// ShouldSendEmailNotification checks if a specific email notification type is enabled
func (s *NotificationSettingsService) ShouldSendEmailNotification(userID string, notificationType string) (bool, error) {
	settings, err := s.GetOrCreateSettings(userID)
	if err != nil {
		return false, err
	}

	if settings.Email == nil {
		return true, nil // Default to enabled if no settings
	}

	switch notificationType {
	case "app_updates":
		return settings.Email.AppUpdates, nil
	case "new_messages":
		return settings.Email.NewMessages, nil
	case "job_recommendations":
		return settings.Email.JobRecommendations, nil
	case "newsletter":
		return settings.Email.Newsletter, nil
	case "marketing_emails":
		return settings.Email.MarketingEmails, nil
	default:
		return true, nil
	}
}

// ShouldSendPushNotification checks if a specific push notification type is enabled
func (s *NotificationSettingsService) ShouldSendPushNotification(userID string, notificationType string) (bool, error) {
	settings, err := s.GetOrCreateSettings(userID)
	if err != nil {
		return false, err
	}

	if settings.Push == nil {
		return true, nil
	}

	switch notificationType {
	case "app_updates":
		return settings.Push.AppUpdates, nil
	case "new_messages":
		return settings.Push.NewMessages, nil
	case "reminders":
		return settings.Push.Reminders, nil
	default:
		return true, nil
	}
}

// ShouldSendInAppNotification checks if a specific in-app notification type is enabled
func (s *NotificationSettingsService) ShouldSendInAppNotification(userID string, notificationType string) (bool, error) {
	settings, err := s.GetOrCreateSettings(userID)
	if err != nil {
		return false, err
	}

	if settings.InApp == nil {
		return true, nil
	}

	switch notificationType {
	case "activity_updates":
		return settings.InApp.ActivityUpdates, nil
	case "mentions":
		return settings.InApp.Mentions, nil
	case "announcements":
		return settings.InApp.Announcements, nil
	default:
		return true, nil
	}
}

// MapNotificationTypeToSetting maps notification types to settings keys
func MapNotificationTypeToSetting(notificationType models.NotificationType) string {
	switch notificationType {
	case models.ApplicationSubmitted, models.ApplicationAccepted, models.ApplicationRejected:
		return "activity_updates"
	case models.NewMessage:
		return "new_messages"
	case models.System:
		return "announcements"
	default:
		return "activity_updates"
	}
}
