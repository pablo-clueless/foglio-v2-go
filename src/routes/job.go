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
	jobs.GET("/applications/:id", handler.GetApplicationsByJob())
	jobs.POST("/applications/:id/accept", handler.AcceptApplication())
	jobs.POST("/applications/:id/reject", handler.RejectApplication())
	jobs.POST("/:id/comment", handler.AddComment())
	jobs.DELETE("/:id/comment", handler.DeleteComment())
	jobs.POST("/:id/reaction/:reaction", handler.AddReaction())
	jobs.DELETE("/:id/reaction", handler.RemoveReaction())

	return jobs
}
