package dto

import (
	"foglio/v2/src/models"
	"time"
)

type CreateJobDto struct {
	Title          string         `json:"title"`
	Company        string         `json:"company"`
	Location       string         `json:"location"`
	Description    string         `json:"description"`
	Requirements   []string       `json:"requirements" gorm:"serializer:json"`
	Salary         *models.Salary `json:"salary,omitempty" gorm:"embedded;embeddedPrefix:salary_"`
	Deadline       *time.Time     `json:"deadline,omitempty"`
	IsRemote       bool           `json:"isRemote"`
	EmploymentType string         `json:"employmentType"`
}

type UpdateJobDto struct {
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Requirements   []string       `json:"requirements" gorm:"serializer:json"`
	Salary         *models.Salary `json:"salary,omitempty" gorm:"embedded;embeddedPrefix:salary_"`
	Deadline       *time.Time     `json:"deadline,omitempty"`
	IsRemote       bool           `json:"isRemote"`
	EmploymentType string         `json:"employmentType"`
}

type JobApplicationDto struct {
	Resume         string    `json:"resume"`
	CoverLetter    string    `json:"coverLetter,omitempty"`
	Status         string    `json:"status"`
	SubmissionDate time.Time `json:"submissionDate"`
	LastUpdated    time.Time `json:"lastUpdated"`
	Notes          *string   `json:"notes,omitempty"`
}
