package services

import (
	"errors"
	"foglio/v2/src/dto"
	"foglio/v2/src/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReviewService struct {
	database *gorm.DB
}

func NewReviewService(database *gorm.DB) *ReviewService {
	return &ReviewService{
		database: database,
	}
}

func (s *ReviewService) CreateReview(userId string, payload dto.CreateReviewDto) (*models.Review, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	var existingReview models.Review
	if err := s.database.Where("user_id = ?", userUUID).First(&existingReview).Error; err == nil {
		return nil, errors.New("you have already submitted a review")
	}

	review := &models.Review{
		UserID:  userUUID,
		Comment: payload.Comment,
		Rating:  payload.Rating,
	}

	if err := s.database.Create(review).Error; err != nil {
		return nil, err
	}

	if err := s.database.Preload("User").First(review, "id = ?", review.ID).Error; err != nil {
		return nil, err
	}

	return review, nil
}

func (s *ReviewService) UpdateReview(id string, userId string, payload dto.UpdateReviewDto) (*models.Review, error) {
	review, err := s.GetReviewById(id)
	if err != nil {
		return nil, err
	}

	if review.UserID.String() != userId {
		return nil, errors.New("you can only update your own review")
	}

	if payload.Comment != nil {
		review.Comment = *payload.Comment
	}
	if payload.Rating != nil {
		if *payload.Rating < 0 || *payload.Rating > 5 {
			return nil, errors.New("rating must be between 0 and 5")
		}
		review.Rating = *payload.Rating
	}

	if err := s.database.Save(review).Error; err != nil {
		return nil, err
	}

	return review, nil
}

func (s *ReviewService) DeleteReview(id string, userId string) error {
	review, err := s.GetReviewById(id)
	if err != nil {
		return err
	}

	if review.UserID.String() != userId {
		return errors.New("you can only delete your own review")
	}

	return s.database.Delete(review).Error
}

func (s *ReviewService) GetReviews(params dto.ReviewPagination) (*dto.PaginatedResponse[models.Review], error) {
	params = normalizeReviewQuery(params)

	var reviews []models.Review
	var totalItems int64

	query := s.database.Model(&models.Review{}).Preload("User")

	if params.Rating != nil && *params.Rating >= 0 && *params.Rating <= 5 {
		query = query.Where("rating = ?", *params.Rating)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.Limit
	if err := query.Offset(offset).Limit(params.Limit).Order("created_at DESC").Find(&reviews).Error; err != nil {
		return nil, err
	}

	totalPages := int(totalItems) / params.Limit
	if int(totalItems)%params.Limit != 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[models.Review]{
		Data:       reviews,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
		Page:       params.Page,
		Limit:      params.Limit,
	}, nil
}

func (s *ReviewService) GetReviewById(id string) (*models.Review, error) {
	var review models.Review
	if err := s.database.Preload("User").First(&review, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("review not found")
		}
		return nil, err
	}
	return &review, nil
}

func (s *ReviewService) GetReviewByUser(userId string) (*models.Review, error) {
	var review models.Review
	if err := s.database.Preload("User").First(&review, "user_id = ?", userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("review not found")
		}
		return nil, err
	}
	return &review, nil
}

func (s *ReviewService) GetAverageRating() (float64, int64, error) {
	var result struct {
		Average float64
		Count   int64
	}

	err := s.database.Model(&models.Review{}).
		Select("COALESCE(AVG(rating), 0) as average, COUNT(*) as count").
		Scan(&result).Error

	if err != nil {
		return 0, 0, err
	}

	return result.Average, result.Count, nil
}

func normalizeReviewQuery(q dto.ReviewPagination) dto.ReviewPagination {
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	return q
}
