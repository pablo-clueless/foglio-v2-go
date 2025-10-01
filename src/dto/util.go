package dto

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type UserPagination struct {
	Pagination
	Username string `json:"username"`
	Language string `json:"language"`
	Location string `json:"location"`
	Skill    string `json:"skill"`
}

type JobPagination struct {
	Pagination
	Company        string `json:"company"`
	Location       string `json:"location"`
	Salary         string `json:"salary"`
	PostedDate     string `json:"posted_date"`
	EmploymentType string `json:"employment_type"`
	Requirement    string `json:"requirement"`
}

type PaginatedResponse[T any] struct {
	Data       []T `json:"data"`
	Limit      int `json:"limit"`
	Page       int `json:"page"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}
