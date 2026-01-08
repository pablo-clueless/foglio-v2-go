package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"message": "Server is healthy",
		})
	})

	return router
}