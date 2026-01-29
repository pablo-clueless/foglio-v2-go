package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionType string

const (
	SubscriptionMonthly  SubscriptionType = "monthly"
	SubscriptionYearly   SubscriptionType = "yearly"
	SubscriptionLifetime SubscriptionType = "lifetime"
)

type SubscriptionTier string

const (
	TierFree     SubscriptionTier = "free"
	TierBasic    SubscriptionTier = "basic"
	TierPremium  SubscriptionTier = "premium"
	TierBusiness SubscriptionTier = "business"
)

type Subscription struct {
	ID               uuid.UUID              `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name             string                 `gorm:"not null" json:"name"`
	Description      *string                `json:"description,omitempty"`
	Type             SubscriptionType       `gorm:"not null" json:"type"`
	Tier             SubscriptionTier       `gorm:"not null" json:"tier"`
	Price            float64                `gorm:"not null" json:"price"`
	Currency         string                 `gorm:"not null;default:'USD'" json:"currency"`
	BillingCycleDays int                    `gorm:"not null" json:"billing_cycle_days"` // 30 for monthly, 365 for yearly
	TrialPeriodDays  int                    `gorm:"not null;default:0" json:"trial_period_days"`
	Features         map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"features,omitempty"`
	MaxProjects      int                    `gorm:"not null;default:5" json:"max_projects"`
	MaxSkills        int                    `gorm:"not null;default:20" json:"max_skills"`
	MaxExperiences   int                    `gorm:"not null;default:10" json:"max_experiences"`
	IsActive         bool                   `gorm:"not null;default:true" json:"is_active"`
	IsPopular        bool                   `gorm:"not null;default:false" json:"is_popular"`
	SortOrder        int                    `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	DeletedAt        gorm.DeletedAt         `gorm:"index" json:"-"`
}

type UserSubscription struct {
	ID                     uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID                 uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	User                   *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	SubscriptionID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"subscription_id"`
	Subscription           *Subscription  `gorm:"foreignKey:SubscriptionID" json:"subscription,omitempty"`
	PaystackCustomerID     *string        `gorm:"index" json:"paystack_customer_id,omitempty"`
	PaystackSubscriptionID *string        `gorm:"index" json:"paystack_subscription_id,omitempty"`
	PaymentMethodID        *string        `json:"payment_method_id,omitempty"`
	LastPaymentAmount      *float64       `json:"last_payment_amount,omitempty"`
	LastPaymentDate        *time.Time     `json:"last_payment_date,omitempty"`
	Status                 string         `gorm:"not null;default:'active'" json:"status"` // active, canceled, expired, trialing
	IsActive               bool           `gorm:"not null;default:true" json:"is_active"`
	CurrentPeriodStart     time.Time      `gorm:"not null" json:"current_period_start"`
	CurrentPeriodEnd       time.Time      `gorm:"not null" json:"current_period_end"`
	CancelAtPeriodEnd      bool           `gorm:"not null;default:false" json:"cancel_at_period_end"`
	TrialStart             *time.Time     `json:"trial_start,omitempty"`
	TrialEnd               *time.Time     `json:"trial_end,omitempty"`
	CancelledAt            *time.Time     `json:"cancelled_at,omitempty"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`
}

type SubscriptionInvoice struct {
	ID                 uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserSubscriptionID uuid.UUID         `gorm:"type:uuid;not null;index" json:"user_subscription_id"`
	UserSubscription   *UserSubscription `gorm:"foreignKey:UserSubscriptionID" json:"user_subscription,omitempty"`
	PaystackReference  string            `gorm:"not null;uniqueIndex" json:"paystack_reference"` // Unique for idempotency
	AmountPaid         float64           `gorm:"not null" json:"amount_paid"`
	Currency           string            `gorm:"not null" json:"currency"`
	Status             string            `gorm:"not null" json:"status"` // paid, failed, void
	InvoicePDF         *string           `json:"invoice_pdf,omitempty"`
	PeriodStart        time.Time         `json:"period_start"`
	PeriodEnd          time.Time         `json:"period_end"`
	PaidAt             *time.Time        `json:"paid_at,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
}

type PaystackPlan struct {
	ID             uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	SubscriptionID uuid.UUID     `gorm:"type:uuid;not null;index" json:"subscription_id"`
	Subscription   *Subscription `gorm:"foreignKey:SubscriptionID" json:"subscription,omitempty"`
	PlanCode       string        `gorm:"not null;uniqueIndex" json:"plan_code"`
	PaystackPlanID int           `gorm:"not null" json:"paystack_plan_id"`
	Interval       string        `gorm:"not null" json:"interval"` // monthly, yearly
	IsActive       bool          `gorm:"not null;default:true" json:"is_active"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}
