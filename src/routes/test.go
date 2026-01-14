package routes

import (
	"foglio/v2/src/handlers"

	"github.com/gin-gonic/gin"
)

func TestingRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	test := router.Group("/test")
	handler := handlers.NewEmailHandler()

	test.GET("/email", handler.TestEmail())

	return test
}
