package services

import (
	"errors"
	"foglio/v2/src/dto"
	"foglio/v2/src/models"

	"gorm.io/gorm"
)

type NotificationService struct {
	database *gorm.DB
}

func NewNotificationService(database *gorm.DB) *NotificationService {
	return &NotificationService{
		database: database,
	}
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
