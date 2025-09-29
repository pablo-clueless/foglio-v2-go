package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID             uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name           string          `gorm:"not null" json:"name"`
	Username       string          `gorm:"not null" json:"username"`
	Email          string          `gorm:"uniqueIndex;not null" json:"email"`
	Password       string          `json:"password"`
	Phone          *string         `json:"phone,omitempty"`
	Headline       *string         `json:"headline,omitempty"`
	Location       *string         `json:"location,omitempty"`
	Image          *string         `json:"image,omitempty"`
	Summary        string          `gorm:"not null" json:"summary"`
	Skills         []Skill         `gorm:"foreignKey:UserID" json:"skills,omitempty"`
	Projects       []Project       `gorm:"foreignKey:UserID" json:"projects,omitempty"`
	Experiences    []Experience    `gorm:"foreignKey:UserID" json:"experiences,omitempty"`
	Education      []Education     `gorm:"foreignKey:UserID" json:"education,omitempty"`
	Certifications []Certification `gorm:"foreignKey:UserID" json:"certifications,omitempty"`
	Languages      []Language      `gorm:"foreignKey:UserID" json:"languages,omitempty"`
	IsRecruiter    bool            `json:"isRecruiter"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
	Verified       bool            `json:"verified"`
	Otp            string          `json:"otp"`
	Company        *Company        `gorm:"foreignKey:CompanyID;references:ID" json:"company,omitempty"` // optional, only for recruiters
}

type Company struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID      *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"` // FK to User (recruiter)
	Name        string     `gorm:"not null" json:"name"`
	Industry    *string    `json:"industry,omitempty"`
	Size        *string    `json:"size,omitempty"` // e.g. "1-10", "11-50", "51-200", "201-500", "500+"
	Website     *string    `json:"website,omitempty"`
	Logo        *string    `json:"logo,omitempty"`
	Description *string    `json:"description,omitempty"`
	Location    *string    `json:"location,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Project struct {
	ID          uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID      uuid.UUID          `gorm:"type:uuid;not null" json:"user_id"`
	Title       string             `gorm:"not null" json:"title"`
	Description string             `gorm:"not null" json:"description"`
	Image       *string            `json:"image,omitempty"`
	URL         *string            `json:"url,omitempty"`
	Stack       []ProjectStack     `gorm:"foreignKey:ProjectID" json:"stack,omitempty"`
	StartDate   *time.Time         `json:"start_date,omitempty"`
	EndDate     *time.Time         `json:"end_date,omitempty"`
	Highlights  []ProjectHighlight `gorm:"foreignKey:ProjectID" json:"highlights,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type Experience struct {
	ID           uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID       uuid.UUID             `gorm:"type:uuid;not null" json:"user_id"`
	CompanyName  string                `gorm:"not null" json:"company_name"`
	Location     *string               `json:"location,omitempty"`
	Role         string                `gorm:"not null" json:"role"`
	Description  string                `gorm:"not null" json:"description"`
	StartDate    time.Time             `gorm:"not null" json:"start_date"`
	EndDate      *time.Time            `json:"end_date,omitempty"`
	Highlights   []ExperienceHighlight `gorm:"foreignKey:ExperienceID" json:"highlights,omitempty"`
	Technologies []ExperienceTech      `gorm:"foreignKey:ExperienceID" json:"technologies,omitempty"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
}

type Education struct {
	ID          uuid.UUID            `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID      uuid.UUID            `gorm:"type:uuid;not null" json:"user_id"`
	Institution string               `gorm:"not null" json:"institution"`
	Degree      string               `gorm:"not null" json:"degree"`
	Field       string               `gorm:"not null" json:"field"`
	Location    *string              `json:"location,omitempty"`
	StartDate   time.Time            `gorm:"not null" json:"start_date"`
	EndDate     *time.Time           `json:"end_date,omitempty"`
	GPA         *float64             `json:"gpa,omitempty"`
	Highlights  []EducationHighlight `gorm:"foreignKey:EducationID" json:"highlights,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

type Certification struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	Name         string     `gorm:"not null" json:"name"`
	Issuer       string     `gorm:"not null" json:"issuer"`
	IssueDate    time.Time  `gorm:"not null" json:"issue_date"`
	ExpiryDate   *time.Time `json:"expiry_date,omitempty"`
	CredentialID *string    `json:"credential_id,omitempty"`
	URL          *string    `json:"url,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type Language struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Name        string    `gorm:"not null" json:"name"`
	Proficiency string    `gorm:"not null" json:"proficiency"` // "Basic" | "Intermediate" | "Advanced" | "Native"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Skill struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProjectStack struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null" json:"project_id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProjectHighlight struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null" json:"project_id"`
	Text      string    `gorm:"not null" json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ExperienceHighlight struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ExperienceID uuid.UUID `gorm:"type:uuid;not null" json:"experience_id"`
	Text         string    `gorm:"not null" json:"text"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ExperienceTech struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ExperienceID uuid.UUID `gorm:"type:uuid;not null" json:"experience_id"`
	Name         string    `gorm:"not null" json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type EducationHighlight struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EducationID uuid.UUID `gorm:"type:uuid;not null" json:"education_id"`
	Text        string    `gorm:"not null" json:"text"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	return nil
}

func (c *Company) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	return nil
}

func (u *User) AfterUpdate(tx *gorm.DB) error {
	now := time.Now()
	u.UpdatedAt = now
	return nil
}

func (u *User) BeforeDelete(tx *gorm.DB) error {
	if err := tx.Where("user_id = ?", u.ID).Delete(&Skill{}).Error; err != nil {
		return err
	}
	if err := tx.Where("user_id = ?", u.ID).Delete(&Language{}).Error; err != nil {
		return err
	}
	if err := tx.Where("user_id = ?", u.ID).Delete(&Certification{}).Error; err != nil {
		return err
	}
	if err := tx.Where("user_id = ?", u.ID).Delete(&Education{}).Error; err != nil {
		return err
	}
	if err := tx.Where("user_id = ?", u.ID).Delete(&Experience{}).Error; err != nil {
		return err
	}
	if err := tx.Where("user_id = ?", u.ID).Delete(&Project{}).Error; err != nil {
		return err
	}
	return nil
}
