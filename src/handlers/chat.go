package handlers

import (
	"foglio/v2/src/config"
	"foglio/v2/src/database"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	service *services.ChatService
}

func NewChatHandler(hub *lib.Hub) *ChatHandler {
	notificationService := services.NewNotificationService(database.GetDatabase(), hub)
	return &ChatHandler{
		service: services.NewChatService(database.GetDatabase(), hub, notificationService),
	}
}

func (h *ChatHandler) SendMessage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var payload dto.SendMessageDto
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		message, err := h.service.SendMessage(userID, payload)
		if err != nil {
			switch err {
			case services.ErrCannotMessageSelf:
				lib.BadRequest(ctx, err.Error(), "CANNOT_MESSAGE_SELF")
			case services.ErrEmptyMessage:
				lib.BadRequest(ctx, err.Error(), "EMPTY_MESSAGE")
			case services.ErrRecipientNotFound:
				lib.NotFound(ctx, err.Error(), "RECIPIENT_NOT_FOUND")
			default:
				lib.InternalServerError(ctx, "Failed to send message: "+err.Error())
			}
			return
		}

		lib.Created(ctx, "Message sent successfully", message)
	}
}

func (h *ChatHandler) GetConversations() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

		conversations, err := h.service.GetConversations(userID, page, limit)
		if err != nil {
			lib.InternalServerError(ctx, "Failed to get conversations: "+err.Error())
			return
		}

		lib.Success(ctx, "Conversations retrieved successfully", conversations)
	}
}

func (h *ChatHandler) GetConversation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		conversationID := ctx.Param("id")

		conversation, err := h.service.GetConversation(userID, conversationID)
		if err != nil {
			switch err {
			case services.ErrConversationNotFound:
				lib.NotFound(ctx, err.Error(), "CONVERSATION_NOT_FOUND")
			case services.ErrNotParticipant:
				lib.Forbidden(ctx, err.Error())
			default:
				lib.InternalServerError(ctx, "Failed to get conversation: "+err.Error())
			}
			return
		}

		lib.Success(ctx, "Conversation retrieved successfully", conversation)
	}
}

func (h *ChatHandler) GetOrCreateConversation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		otherUserID := ctx.Param("userId")

		conversation, err := h.service.GetOrCreateConversation(userID, otherUserID)
		if err != nil {
			switch err {
			case services.ErrCannotMessageSelf:
				lib.BadRequest(ctx, err.Error(), "CANNOT_MESSAGE_SELF")
			case services.ErrRecipientNotFound:
				lib.NotFound(ctx, "User not found", "USER_NOT_FOUND")
			default:
				lib.InternalServerError(ctx, "Failed to get conversation: "+err.Error())
			}
			return
		}

		lib.Success(ctx, "Conversation retrieved successfully", conversation)
	}
}

func (h *ChatHandler) GetMessages() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		conversationID := ctx.Param("id")
		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "50"))

		messages, err := h.service.GetMessages(userID, conversationID, page, limit)
		if err != nil {
			switch err {
			case services.ErrConversationNotFound:
				lib.NotFound(ctx, err.Error(), "CONVERSATION_NOT_FOUND")
			case services.ErrNotParticipant:
				lib.Forbidden(ctx, err.Error())
			default:
				lib.InternalServerError(ctx, "Failed to get messages: "+err.Error())
			}
			return
		}

		lib.Success(ctx, "Messages retrieved successfully", messages)
	}
}

func (h *ChatHandler) MarkAsRead() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		conversationID := ctx.Param("id")

		err := h.service.MarkMessagesAsRead(userID, conversationID)
		if err != nil {
			switch err {
			case services.ErrConversationNotFound:
				lib.NotFound(ctx, err.Error(), "CONVERSATION_NOT_FOUND")
			case services.ErrNotParticipant:
				lib.Forbidden(ctx, err.Error())
			default:
				lib.InternalServerError(ctx, "Failed to mark messages as read: "+err.Error())
			}
			return
		}

		lib.Success(ctx, "Messages marked as read", nil)
	}
}

func (h *ChatHandler) DeleteConversation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		conversationID := ctx.Param("id")

		err := h.service.DeleteConversation(userID, conversationID)
		if err != nil {
			switch err {
			case services.ErrConversationNotFound:
				lib.NotFound(ctx, err.Error(), "CONVERSATION_NOT_FOUND")
			case services.ErrNotParticipant:
				lib.Forbidden(ctx, err.Error())
			default:
				lib.InternalServerError(ctx, "Failed to delete conversation: "+err.Error())
			}
			return
		}

		lib.Success(ctx, "Conversation deleted successfully", nil)
	}
}

func (h *ChatHandler) GetUnreadCount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(config.AppConfig.CurrentUserId)
		if userID == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		count, err := h.service.GetUnreadCount(userID)
		if err != nil {
			lib.InternalServerError(ctx, "Failed to get unread count: "+err.Error())
			return
		}

		lib.Success(ctx, "Unread count retrieved", map[string]int64{"unread_count": count})
	}
}
