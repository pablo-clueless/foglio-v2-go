package dto

type CreateReviewDto struct {
	Comment string `json:"comment" binding:"required"`
	Rating  int    `json:"rating" binding:"required,min=0,max=5"`
}

type UpdateReviewDto struct {
	Comment *string `json:"comment,omitempty"`
	Rating  *int    `json:"rating,omitempty" binding:"omitempty,min=0,max=5"`
}

type ReviewPagination struct {
	Pagination
	Rating *int `json:"rating" form:"rating"`
}
