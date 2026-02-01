package services

import (
	"errors"
	"fmt"
	"foglio/v2/src/dto"
	"foglio/v2/src/models"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionService struct {
	database *gorm.DB
}

func NewSubscriptionService(database *gorm.DB) *SubscriptionService {
	return &SubscriptionService{
		database: database,
	}
}

func (s *SubscriptionService) CreateSubscriptionTier(payload dto.CreateSubscriptionDto) (*models.Subscription, error) {
	subscription := &models.Subscription{
		Name:             payload.Name,
		Description:      payload.Description,
		Type:             payload.Type,
		Tier:             payload.Tier,
		Price:            payload.Price,
		Currency:         payload.Currency,
		BillingCycleDays: payload.BillingCycleDays,
		TrialPeriodDays:  payload.TrialPeriodDays,
		Features:         payload.Features,
		SortOrder:        payload.SortOrder,
		IsActive:         true,
	}

	if err := s.database.Create(subscription).Error; err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *SubscriptionService) UpdateSubscriptionTier(id string, payload dto.UpdateSubscriptionDto) (*models.Subscription, error) {
	var subscription models.Subscription
	if err := s.database.First(&subscription, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if payload.Name != nil {
		subscription.Name = *payload.Name
	}
	if payload.Description != nil {
		subscription.Description = payload.Description
	}
	if payload.Type != nil {
		subscription.Type = *payload.Type
	}
	if payload.Tier != nil {
		subscription.Tier = *payload.Tier
	}
	if payload.Price != nil {
		subscription.Price = *payload.Price
	}
	if payload.Currency != nil {
		subscription.Currency = *payload.Currency
	}
	if payload.BillingCycleDays != nil {
		subscription.BillingCycleDays = *payload.BillingCycleDays
	}
	if payload.TrialPeriodDays != nil {
		subscription.TrialPeriodDays = *payload.TrialPeriodDays
	}
	if payload.Features != nil {
		subscription.Features = *payload.Features
	}
	if payload.SortOrder != nil {
		subscription.SortOrder = *payload.SortOrder
	}

	if err := s.database.Save(&subscription).Error; err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (s *SubscriptionService) DeleteSubscriptionTier(id string) error {
	var subscription models.Subscription
	if err := s.database.First(&subscription, "id = ?", id).Error; err != nil {
		return err
	}

	return s.database.Delete(&subscription).Error
}

func (s *SubscriptionService) GetSubscriptions(params *dto.Pagination) (*dto.PaginatedResponse[models.Subscription], error) {
	q := normalizeSubscriptionQuery(*params)

	var tiers []models.Subscription
	var totalItems int64

	db := s.database.Model(&models.Subscription{})
	db.Count(&totalItems)

	err := db.Limit(q.Limit).Offset((q.Page - 1) * q.Limit).Find(&tiers).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(totalItems) / q.Limit
	if int(totalItems)%q.Limit != 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[models.Subscription]{
		Data:       tiers,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
		Page:       q.Page,
		Limit:      q.Limit,
	}, nil
}

func (s *SubscriptionService) GetSubscriptionById(id string) (*models.Subscription, error) {
	var subscription models.Subscription
	err := s.database.First(&subscription, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (s *SubscriptionService) CreateUserSubscription(userSub *models.UserSubscription) error {
	return s.database.Create(userSub).Error
}

func (s *SubscriptionService) UpdateUserSubscription(id string, userSub *models.UserSubscription) error {
	return s.database.Model(&models.UserSubscription{}).Where("id = ?", id).Updates(userSub).Error
}

func (s *SubscriptionService) DeleteUserSubscription(id string) error {
	return s.database.Delete(&models.UserSubscription{}, "id = ?", id).Error
}

func (s *SubscriptionService) GetUserSubscriptions(userId string, params *dto.Pagination) (*dto.PaginatedResponse[models.UserSubscription], error) {
	q := normalizeSubscriptionQuery(*params)

	var subs []models.UserSubscription
	var totalItems int64

	db := s.database.Model(&models.UserSubscription{}).Where("user_id = ?", userId)
	db.Count(&totalItems)

	err := db.Limit(q.Limit).Offset((q.Page - 1) * q.Limit).Find(&subs).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(totalItems) / q.Limit
	if int(totalItems)%q.Limit != 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[models.UserSubscription]{
		Data:       subs,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
		Page:       q.Page,
		Limit:      q.Limit,
	}, nil
}

func (s *SubscriptionService) GetUserSubscriptionById(id string) (*models.UserSubscription, error) {
	var sub models.UserSubscription
	err := s.database.First(&sub, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (s *SubscriptionService) SubscribeUser(userId string, tierId string) error {
	tx := s.database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingSub models.UserSubscription
	err := tx.Where("user_id = ?", userId).First(&existingSub).Error
	if err == nil {
		tx.Rollback()
		return errors.New("user already has an active subscription")
	}

	var user models.User
	if err := tx.First(&user, "id = ?", userId).Error; err != nil {
		tx.Rollback()
		return err
	}

	newSub := models.UserSubscription{
		UserID:         uuid.MustParse(userId),
		SubscriptionID: uuid.MustParse(tierId),
		Status:         "active",
	}

	if err := tx.Create(&newSub).Error; err != nil {
		tx.Rollback()
		return err
	}

	if user.Domain == nil || user.Domain.Subdomain == "" {
		subdomain := s.generateUniqueSubdomain(tx, user.Username)
		if user.Domain == nil {
			user.Domain = &models.Domain{}
		}
		user.Domain.Subdomain = subdomain
		if err := tx.Save(&user).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (s *SubscriptionService) generateUniqueSubdomain(tx *gorm.DB, username string) string {
	baseSubdomain := strings.ToLower(username)
	subdomain := baseSubdomain

	counter := 1
	for {
		var count int64
		tx.Model(&models.User{}).
			Where("domain->>'subdomain' = ?", subdomain).
			Count(&count)
		if count == 0 {
			break
		}
		subdomain = fmt.Sprintf("%s%d", baseSubdomain, counter)
		counter++
	}

	return subdomain
}

func (s *SubscriptionService) UpgradeUserSubscription(userId string, newTierId string) error {
	tx := s.database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var sub models.UserSubscription
	err := tx.Where("user_id = ? AND status = ?", userId, "active").First(&sub).Error
	if err != nil {
		tx.Rollback()
		return errors.New("no active subscription found")
	}

	sub.SubscriptionID = uuid.MustParse(newTierId)

	if err := tx.Save(&sub).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *SubscriptionService) DowngradeUserSubscription(userId string, newTierId string) error {
	tx := s.database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var sub models.UserSubscription
	err := tx.Where("user_id = ? AND status = ?", userId, "active").First(&sub).Error
	if err != nil {
		tx.Rollback()
		return errors.New("no active subscription found")
	}

	sub.SubscriptionID = uuid.MustParse(newTierId)

	if err := tx.Save(&sub).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *SubscriptionService) UnsubscribeUser(userId string) error {
	tx := s.database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var sub models.UserSubscription
	err := tx.Where("user_id = ? AND status = ?", userId, "active").First(&sub).Error
	if err != nil {
		tx.Rollback()
		return errors.New("no active subscription found")
	}

	sub.Status = "cancelled"

	if err := tx.Save(&sub).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func normalizeSubscriptionQuery(q dto.Pagination) dto.Pagination {
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Page <= 0 {
		q.Page = 1
	}

	return q
}

func (s *SubscriptionService) ProcessExpiredSubscriptions() error {
	now := time.Now()

	var expiredSubs []models.UserSubscription
	err := s.database.
		Where("status = ?", "active").
		Where("is_active = ?", true).
		Where("current_period_end < ?", now).
		Find(&expiredSubs).Error

	if err != nil {
		return err
	}

	if len(expiredSubs) == 0 {
		return nil
	}

	log.Printf("Found %d expired subscriptions to process", len(expiredSubs))

	for _, sub := range expiredSubs {
		tx := s.database.Begin()
		sub.Status = "expired"
		sub.IsActive = false
		if err := tx.Save(&sub).Error; err != nil {
			tx.Rollback()
			log.Printf("Failed to update subscription %s: %v", sub.ID, err)
			continue
		}
		if err := tx.Model(&models.User{}).
			Where("id = ?", sub.UserID).
			Update("is_premium", false).Error; err != nil {
			tx.Rollback()
			log.Printf("Failed to update user %s premium status: %v", sub.UserID, err)
			continue
		}

		if err := tx.Commit().Error; err != nil {
			log.Printf("Failed to commit transaction for subscription %s: %v", sub.ID, err)
			continue
		}

		log.Printf("Expired subscription %s for user %s", sub.ID, sub.UserID)
	}

	return nil
}
