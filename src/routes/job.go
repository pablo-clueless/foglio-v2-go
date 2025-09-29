package routes

import "github.com/gin-gonic/gin"

func JobRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	jobs := router.Group("/jobs")

	return jobs
}
