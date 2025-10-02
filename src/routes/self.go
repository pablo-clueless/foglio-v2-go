package routes

import (
	"foglio/v2/src/handlers"

	"github.com/gin-gonic/gin"
)

func SelfRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	self := router.Group("/")
	user := handlers.NewUserHandler()
	job := handlers.NewJobHandler()

	self.GET("/me", user.GetMe())
	self.GET("/me/jobs", job.GetJobsByUser())

	return self
}
