package routes

import (
	"foglio/v2/src/handlers"
	"foglio/v2/src/lib"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(router *gin.RouterGroup, hub *lib.Hub) *gin.RouterGroup {
	handler := handlers.NewChatHandler(hub)

	chat := router.Group("/chat")

	chat.GET("/conversations", handler.GetConversations())
	chat.GET("/conversations/:id", handler.GetConversation())
	chat.GET("/conversations/user/:userId", handler.GetOrCreateConversation())
	chat.DELETE("/conversations/:id", handler.DeleteConversation())

	chat.POST("/messages", handler.SendMessage())
	chat.GET("/conversations/:id/messages", handler.GetMessages())
	chat.PUT("/conversations/:id/read", handler.MarkAsRead())

	chat.GET("/unread", handler.GetUnreadCount())

	return chat
}
