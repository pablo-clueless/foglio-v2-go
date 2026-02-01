package handlers

import (
	"foglio/v2/src/config"
	"foglio/v2/src/database"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/services"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	service *services.ReviewService
}

func NewReviewHandler() *ReviewHandler {
	return &ReviewHandler{
		service: services.NewReviewService(database.GetDatabase()),
	}
}

func (h *ReviewHandler) CreateReview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		var payload dto.CreateReviewDto

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		review, err := h.service.CreateReview(userId, payload)
		if err != nil {
			if err.Error() == "you have already submitted a review" {
				lib.BadRequest(ctx, err.Error(), "400")
				return
			}
			lib.InternalServerError(ctx, "Internal server error: "+err.Error())
			return
		}

		lib.Success(ctx, "Review created successfully", review)
	}
}

func (h *ReviewHandler) UpdateReview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		id := ctx.Param("id")
		var payload dto.UpdateReviewDto

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		review, err := h.service.UpdateReview(id, userId, payload)
		if err != nil {
			if err.Error() == "review not found" {
				lib.NotFound(ctx, err.Error(), "404")
				return
			}
			if err.Error() == "you can only update your own review" {
				lib.Forbidden(ctx, err.Error())
				return
			}
			lib.InternalServerError(ctx, "Internal server error: "+err.Error())
			return
		}

		lib.Success(ctx, "Review updated successfully", review)
	}
}

func (h *ReviewHandler) DeleteReview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		id := ctx.Param("id")

		err := h.service.DeleteReview(id, userId)
		if err != nil {
			if err.Error() == "review not found" {
				lib.NotFound(ctx, err.Error(), "404")
				return
			}
			if err.Error() == "you can only delete your own review" {
				lib.Forbidden(ctx, err.Error())
				return
			}
			lib.InternalServerError(ctx, "Internal server error: "+err.Error())
			return
		}

		lib.Success(ctx, "Review deleted successfully", nil)
	}
}

func (h *ReviewHandler) GetReviews() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var query dto.ReviewPagination

		if err := ctx.ShouldBindQuery(&query); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		reviews, err := h.service.GetReviews(query)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error: "+err.Error())
			return
		}

		lib.Success(ctx, "Reviews fetched successfully", reviews)
	}
}

func (h *ReviewHandler) GetReview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		review, err := h.service.GetReviewById(id)
		if err != nil {
			if err.Error() == "review not found" {
				lib.NotFound(ctx, err.Error(), "404")
				return
			}
			lib.InternalServerError(ctx, "Internal server error: "+err.Error())
			return
		}

		lib.Success(ctx, "Review fetched successfully", review)
	}
}

func (h *ReviewHandler) GetMyReview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)

		review, err := h.service.GetReviewByUser(userId)
		if err != nil {
			if err.Error() == "review not found" {
				lib.NotFound(ctx, err.Error(), "404")
				return
			}
			lib.InternalServerError(ctx, "Internal server error: "+err.Error())
			return
		}

		lib.Success(ctx, "Review fetched successfully", review)
	}
}

func (h *ReviewHandler) GetAverageRating() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		average, count, err := h.service.GetAverageRating()
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error: "+err.Error())
			return
		}

		lib.Success(ctx, "Average rating fetched successfully", map[string]interface{}{
			"average_rating": average,
			"total_reviews":  count,
		})
	}
}
