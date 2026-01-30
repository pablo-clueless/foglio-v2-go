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
	IsRemote       bool           `json:"is_remote"`
	EmploymentType string         `json:"employment_type"`
}

type UpdateJobDto struct {
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Requirements   []string       `json:"requirements" gorm:"serializer:json"`
	Salary         *models.Salary `json:"salary,omitempty" gorm:"embedded;embeddedPrefix:salary_"`
	Deadline       *time.Time     `json:"deadline,omitempty"`
	IsRemote       bool           `json:"is_remote"`
	EmploymentType string         `json:"employment_type"`
}

type JobApplicationDto struct {
	Resume         string    `json:"resume"`
	CoverLetter    string    `json:"coverLetter,omitempty"`
	Status         string    `json:"status"`
	SubmissionDate time.Time `json:"submissionDate"`
	LastUpdated    time.Time `json:"lastUpdated"`
	Notes          *string   `json:"notes,omitempty"`
}

type JobSearch struct {
	Keywords       []string `json:"keywords,omitempty"`
	Location       *string  `json:"location,omitempty"`
	EmploymentType []string `json:"employment_type,omitempty"`
	SalaryMin      *int64   `json:"salary_min,omitempty"`
	SalaryMax      *int64   `json:"salary_max,omitempty"`
	IsRemote       *bool    `json:"is_remote,omitempty"`
	SortBy         *string  `json:"sortBy,omitempty"`
	Page           *int     `json:"page,omitempty"`
	Limit          *int     `json:"limit,omitempty"`
}

type ApplicationStatusDto struct {
	Reason *string `json:"reason,omitempty"`
}

type CommentDto struct {
	Content string `json:"content"`
}

type JobApplicationPagination struct {
	Pagination
	SubmissionDate *string `json:"submission_date"`
	Status         *string `json:"status"`
}
