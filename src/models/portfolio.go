package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PortfolioStatus string

const (
	PortfolioStatusDraft     PortfolioStatus = "draft"
	PortfolioStatusPublished PortfolioStatus = "published"
	PortfolioStatusArchived  PortfolioStatus = "archived"
)

type Portfolio struct {
	ID         uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID     uuid.UUID          `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	Title      string             `gorm:"not null" json:"title"`
	Slug       string             `gorm:"uniqueIndex;not null" json:"slug"`
	Tagline    *string            `json:"tagline,omitempty"`
	Bio        *string            `gorm:"type:text" json:"bio,omitempty"`
	CoverImage *string            `json:"cover_image,omitempty"`
	Logo       *string            `json:"logo,omitempty"`
	Template   string             `gorm:"not null;default:'default'" json:"template"`
	Theme      *PortfolioTheme    `gorm:"type:jsonb;serializer:json" json:"theme,omitempty"`
	CustomCSS  *string            `gorm:"type:text" json:"custom_css,omitempty"`
	Status     PortfolioStatus    `gorm:"not null;default:'draft'" json:"status"`
	IsPublic   bool               `gorm:"not null;default:true" json:"is_public"`
	ViewCount  int                `gorm:"not null;default:0" json:"view_count"`
	Sections   []PortfolioSection `gorm:"foreignKey:PortfolioID;constraint:OnDelete:CASCADE" json:"sections,omitempty"`
	SEO        *PortfolioSEO      `gorm:"type:jsonb;serializer:json" json:"seo,omitempty"`
	Settings   *PortfolioSettings `gorm:"type:jsonb;serializer:json" json:"settings,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	DeletedAt  gorm.DeletedAt     `gorm:"index" json:"-"`
}

type PortfolioSection struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PortfolioID uuid.UUID      `gorm:"type:uuid;not null;index" json:"portfolio_id"`
	Title       string         `gorm:"not null" json:"title"`
	Type        string         `gorm:"not null" json:"type"` // hero, about, projects, experience, skills, contact, custom
	Content     *string        `gorm:"type:text" json:"content,omitempty"`
	Settings    *string        `gorm:"type:jsonb" json:"settings,omitempty"`
	SortOrder   int            `gorm:"not null;default:0" json:"sort_order"`
	IsVisible   bool           `gorm:"not null;default:true" json:"is_visible"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type PortfolioTheme struct {
	PrimaryColor    string `json:"primary_color,omitempty"`
	SecondaryColor  string `json:"secondary_color,omitempty"`
	AccentColor     string `json:"accent_color,omitempty"`
	TextColor       string `json:"text_color,omitempty"`
	BackgroundColor string `json:"background_color,omitempty"`
	FontFamily      string `json:"font_family,omitempty"`
	FontSize        string `json:"font_size,omitempty"`
}

type PortfolioSEO struct {
	MetaTitle       *string `json:"meta_title,omitempty"`
	MetaDescription *string `json:"meta_description,omitempty"`
	MetaKeywords    *string `json:"meta_keywords,omitempty"`
	OgImage         *string `json:"og_image,omitempty"`
	Canonical       *string `json:"canonical,omitempty"`
}

type PortfolioSettings struct {
	ShowProjects       bool `json:"show_projects"`
	ShowExperiences    bool `json:"show_experiences"`
	ShowEducation      bool `json:"show_education"`
	ShowSkills         bool `json:"show_skills"`
	ShowCertifications bool `json:"show_certifications"`
	ShowContact        bool `json:"show_contact"`
	ShowSocialLinks    bool `json:"show_social_links"`
	EnableAnalytics    bool `json:"enable_analytics"`
	EnableComments     bool `json:"enable_comments"`
}

func (pt *PortfolioTheme) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), pt)
}

func (pt PortfolioTheme) Value() (driver.Value, error) {
	return json.Marshal(pt)
}

func (ps *PortfolioSEO) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), ps)
}

func (ps PortfolioSEO) Value() (driver.Value, error) {
	return json.Marshal(ps)
}

func (s *PortfolioSettings) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), s)
}

func (s PortfolioSettings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (p *Portfolio) BeforeCreate(tx *gorm.DB) error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Portfolio) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}

func (ps *PortfolioSection) BeforeCreate(tx *gorm.DB) error {
	ps.CreatedAt = time.Now()
	ps.UpdatedAt = time.Now()
	return nil
}

func (ps *PortfolioSection) BeforeUpdate(tx *gorm.DB) error {
	ps.UpdatedAt = time.Now()
	return nil
}
