package dto

import (
	"foglio/v2/src/models"
	"time"
)

type MediaDto struct {
	ID        string           `json:"id,omitempty"`
	Type      models.MediaType `json:"type" binding:"required,oneof=IMAGE VIDEO AUDIO DOCUMENT FILE"`
	URL       string           `json:"url" binding:"required,url"`
	FileName  string           `json:"file_name,omitempty"`
	FileSize  int64            `json:"file_size,omitempty"`
	MimeType  string           `json:"mime_type,omitempty"`
	Width     int              `json:"width,omitempty"`
	Height    int              `json:"height,omitempty"`
	Duration  int              `json:"duration,omitempty"`
	Thumbnail string           `json:"thumbnail,omitempty"`
}

type SendMessageDto struct {
	RecipientID string     `json:"recipient_id" binding:"required,uuid"`
	Content     string     `json:"content" binding:"max=5000"`
	Media       []MediaDto `json:"media,omitempty" binding:"max=10,dive"`
}

type WebSocketSendMessageDto struct {
	Action      string     `json:"action"` // "send_message", "typing", "stop_typing", "mark_read"
	RecipientID string     `json:"recipient_id,omitempty"`
	Content     string     `json:"content,omitempty"`
	Media       []MediaDto `json:"media,omitempty"`
	MessageID   string     `json:"message_id,omitempty"`
}

type MessageResponse struct {
	ID             string               `json:"id"`
	ConversationID string               `json:"conversation_id"`
	SenderID       string               `json:"sender_id"`
	Sender         *UserSummary         `json:"sender,omitempty"`
	RecipientID    string               `json:"recipient_id"`
	Recipient      *UserSummary         `json:"recipient,omitempty"`
	Content        string               `json:"content"`
	Media          []MediaDto           `json:"media,omitempty"`
	Status         models.MessageStatus `json:"status"`
	ReadAt         *time.Time           `json:"read_at,omitempty"`
	DeliveredAt    *time.Time           `json:"delivered_at,omitempty"`
	CreatedAt      time.Time            `json:"created_at"`
}

type UserSummary struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Image    *string `json:"image,omitempty"`
}

type ConversationResponse struct {
	ID          string           `json:"id"`
	OtherUser   *UserSummary     `json:"other_user"`
	LastMessage *MessageResponse `json:"last_message,omitempty"`
	UnreadCount int              `json:"unread_count"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type ConversationListResponse struct {
	Data       []ConversationResponse `json:"data"`
	TotalItems int                    `json:"total_items"`
	TotalPages int                    `json:"total_pages"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
}

type MessageListResponse struct {
	Data       []MessageResponse `json:"data"`
	TotalItems int               `json:"total_items"`
	TotalPages int               `json:"total_pages"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
}

type ChatQueryParams struct {
	Page   int    `json:"page" form:"page"`
	Limit  int    `json:"limit" form:"limit"`
	Before string `json:"before" form:"before"` // Message ID to fetch messages before
	After  string `json:"after" form:"after"`   // Message ID to fetch messages after
}

type WebSocketMessageEvent struct {
	Type           string           `json:"type"` // "new_message", "message_read", "message_delivered", "typing", "stop_typing"
	Message        *MessageResponse `json:"message,omitempty"`
	UserID         string           `json:"user_id,omitempty"`
	ConversationID string           `json:"conversation_id,omitempty"`
}

type WebSocketResponse struct {
	Success bool        `json:"success"`
	Type    string      `json:"type"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
