package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"foglio/v2/src/config"
	"foglio/v2/src/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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

		database, err := gorm.Open(postgres.Open(url), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
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

		if !config.AppConfig.IsDevMode {
			if err := runMigrations(database); err != nil {
				log.Printf("Failed to run migrations: %v", err)
				initErr = err
				return
			}
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

	return nil
}

func runMigrations(db *gorm.DB) error {
	migrations := []struct {
		name  string
		model interface{}
	}{
		{"001_create_companies", &models.Company{}},
		{"002_create_users", &models.User{}},
		{"003_create_skills", &models.Skill{}},
		{"004_create_languages", &models.Language{}},
		{"005_create_certifications", &models.Certification{}},
		{"006_create_education", &models.Education{}},
		{"007_create_education_highlights", &models.EducationHighlight{}},
		{"008_create_experiences", &models.Experience{}},
		{"009_create_experience_highlights", &models.ExperienceHighlight{}},
		{"010_create_experience_tech", &models.ExperienceTech{}},
		{"011_create_projects", &models.Project{}},
		{"012_create_project_stacks", &models.ProjectStack{}},
		{"013_create_project_highlights", &models.ProjectHighlight{}},
		{"014_create_jobs", &models.Job{}},
		{"015_create_job_applications", &models.JobApplication{}},
	}

	for _, migration := range migrations {
		log.Printf("Running migration: %s", migration.name)

		if err := db.AutoMigrate(migration.model); err != nil {
			return fmt.Errorf("migration %s failed: %w", migration.name, err)
		}

		log.Printf("Migration %s completed successfully", migration.name)
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
