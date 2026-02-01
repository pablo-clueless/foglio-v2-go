package routes

import (
	"foglio/v2/src/handlers"

	"github.com/gin-gonic/gin"
)

func ReviewRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	reviews := router.Group("/reviews")
	handler := handlers.NewReviewHandler()

	reviews.POST("", handler.CreateReview())
	reviews.GET("", handler.GetReviews())
	reviews.GET("/me", handler.GetMyReview())
	reviews.GET("/stats", handler.GetAverageRating())
	reviews.GET("/:id", handler.GetReview())
	reviews.PUT("/:id", handler.UpdateReview())
	reviews.DELETE("/:id", handler.DeleteReview())

	return reviews
}
