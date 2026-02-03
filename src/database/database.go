package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"foglio/v2/src/config"
	"foglio/v2/src/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SchemaMigration struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"uniqueIndex;not null"`
	AppliedAt time.Time `gorm:"not null"`
}

var (
	client                    *gorm.DB
	once                      sync.Once
	ErrDatabaseNotInitialized = errors.New("database not initialized")
)

const (
	maxOpenConns    = 25
	maxIdleConns    = 10
	connMaxLifetime = 5 * time.Minute
	connMaxIdleTime = 30 * time.Second
	pingTimeout     = 5 * time.Second
)

func InitializeDatabase() error {
	var initErr error
	once.Do(func() {
		url := config.AppConfig.PostgresUrl
		if url == "" {
			initErr = errors.New("database URL not configured")
			return
		}

		logLevel := logger.Error
		if config.AppConfig.IsDevMode {
			logLevel = logger.Info
		}

		database, err := gorm.Open(postgres.Open(url), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
		})
		if err != nil {
			initErr = err
			return
		}

		if err := configureConnectionPool(database); err != nil {
			initErr = err
			return
		}

		if err := EnableUUIDExtension(database); err != nil {
			log.Printf("Failed to enable UUID extension: %v", err)
			initErr = err
			return
		}

		if err := runMigrations(database); err != nil {
			log.Printf("Failed to run migrations: %v", err)
			initErr = err
			return
		}

		if err := PingDatabase(database); err != nil {
			initErr = err
			return
		}

		client = database
		log.Println("Database connection established successfully")
	})

	return initErr
}

func GetDatabase() *gorm.DB {
	if client == nil {
		log.Fatal("Database not initialized. Call InitializeDatabase() first")
	}
	return client
}

func CloseDatabase() error {
	if client == nil {
		return ErrDatabaseNotInitialized
	}

	sqlDB, err := client.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Close(); err != nil {
		return err
	}

	log.Println("Database connection closed")
	return nil
}

func HealthCheck() error {
	if client == nil {
		return ErrDatabaseNotInitialized
	}
	return PingDatabase(client)
}

func configureConnectionPool(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)

	return nil
}

func EnableUUIDExtension(db *gorm.DB) error {
	extensions := []string{
		"CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"",
		"CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"",
	}

	for _, ext := range extensions {
		if err := db.Exec(ext).Error; err != nil {
			log.Printf("Warning: Failed to create extension with query '%s': %v", ext, err)
		}
	}

	if err := createEnumTypes(db); err != nil {
		return err
	}

	return nil
}

func createEnumTypes(db *gorm.DB) error {
	enumTypes := []struct {
		name   string
		values []string
	}{
		{
			name:   "verification_type",
			values: []string{"DRIVERS_LICENSE", "INTERNATIONAL_PASSPORT", "NATIONAL_ID_CARD", "VOTERS_CARD"},
		},
	}

	for _, enum := range enumTypes {
		var exists bool
		db.Raw("SELECT EXISTS (SELECT 1 FROM pg_type WHERE typname = ?)", enum.name).Scan(&exists)
		if exists {
			continue
		}

		values := make([]string, len(enum.values))
		for i, v := range enum.values {
			values[i] = fmt.Sprintf("'%s'", v)
		}
		sql := fmt.Sprintf("CREATE TYPE %s AS ENUM (%s)", enum.name, strings.Join(values, ", "))
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("Warning: Failed to create enum type %s: %v", enum.name, err)
		} else {
			log.Printf("Created enum type: %s", enum.name)
		}
	}

	return nil
}

func runMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(&SchemaMigration{}); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	var appliedMigrations []SchemaMigration
	if err := db.Find(&appliedMigrations).Error; err != nil {
		return fmt.Errorf("failed to fetch applied migrations: %w", err)
	}

	appliedMap := make(map[string]bool)
	for _, m := range appliedMigrations {
		appliedMap[m.Name] = true
	}

	migrations := []struct {
		name  string
		model any
	}{
		{"001_create_companies", &models.Company{}},
		{"002_create_users", &models.User{}},
		{"003_create_languages", &models.Language{}},
		{"004_create_certifications", &models.Certification{}},
		{"005_create_education", &models.Education{}},
		{"006_create_education_highlights", &models.EducationHighlight{}},
		{"007_create_experiences", &models.Experience{}},
		{"008_create_experience_highlights", &models.ExperienceHighlight{}},
		{"009_create_experience_tech", &models.ExperienceTech{}},
		{"010_create_projects", &models.Project{}},
		{"011_create_project_stacks", &models.ProjectStack{}},
		{"012_create_project_highlights", &models.ProjectHighlight{}},
		{"013_create_jobs", &models.Job{}},
		{"014_create_job_applications", &models.JobApplication{}},
		{"015_create_comments", &models.Comment{}},
		{"016_create_reactions", &models.Reaction{}},
		{"017_create_notifications", &models.Notification{}},
		{"018_create_social_media", &models.SocialMedia{}},
		{"019_create_subscription", &models.Subscription{}},
		{"020_create_subscription_invoices", &models.SubscriptionInvoice{}},
		{"021_create_user_subscriptions", &models.UserSubscription{}},
		{"022_create_paystack_plans", &models.PaystackPlan{}},
		{"023_create_portfolios", &models.Portfolio{}},
		{"024_create_portfolio_sections", &models.PortfolioSection{}},
		{"025_create_page_views", &models.PageView{}},
		{"026_create_job_views", &models.JobView{}},
		{"027_create_profile_views", &models.ProfileView{}},
		{"028_create_portfolio_views", &models.PortfolioView{}},
		{"029_create_analytics_events", &models.AnalyticsEvent{}},
		{"030_create_daily_stats", &models.DailyStats{}},
		{"031_add_two_factor_fields", &models.User{}},
		{"032_create_notification_settings", &models.NotificationSettings{}},
		{"033_create_announcements", &models.Announcement{}},
		{"034_create_user_announcement_status", &models.UserAnnouncementStatus{}},
		{"035_create_conversations", &models.Conversation{}},
		{"036_create_messages", &models.Message{}},
		{"038_add_message_media", &models.Message{}},
		{"039_create_reviews", &models.Review{}},
	}

	pendingCount := 0
	for _, migration := range migrations {
		if !appliedMap[migration.name] {
			pendingCount++
		}
	}

	if pendingCount > 0 {
		log.Printf("Running %d pending migration(s)...", pendingCount)

		for _, migration := range migrations {
			if appliedMap[migration.name] {
				continue
			}

			log.Printf("Applying migration: %s", migration.name)
			if err := db.AutoMigrate(migration.model); err != nil {
				return fmt.Errorf("migration %s failed: %w", migration.name, err)
			}

			record := SchemaMigration{
				Name:      migration.name,
				AppliedAt: time.Now(),
			}
			if err := db.Create(&record).Error; err != nil {
				return fmt.Errorf("failed to record migration %s: %w", migration.name, err)
			}

			log.Printf("Migration %s applied successfully", migration.name)
		}

		log.Println("All model migrations completed successfully")
	} else {
		log.Println("No pending model migrations")
	}

	if err := runCustomMigrations(db, appliedMap); err != nil {
		return err
	}

	return nil
}

func runCustomMigrations(db *gorm.DB, appliedMap map[string]bool) error {
	customMigrations := []struct {
		name string
		sql  string
	}{
		{
			name: "025_add_unique_active_subscription_index",
			sql: `CREATE UNIQUE INDEX IF NOT EXISTS idx_user_subscriptions_user_active
				  ON user_subscriptions (user_id)
				  WHERE status = 'active' AND deleted_at IS NULL`,
		},
		{
			name: "035_add_user_announcement_status_unique_index",
			sql: `CREATE UNIQUE INDEX IF NOT EXISTS idx_user_announcement_status_unique
				  ON user_announcement_statuses (user_id, announcement_id)`,
		},
		{
			name: "037_add_conversation_participants_unique_index",
			sql: `CREATE UNIQUE INDEX IF NOT EXISTS idx_conversation_participants_unique
				  ON conversations (LEAST(participant1, participant2), GREATEST(participant1, participant2))
				  WHERE deleted_at IS NULL`,
		},
		{
			name: "040_migrate_jobs_company_to_company_id",
			sql: `DO $$
				  BEGIN
				      IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'jobs' AND column_name = 'company_id') THEN
				          ALTER TABLE jobs ADD COLUMN company_id uuid;
				      END IF;
				      IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'fk_jobs_company' AND table_name = 'jobs') THEN
				          ALTER TABLE jobs ADD CONSTRAINT fk_jobs_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE;
				      END IF;
				      IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_jobs_company_id') THEN
				          CREATE INDEX idx_jobs_company_id ON jobs(company_id);
				      END IF;
				      IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'jobs' AND column_name = 'company') THEN
				          ALTER TABLE jobs DROP COLUMN company;
				      END IF;
				  END $$;`,
		},
	}

	for _, migration := range customMigrations {
		if appliedMap[migration.name] {
			continue
		}

		log.Printf("Applying custom migration: %s", migration.name)
		if err := db.Exec(migration.sql).Error; err != nil {
			// Ignore if index already exists
			if !strings.Contains(err.Error(), "already exists") {
				return fmt.Errorf("custom migration %s failed: %w", migration.name, err)
			}
		}

		record := SchemaMigration{
			Name:      migration.name,
			AppliedAt: time.Now(),
		}
		if err := db.Create(&record).Error; err != nil {
			// Ignore duplicate key error (already recorded)
			if !strings.Contains(err.Error(), "duplicate key") {
				return fmt.Errorf("failed to record migration %s: %w", migration.name, err)
			}
		}

		log.Printf("Custom migration %s applied successfully", migration.name)
	}

	return nil
}

func PingDatabase(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	return sqlDB.PingContext(ctx)
}

func GetStats() (sql.DBStats, error) {
	if client == nil {
		return sql.DBStats{}, ErrDatabaseNotInitialized
	}

	sqlDB, err := client.DB()
	if err != nil {
		return sql.DBStats{}, err
	}

	return sqlDB.Stats(), nil
}
