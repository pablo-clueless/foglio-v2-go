package routes

import (
	"foglio/v2/src/handlers"

	"github.com/gin-gonic/gin"
)

func JobRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	jobs := router.Group("/jobs")
	handler := handlers.NewJobHandler()

	jobs.POST("", handler.CreateJob())
	jobs.GET("", handler.GetJobs())
	jobs.GET("/:id", handler.GetJob())
	jobs.PUT("/:id", handler.UpdateJob())
	jobs.DELETE("/:id", handler.DeleteJob())
	jobs.POST("/:id/apply", handler.ApplyToJob())
	jobs.GET("/applications/user", handler.GetApplicationsByUser())
	jobs.GET("/applications/recruiter", handler.GetApplicationsByRecruiter())
	jobs.GET("/applications/job/:id", handler.GetApplicationsByJob())
	jobs.GET("/applications/:applicationId", handler.GetApplication())
	jobs.POST("/applications/:id/accept", handler.AcceptApplication())
	jobs.POST("/applications/:id/reject", handler.RejectApplication())
	jobs.POST("/applications/:id/review", handler.ReviewApplication())
	jobs.POST("/applications/:id/hire", handler.HireApplication())
	jobs.POST("/:id/comment", handler.AddComment())
	jobs.DELETE("/:id/comment", handler.DeleteComment())
	jobs.POST("/:id/reaction/:reaction", handler.AddReaction())
	jobs.DELETE("/:id/reaction", handler.RemoveReaction())

	return jobs
}
