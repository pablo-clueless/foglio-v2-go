package database

import (
	"log"

	"foglio/v2/src/models"

	"gorm.io/gorm"
)

func SeedSubscriptions(db *gorm.DB) error {
	subscriptions := []models.Subscription{
		{
			Name:             "Free",
			Description:      strPtr("Get started with basic features"),
			Type:             models.SubscriptionMonthly,
			Tier:             models.TierFree,
			Price:            0,
			Currency:         "NGN",
			BillingCycleDays: 30,
			TrialPeriodDays:  0,
			MaxProjects:      3,
			MaxSkills:        10,
			MaxExperiences:   5,
			IsActive:         true,
			IsPopular:        false,
			SortOrder:        1,
			Features: map[string]interface{}{
				"basic_profile":    true,
				"job_applications": 5,
				"resume_downloads": 1,
			},
		},
		{
			Name:             "Pro Monthly",
			Description:      strPtr("Unlock all features with monthly billing"),
			Type:             models.SubscriptionMonthly,
			Tier:             models.TierPremium,
			Price:            1200,
			Currency:         "NGN",
			BillingCycleDays: 30,
			TrialPeriodDays:  7,
			MaxProjects:      20,
			MaxSkills:        50,
			MaxExperiences:   20,
			IsActive:         true,
			IsPopular:        true,
			SortOrder:        2,
			Features: map[string]interface{}{
				"basic_profile":      true,
				"premium_profile":    true,
				"job_applications":   -1,
				"resume_downloads":   -1,
				"priority_support":   true,
				"analytics":          true,
				"custom_domain":      false,
				"verified_badge":     true,
			},
		},
		{
			Name:             "Pro Yearly",
			Description:      strPtr("Save 20% with annual billing"),
			Type:             models.SubscriptionYearly,
			Tier:             models.TierPremium,
			Price:            11520, // 1200 * 12 * 0.8 = 20% discount
			Currency:         "NGN",
			BillingCycleDays: 365,
			TrialPeriodDays:  14,
			MaxProjects:      20,
			MaxSkills:        50,
			MaxExperiences:   20,
			IsActive:         true,
			IsPopular:        false,
			SortOrder:        3,
			Features: map[string]interface{}{
				"basic_profile":      true,
				"premium_profile":    true,
				"job_applications":   -1,
				"resume_downloads":   -1,
				"priority_support":   true,
				"analytics":          true,
				"custom_domain":      true,
				"verified_badge":     true,
			},
		},
	}

	for _, sub := range subscriptions {
		var existing models.Subscription
		result := db.Where("name = ? AND type = ?", sub.Name, sub.Type).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&sub).Error; err != nil {
				log.Printf("Failed to seed subscription %s: %v", sub.Name, err)
				return err
			}
			log.Printf("Seeded subscription: %s", sub.Name)
		} else if result.Error != nil {
			return result.Error
		} else {
			log.Printf("Subscription %s already exists, skipping", sub.Name)
		}
	}

	return nil
}

func strPtr(s string) *string {
	return &s
}

func RunSeeds() error {
	db := GetDatabase()

	log.Println("Running database seeds...")

	if err := SeedSubscriptions(db); err != nil {
		return err
	}

	log.Println("Database seeds completed successfully")
	return nil
}
