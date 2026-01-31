package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	ID                  uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name                string             `gorm:"not null" json:"name"`
	Username            string             `gorm:"uniqueIndex;not null" json:"username"`
	Email               string             `gorm:"uniqueIndex;not null" json:"email"`
	Password            string             `gorm:"null" json:"-"`                   // Nullable for OAuth users
	Provider            string             `gorm:"default:'local'" json:"provider"` // local, google, github
	ProviderID          string             `gorm:"null" json:"-"`                   // Provider's user ID
	Role                *string            `json:"role"`
	Headline            *string            `json:"headline"`
	Phone               *string            `gorm:"index" json:"phone"`
	Location            *string            `gorm:"index" json:"location"`
	Image               *string            `json:"image"`
	Domain              *Domain            `gorm:"type:jsonb;serializer:json" json:"domain,omitempty"`
	Portfolio           *Portfolio         `gorm:"foreignKey:UserID" json:"portfolio,omitempty"`
	Summary             *string            `gorm:"null" json:"summary"`
	SocialMedia         *SocialMedia       `gorm:"type:jsonb;serializer:json" json:"social_media,omitempty"`
	CompanyID           *uuid.UUID         `gorm:"type:uuid;index" json:"company_id,omitempty"`
	Company             *Company           `gorm:"foreignKey:CompanyID;references:ID" json:"company,omitempty"`
	CurrentSubscription *UserSubscription  `gorm:"foreignKey:UserID" json:"current_subscription,omitempty"`
	SubscriptionHistory []UserSubscription `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"subscription_history,omitempty"`
	Skills              pq.StringArray     `gorm:"type:text[]" json:"skills"`
	Projects            []Project          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"projects,"`
	Experiences         []Experience       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"experiences,"`
	Education           []Education        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"education,"`
	Certifications      []Certification    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"certifications,"`
	Languages           []Language         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"languages,"`
	IsAdmin             bool               `json:"is_admin"`
	IsRecruiter         bool               `json:"is_recruiter"`
	IsPremium           bool               `json:"is_premium"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
	DeletedAt           gorm.DeletedAt     `gorm:"index" json:"-"`
	Verified            bool               `json:"verified"`
	Otp                 string             `json:"otp"`
}

type Company struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Industry    *string        `json:"industry,omitempty"`
	Size        *string        `json:"size,omitempty"`
	Website     *string        `json:"website,omitempty"`
	Logo        *string        `json:"logo,omitempty"`
	Description *string        `json:"description,omitempty"`
	Location    *string        `json:"location,omitempty"`
	Users       []User         `gorm:"foreignKey:CompanyID" json:"users,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Project struct {
	ID          uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID      uuid.UUID          `gorm:"type:uuid;not null;index" json:"user_id"`
	Title       string             `gorm:"not null" json:"title"`
	Description string             `gorm:"not null" json:"description"`
	Image       *string            `json:"image,omitempty"`
	URL         *string            `json:"url,omitempty"`
	Stack       []ProjectStack     `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"stack,omitempty"`
	StartDate   *time.Time         `json:"start_date,omitempty"`
	EndDate     *time.Time         `json:"end_date,omitempty"`
	Highlights  []ProjectHighlight `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"highlights,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	DeletedAt   gorm.DeletedAt     `gorm:"index" json:"-"`
}

type Experience struct {
	ID           uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID       uuid.UUID             `gorm:"type:uuid;not null;index" json:"user_id"`
	CompanyName  string                `gorm:"not null" json:"company_name"`
	Location     *string               `json:"location,omitempty"`
	Role         string                `gorm:"not null" json:"role"`
	Description  string                `gorm:"not null" json:"description"`
	StartDate    time.Time             `gorm:"not null" json:"start_date"`
	EndDate      *time.Time            `json:"end_date,omitempty"`
	Highlights   []ExperienceHighlight `gorm:"foreignKey:ExperienceID;constraint:OnDelete:CASCADE" json:"highlights,omitempty"`
	Technologies []ExperienceTech      `gorm:"foreignKey:ExperienceID;constraint:OnDelete:CASCADE" json:"technologies,omitempty"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
	DeletedAt    gorm.DeletedAt        `gorm:"index" json:"-"`
}

type Education struct {
	ID          uuid.UUID            `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID      uuid.UUID            `gorm:"type:uuid;not null;index" json:"user_id"`
	Institution string               `gorm:"not null" json:"institution"`
	Degree      string               `gorm:"not null" json:"degree"`
	Field       string               `gorm:"not null" json:"field"`
	Location    *string              `json:"location,omitempty"`
	StartDate   time.Time            `gorm:"not null" json:"start_date"`
	EndDate     *time.Time           `json:"end_date,omitempty"`
	GPA         *float64             `json:"gpa,omitempty"`
	Highlights  []EducationHighlight `gorm:"foreignKey:EducationID;constraint:OnDelete:CASCADE" json:"highlights,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	DeletedAt   gorm.DeletedAt       `gorm:"index" json:"-"`
}

type Certification struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Name         string         `gorm:"not null" json:"name"`
	Issuer       string         `gorm:"not null" json:"issuer"`
	IssueDate    time.Time      `gorm:"not null" json:"issue_date"`
	ExpiryDate   *time.Time     `json:"expiry_date,omitempty"`
	CredentialID *string        `json:"credential_id,omitempty"`
	URL          *string        `json:"url,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Language struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Name        string         `gorm:"not null" json:"name"`
	Proficiency string         `gorm:"not null" json:"proficiency"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type ProjectStack struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ProjectID uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	Name      string         `gorm:"not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ProjectHighlight struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ProjectID uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	Text      string         `gorm:"not null" json:"text"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ExperienceHighlight struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ExperienceID uuid.UUID      `gorm:"type:uuid;not null;index" json:"experience_id"`
	Text         string         `gorm:"not null" json:"text"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type ExperienceTech struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ExperienceID uuid.UUID      `gorm:"type:uuid;not null;index" json:"experience_id"`
	Name         string         `gorm:"not null" json:"name"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type EducationHighlight struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EducationID uuid.UUID      `gorm:"type:uuid;not null;index" json:"education_id"`
	Text        string         `gorm:"not null" json:"text"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type SocialMedia struct {
	LinkedIn  *string `json:"linkedin,omitempty"`
	GitHub    *string `json:"github,omitempty"`
	Twitter   *string `json:"twitter,omitempty"`
	Instagram *string `json:"instagram,omitempty"`
	Facebook  *string `json:"facebook,omitempty"`
	Medium    *string `json:"medium,omitempty"`
	YouTube   *string `json:"youtube,omitempty"`
	Blog      *string `json:"blog,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

func (c *Company) BeforeCreate(tx *gorm.DB) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Company) BeforeUpdate(tx *gorm.DB) error {
	c.UpdatedAt = time.Now()
	return nil
}

func (sm *SocialMedia) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), sm)
}

func (sm SocialMedia) Value() (driver.Value, error) {
	return json.Marshal(sm)
}

func (u *User) IsSubscriptionActive() bool {
	if u.CurrentSubscription == nil {
		return false
	}
	return u.CurrentSubscription.IsActive &&
		u.CurrentSubscription.Status == "active" &&
		u.CurrentSubscription.CurrentPeriodEnd.After(time.Now())
}

func (u *User) GetSubscriptionTier() SubscriptionTier {
	if u.IsSubscriptionActive() && u.CurrentSubscription.Subscription != nil {
		return u.CurrentSubscription.Subscription.Tier
	}
	return TierFree
}

func (u *User) CanAddProject() bool {
	if !u.IsSubscriptionActive() {
		return len(u.Projects) < 3
	}

	if u.CurrentSubscription.Subscription != nil {
		return len(u.Projects) < u.CurrentSubscription.Subscription.MaxProjects
	}
	return false
}

func (u *User) CanAddExperience() bool {
	if !u.IsSubscriptionActive() {
		return len(u.Experiences) < 5
	}

	if u.CurrentSubscription.Subscription != nil {
		return len(u.Experiences) < u.CurrentSubscription.Subscription.MaxExperiences
	}
	return false
}

func (u *User) IsInTrialPeriod() bool {
	if u.CurrentSubscription == nil || u.CurrentSubscription.TrialEnd == nil {
		return false
	}
	return time.Now().Before(*u.CurrentSubscription.TrialEnd)
}

func (u *User) CanUseCustomDomain() bool {
	if u.IsPremium {
		return true
	}

	tier := u.GetSubscriptionTier()
	return tier == TierBasic || tier == TierPremium || tier == TierBusiness
}
