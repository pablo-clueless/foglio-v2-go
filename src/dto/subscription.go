package dto

import (
	"foglio/v2/src/models"
)

type CreateSubscriptionDto struct {
	Name             string                  `gorm:"not null" json:"name"`
	Description      *string                 `json:"description,omitempty"`
	Type             models.SubscriptionType `gorm:"not null" json:"type"`
	Tier             models.SubscriptionTier `gorm:"not null" json:"tier"`
	Price            float64                 `gorm:"not null" json:"price"`
	Currency         string                  `gorm:"not null;default:'NGN'" json:"currency"`
	BillingCycleDays int                     `gorm:"not null" json:"billing_cycle_days"` // 30 for monthly, 365 for yearly
	TrialPeriodDays  int                     `gorm:"not null;default:0" json:"trial_period_days"`
	Features         map[string]interface{}  `gorm:"type:jsonb;serializer:json" json:"features,omitempty"`
	SortOrder        int                     `gorm:"not null;default:0" json:"sort_order"`
}

type UpdateSubscriptionDto struct {
	Name             *string                  `json:"name,omitempty"`
	Description      *string                  `json:"description,omitempty"`
	Type             *models.SubscriptionType `json:"type,omitempty"`
	Tier             *models.SubscriptionTier `json:"tier,omitempty"`
	Price            *float64                 `json:"price,omitempty"`
	Currency         *string                  `json:"currency,omitempty"`
	BillingCycleDays *int                     `json:"billing_cycle_days,omitempty"` // 30 for monthly, 365 for yearly
	TrialPeriodDays  *int                     `json:"trial_period_days,omitempty"`
	Features         *map[string]interface{}  `gorm:"type:jsonb;serializer:json" json:"features,omitempty"`
	SortOrder        *int                     `gorm:"not null;default:0" json:"sort_order"`
}
