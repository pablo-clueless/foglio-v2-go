package dto

import "foglio/v2/src/models"

type SigninDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserDto struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserDto struct {
	Username       *string                `json:"username,omitempty"`
	Phone          *string                `json:"phone,omitempty"`
	Headline       *string                `json:"headline,omitempty"`
	Location       *string                `json:"location,omitempty"`
	Summary        *string                `json:"summary,omitempty"`
	Skills         []models.Skill         `json:"skills,omitempty"`
	Projects       []models.Project       `json:"projects,omitempty"`
	Experiences    []models.Experience    `json:"experiences,omitempty"`
	Education      []models.Education     `json:"education,omitempty"`
	Certifications []models.Certification `json:"certifications,omitempty"`
	Languages      []models.Language      `json:"languages,omitempty"`
}

type ChangePasswordDto struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ResetPasswordDto struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	Token           string `json:"token"`
}
