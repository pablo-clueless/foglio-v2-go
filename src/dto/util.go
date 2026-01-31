package dto

type Pagination struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

type UserPagination struct {
	Pagination
	Query    *string `json:"query" form:"query"`
	UserType *string `json:"user_type" form:"user_type"`
}

type JobPagination struct {
	Pagination
	Company        string `json:"company" form:"company"`
	Location       string `json:"location" form:"location"`
	Salary         string `json:"salary" form:"salary"`
	PostedDate     string `json:"posted_date" form:"posted_date"`
	EmploymentType string `json:"employment_type" form:"employment_type"`
	Requirement    string `json:"requirement" form:"requirement"`
}

type PaginatedResponse[T any] struct {
	Data       []T `json:"data"`
	Limit      int `json:"limit"`
	Page       int `json:"page"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}
