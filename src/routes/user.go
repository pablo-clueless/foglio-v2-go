package routes

import (
	"foglio/v2/src/handlers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	users := router.Group("/users")
	handler := handlers.NewUserHandler()

	users.GET("", handler.GetUsers())
	users.GET("/:id", handler.GetUser())
	users.PUT("/:id", handler.UpdateUser())
	users.PUT("/:id/avatar", handler.UpdateAvatar())
	users.DELETE("/:id", handler.DeleteUser())

	user := router.Group("/user")
	user.GET("/profile", handler.GetMe())

	return users
}
