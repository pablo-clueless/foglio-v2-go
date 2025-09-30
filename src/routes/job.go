package routes

import (
	"foglio/v2/src/handlers"

	"github.com/gin-gonic/gin"
)

func JobRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	jobs := router.Group("/jobs")
	handler := handlers.NewJobHandler()

	jobs.POST("/", handler.CreateJob())
	jobs.GET("/", handler.GetJobs())
	jobs.GET("/:id", handler.GetJob())
	jobs.PUT("/:id", handler.UpdateJob())
	jobs.DELETE("/:id", handler.DeleteJob())
	jobs.POST("/:id/apply", handler.ApplyToJob())

	return jobs
}
