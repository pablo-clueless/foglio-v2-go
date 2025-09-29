package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Job struct {
	ID             uuid.UUID        `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title          string           `json:"title" gorm:"not null"`
	Company        string           `json:"company" gorm:"not null"`
	Location       string           `json:"location" gorm:"not null"`
	Description    string           `json:"description" gorm:"not null"`
	Requirements   []string         `json:"requirements" gorm:"serializer:json"`
	Salary         *Salary          `json:"salary,omitempty" gorm:"embedded;embeddedPrefix:salary_"`
	PostedDate     time.Time        `json:"posted_date" gorm:"not null"`
	Deadline       *time.Time       `json:"deadline,omitempty"`
	IsRemote       bool             `json:"is_remote" gorm:"default:false"`
	EmploymentType string           `json:"employment_type" gorm:"not null"`
	CreatedBy      uuid.UUID        `json:"created_by" gorm:"type:uuid;not null;index"`
	CreatedByUser  User             `json:"created_by_user" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnDelete:CASCADE"`
	Applications   []JobApplication `json:"applications,omitempty" gorm:"foreignKey:JobID;constraint:OnDelete:CASCADE"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	DeletedAt      gorm.DeletedAt   `json:"-" gorm:"index"`
}

type Salary struct {
	Min      int64  `json:"min" gorm:"column:min"`
	Max      int64  `json:"max" gorm:"column:max"`
	Currency string `json:"currency" gorm:"column:currency"`
}

type JobApplication struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	JobID          uuid.UUID      `json:"job_id" gorm:"type:uuid;not null;index"`
	Job            Job            `json:"job" gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE"`
	ApplicantID    uuid.UUID      `json:"applicantId" gorm:"type:uuid;not null;index"`
	Applicant      User           `json:"applicant" gorm:"foreignKey:ApplicantID;references:ID;constraint:OnDelete:CASCADE"`
	Resume         string         `json:"resume" gorm:"not null"`
	CoverLetter    *string        `json:"coverLetter,omitempty"`
	Status         string         `json:"status" gorm:"not null;default:'pending'"`
	SubmissionDate time.Time      `json:"submission_ate" gorm:"not null"`
	LastUpdated    time.Time      `json:"last_updated" gorm:"not null"`
	Notes          *string        `json:"notes,omitempty"`
	CreatedAt      time.Time      `json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (j *Job) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	j.PostedDate = now
	j.CreatedAt = now
	j.UpdatedAt = now
	return nil
}

func (j *Job) BeforeUpdate(tx *gorm.DB) error {
	j.UpdatedAt = time.Now()
	return nil
}

func (a *JobApplication) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if a.Status == "" {
		a.Status = "pending"
	}
	a.SubmissionDate = now
	a.LastUpdated = now
	a.CreatedAt = now
	a.UpdatedAt = now
	return nil
}

func (a *JobApplication) BeforeUpdate(tx *gorm.DB) error {
	now := time.Now()
	a.LastUpdated = now
	a.UpdatedAt = now
	return nil
}
