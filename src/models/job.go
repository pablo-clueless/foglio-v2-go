package models

import (
	"time"

	"gorm.io/gorm"
)

type Job struct {
	ID             string         `json:"id" gorm:"primaryKey"`
	Title          string         `json:"title"`
	Company        string         `json:"company"`
	Location       string         `json:"location"`
	Description    string         `json:"description"`
	Requirements   []string       `json:"requirements" gorm:"serializer:json"`
	Salary         *Salary        `json:"salary,omitempty" gorm:"embedded;embeddedPrefix:salary_"`
	PostedDate     time.Time      `json:"postedDate"`
	Deadline       *time.Time     `json:"deadline,omitempty"`
	IsRemote       bool           `json:"isRemote"`
	EmploymentType string         `json:"employmentType"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
	CreatedBy      string         `json:"created_by" gorm:"column:created_by;not null" validate:"required"`
	CreatedByUser  User           `json:"created_by_user" gorm:"foreignKey:CreatedBy;references:ID"`
}

type Salary struct {
	Min      int64  `json:"min" gorm:"column:min"`
	Max      int64  `json:"max" gorm:"column:max"`
	Currency string `json:"currency" gorm:"column:currency"`
}

type JobSearch struct {
	Keywords       []string `json:"keywords,omitempty"`
	Location       *string  `json:"location,omitempty"`
	EmploymentType []string `json:"employmentType,omitempty"`
	SalaryMin      *int64   `json:"salaryMin,omitempty"`
	SalaryMax      *int64   `json:"salaryMax,omitempty"`
	IsRemote       *bool    `json:"isRemote,omitempty"`
	SortBy         *string  `json:"sortBy,omitempty"`
	Page           *int     `json:"page,omitempty"`
	Limit          *int     `json:"limit,omitempty"`
}

type JobApplication struct {
	ID             string         `json:"id" gorm:"primaryKey"`
	JobID          string         `json:"jobId" gorm:"index"`
	ApplicantID    string         `json:"applicantId" gorm:"index"`
	Resume         string         `json:"resume"`
	CoverLetter    *string        `json:"coverLetter,omitempty"`
	Status         string         `json:"status"` // e.g. "pending", "reviewed", "rejected", "hired"
	SubmissionDate time.Time      `json:"submissionDate"`
	LastUpdated    time.Time      `json:"lastUpdated"`
	Notes          *string        `json:"notes,omitempty"`
	CreatedAt      time.Time      `json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (j *Job) BeforeCreate(tx *gorm.DB) {
	now := time.Now()
	j.PostedDate = now
	j.CreatedAt = now
	j.UpdatedAt = now
}

func (a *JobApplication) BeforeCreate(tx *gorm.DB) {
	now := time.Now()
	a.SubmissionDate = now
	a.CreatedAt = now
	a.UpdatedAt = now

	a.Status = "pending"
}

func (j *Job) AfterUpdate(tx *gorm.DB) {
	now := time.Now()
	j.UpdatedAt = now
}

func (a *JobApplication) AfterUpdate(tx *gorm.DB) {
	now := time.Now()
	a.LastUpdated = now
	a.UpdatedAt = now
}
