package services

import (
	"errors"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationService struct {
	database *gorm.DB
	hub      *lib.Hub
}

func NewNotificationService(database *gorm.DB, hub *lib.Hub) *NotificationService {
	return &NotificationService{
		database: database,
		hub:      hub,
	}
}

func (s *NotificationService) SendRealTimeNotification(userID, title, message string, notificationType models.NotificationType, data map[string]interface{}) error {
	notification := models.Notification{
		OwnerID: uuid.Must(uuid.Parse(userID)),
		Type:    notificationType,
		Title:   title,
		Content: message,
		IsRead:  false,
	}

	if err := s.database.Create(&notification).Error; err != nil {
		return err
	}

	s.hub.SendToUser(userID, models.Notification{
		ID:        notification.ID,
		Type:      notification.Type,
		Title:     notification.Title,
		Content:   notification.Content,
		OwnerID:   notification.OwnerID,
		IsRead:    notification.IsRead,
		CreatedAt: notification.CreatedAt,
	})

	return nil
}

func (s *NotificationService) NotifyJobApplication(jobPosterID, applicantID, jobID, jobTitle, applicantName string) error {
	data := map[string]interface{}{
		"job_id":         jobID,
		"job_title":      jobTitle,
		"applicant_id":   applicantID,
		"applicant_name": applicantName,
	}

	return s.SendRealTimeNotification(
		jobPosterID,
		"APPLICATION_SUBMITTED",
		"New Job Application",
		"APPLICATION_SUBMITTED",
		data,
	)
}

func (s *NotificationService) NotifyApplicationAccepted(applicantID, jobPosterID, jobID, jobTitle string) error {
	data := map[string]interface{}{
		"job_id":      jobID,
		"job_title":   jobTitle,
		"employer_id": jobPosterID,
	}

	return s.SendRealTimeNotification(
		applicantID,
		"Your application for "+jobTitle+" has been accepted",
		"Application Accepted!",
		"APPLICATION_ACCEPTED",
		data,
	)
}

func (s *NotificationService) NotifyApplicationRejected(applicantID, jobPosterID, jobID, jobTitle string) error {
	data := map[string]interface{}{
		"job_id":      jobID,
		"job_title":   jobTitle,
		"employer_id": jobPosterID,
	}

	return s.SendRealTimeNotification(
		applicantID,
		"Your application for "+jobTitle+" was not selected",
		"Application Update",
		"APPLICATION_REJECTED",
		data,
	)
}

func (s *NotificationService) GetNotifications(id string, params dto.Pagination) (*dto.PaginatedResponse[models.Notification], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	var notifications []models.Notification
	var totalItems int64

	query := s.database.Model(models.Notification{})

	if err := query.Count(&totalItems).Error; err != nil {
		return &dto.PaginatedResponse[models.Notification]{
			Data:       []models.Notification{},
			Limit:      params.Limit,
			Page:       params.Page,
			TotalItems: 0,
			TotalPages: 0,
		}, err
	}

	offset := (params.Page - 1) * params.Limit

	if err := query.Where("owner_id = ?", id).Offset(offset).Order("created_at DESC").Limit(params.Limit).
		Find(&notifications).Error; err != nil {
		return nil, err
	}

	totalPages := int64(0)
	if params.Limit > 0 {
		totalPages = (totalItems + int64(params.Limit) - 1) / int64(params.Limit)
	}

	return &dto.PaginatedResponse[models.Notification]{
		Data:       notifications,
		TotalItems: int(totalItems),
		TotalPages: int(totalPages),
		Page:       params.Page,
		Limit:      params.Limit,
	}, nil
}

func (s *NotificationService) GetNotification(id string) (*models.Notification, error) {
	var notification *models.Notification

	if err := s.database.Where("id = ?", id).First(&notification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("notification not found")
		}
		return nil, err
	}

	return notification, nil
}

func (s *NotificationService) DeleteNotification(id string) error {
	var notification *models.Notification

	if err := s.database.Where("id = ?", id).First(&notification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("notification not found")
		}
		return err
	}

	if err := s.database.Delete(&notification).Error; err != nil {
		return err
	}

	return nil
}

func (s *NotificationService) ReadNotification(id string) error {
	var notification *models.Notification

	if err := s.database.Where("id = ?", id).First(&notification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("notification not found")
		}
		return err
	}

	notification.IsRead = true

	if err := s.database.Save(&notification).Error; err != nil {
		return err
	}

	return nil
}
