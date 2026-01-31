package services

import (
	"errors"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/models"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrConversationNotFound = errors.New("conversation not found")
	ErrMessageNotFound      = errors.New("message not found")
	ErrNotParticipant       = errors.New("you are not a participant in this conversation")
	ErrCannotMessageSelf    = errors.New("you cannot send a message to yourself")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrEmptyMessage         = errors.New("message must have content or media")
)

type ChatService struct {
	database            *gorm.DB
	hub                 *lib.Hub
	notificationService *NotificationService
}

func NewChatService(database *gorm.DB, hub *lib.Hub, notificationService *NotificationService) *ChatService {
	return &ChatService{
		database:            database,
		hub:                 hub,
		notificationService: notificationService,
	}
}

func (s *ChatService) SendMessage(senderID string, payload dto.SendMessageDto) (*models.Message, error) {
	senderUUID := uuid.Must(uuid.Parse(senderID))
	recipientUUID, err := uuid.Parse(payload.RecipientID)
	if err != nil {
		return nil, errors.New("invalid recipient ID")
	}

	if senderUUID == recipientUUID {
		return nil, ErrCannotMessageSelf
	}

	if payload.Content == "" && len(payload.Media) == 0 {
		return nil, ErrEmptyMessage
	}

	var recipient models.User
	if err := s.database.Where("id = ?", recipientUUID).First(&recipient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecipientNotFound
		}
		return nil, err
	}

	conversation, err := s.findOrCreateConversation(senderUUID, recipientUUID)
	if err != nil {
		return nil, err
	}

	var media models.MessageMediaList
	for _, m := range payload.Media {
		media = append(media, models.MessageMedia{
			ID:        m.ID,
			Type:      m.Type,
			URL:       m.URL,
			FileName:  m.FileName,
			FileSize:  m.FileSize,
			MimeType:  m.MimeType,
			Width:     m.Width,
			Height:    m.Height,
			Duration:  m.Duration,
			Thumbnail: m.Thumbnail,
		})
	}

	message := &models.Message{
		ConversationID: conversation.ID,
		SenderID:       senderUUID,
		RecipientID:    recipientUUID,
		Content:        payload.Content,
		Media:          media,
		Status:         models.MessageStatusSent,
	}

	if err := s.database.Create(message).Error; err != nil {
		return nil, err
	}

	s.database.Model(conversation).Update("updated_at", time.Now())

	s.database.Preload("Sender").First(message, "id = ?", message.ID)

	go s.sendMessageNotification(message, &recipient)

	return message, nil
}

func (s *ChatService) GetConversations(userID string, page, limit int) (*dto.ConversationListResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}

	userUUID := uuid.Must(uuid.Parse(userID))
	var conversations []models.Conversation
	var totalItems int64

	query := s.database.Model(&models.Conversation{}).
		Where("participant1 = ? OR participant2 = ?", userUUID, userUUID)

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * limit
	if err := query.
		Preload("User1").
		Preload("User2").
		Order("updated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&conversations).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.ConversationResponse, len(conversations))
	for i, conv := range conversations {
		responses[i] = s.toConversationResponse(&conv, userUUID)
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int((totalItems + int64(limit) - 1) / int64(limit))
	}

	return &dto.ConversationListResponse{
		Data:       responses,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
		Page:       page,
		Limit:      limit,
	}, nil
}

func (s *ChatService) GetConversation(userID, conversationID string) (*dto.ConversationResponse, error) {
	userUUID := uuid.Must(uuid.Parse(userID))

	var conversation models.Conversation
	if err := s.database.
		Preload("User1").
		Preload("User2").
		Where("id = ?", conversationID).
		First(&conversation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConversationNotFound
		}
		return nil, err
	}

	if !conversation.IsParticipant(userUUID) {
		return nil, ErrNotParticipant
	}

	response := s.toConversationResponse(&conversation, userUUID)
	return &response, nil
}

func (s *ChatService) GetOrCreateConversation(userID, otherUserID string) (*dto.ConversationResponse, error) {
	userUUID := uuid.Must(uuid.Parse(userID))
	otherUUID, err := uuid.Parse(otherUserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	if userUUID == otherUUID {
		return nil, ErrCannotMessageSelf
	}

	var otherUser models.User
	if err := s.database.Where("id = ?", otherUUID).First(&otherUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecipientNotFound
		}
		return nil, err
	}

	conversation, err := s.findOrCreateConversation(userUUID, otherUUID)
	if err != nil {
		return nil, err
	}

	s.database.Preload("User1").Preload("User2").First(conversation, "id = ?", conversation.ID)

	response := s.toConversationResponse(conversation, userUUID)
	return &response, nil
}

func (s *ChatService) GetMessages(userID, conversationID string, page, limit int) (*dto.MessageListResponse, error) {
	if limit <= 0 {
		limit = 50
	}
	if page <= 0 {
		page = 1
	}

	userUUID := uuid.Must(uuid.Parse(userID))

	var conversation models.Conversation
	if err := s.database.Where("id = ?", conversationID).First(&conversation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConversationNotFound
		}
		return nil, err
	}

	if !conversation.IsParticipant(userUUID) {
		return nil, ErrNotParticipant
	}

	var messages []models.Message
	var totalItems int64

	query := s.database.Model(&models.Message{}).Where("conversation_id = ?", conversationID)

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * limit
	if err := query.
		Preload("Sender").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&messages).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.MessageResponse, len(messages))
	for i, msg := range messages {
		responses[i] = s.toMessageResponse(&msg)
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int((totalItems + int64(limit) - 1) / int64(limit))
	}

	return &dto.MessageListResponse{
		Data:       responses,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
		Page:       page,
		Limit:      limit,
	}, nil
}

func (s *ChatService) MarkMessagesAsRead(userID, conversationID string) error {
	userUUID := uuid.Must(uuid.Parse(userID))

	var conversation models.Conversation
	if err := s.database.Where("id = ?", conversationID).First(&conversation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrConversationNotFound
		}
		return err
	}

	if !conversation.IsParticipant(userUUID) {
		return ErrNotParticipant
	}

	now := time.Now()

	result := s.database.Model(&models.Message{}).
		Where("conversation_id = ? AND recipient_id = ? AND status != ?", conversationID, userUUID, models.MessageStatusRead).
		Updates(map[string]interface{}{
			"status":  models.MessageStatusRead,
			"read_at": now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		otherUserID := conversation.GetOtherParticipant(userUUID)
		go s.sendReadReceipt(conversationID, userID, otherUserID.String())
	}

	return nil
}

func (s *ChatService) DeleteConversation(userID, conversationID string) error {
	userUUID := uuid.Must(uuid.Parse(userID))

	var conversation models.Conversation
	if err := s.database.Where("id = ?", conversationID).First(&conversation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrConversationNotFound
		}
		return err
	}

	if !conversation.IsParticipant(userUUID) {
		return ErrNotParticipant
	}

	return s.database.Delete(&conversation).Error
}

func (s *ChatService) GetUnreadCount(userID string) (int64, error) {
	userUUID := uuid.Must(uuid.Parse(userID))
	var count int64

	err := s.database.Model(&models.Message{}).
		Where("recipient_id = ? AND status != ?", userUUID, models.MessageStatusRead).
		Count(&count).Error

	return count, err
}

func (s *ChatService) findOrCreateConversation(user1, user2 uuid.UUID) (*models.Conversation, error) {
	var conversation models.Conversation

	p1, p2 := user1, user2
	if p1.String() > p2.String() {
		p1, p2 = p2, p1
	}

	err := s.database.Where(
		"(participant1 = ? AND participant2 = ?) OR (participant1 = ? AND participant2 = ?)",
		p1, p2, p2, p1,
	).First(&conversation).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			conversation = models.Conversation{
				Participant1: p1,
				Participant2: p2,
			}
			if err := s.database.Create(&conversation).Error; err != nil {
				return nil, err
			}
			return &conversation, nil
		}
		return nil, err
	}

	return &conversation, nil
}

func (s *ChatService) toConversationResponse(conv *models.Conversation, userID uuid.UUID) dto.ConversationResponse {
	var otherUser *models.User
	if conv.Participant1 == userID {
		otherUser = &conv.User2
	} else {
		otherUser = &conv.User1
	}

	var lastMessage models.Message
	s.database.Where("conversation_id = ?", conv.ID).Order("created_at DESC").First(&lastMessage)

	var unreadCount int64
	s.database.Model(&models.Message{}).
		Where("conversation_id = ? AND recipient_id = ? AND status != ?", conv.ID, userID, models.MessageStatusRead).
		Count(&unreadCount)

	response := dto.ConversationResponse{
		ID: conv.ID.String(),
		OtherUser: &dto.UserSummary{
			ID:       otherUser.ID.String(),
			Name:     otherUser.Name,
			Username: otherUser.Username,
			Image:    otherUser.Image,
		},
		UnreadCount: int(unreadCount),
		CreatedAt:   conv.CreatedAt,
		UpdatedAt:   conv.UpdatedAt,
	}

	if lastMessage.ID != uuid.Nil {
		msgResp := s.toMessageResponse(&lastMessage)
		response.LastMessage = &msgResp
	}

	return response
}

func (s *ChatService) toMessageResponse(msg *models.Message) dto.MessageResponse {
	response := dto.MessageResponse{
		ID:             msg.ID.String(),
		ConversationID: msg.ConversationID.String(),
		SenderID:       msg.SenderID.String(),
		RecipientID:    msg.RecipientID.String(),
		Content:        msg.Content,
		Status:         msg.Status,
		ReadAt:         msg.ReadAt,
		DeliveredAt:    msg.DeliveredAt,
		CreatedAt:      msg.CreatedAt,
	}

	if len(msg.Media) > 0 {
		response.Media = make([]dto.MediaDto, len(msg.Media))
		for i, m := range msg.Media {
			response.Media[i] = dto.MediaDto{
				ID:        m.ID,
				Type:      m.Type,
				URL:       m.URL,
				FileName:  m.FileName,
				FileSize:  m.FileSize,
				MimeType:  m.MimeType,
				Width:     m.Width,
				Height:    m.Height,
				Duration:  m.Duration,
				Thumbnail: m.Thumbnail,
			}
		}
	}

	if msg.Sender.ID != uuid.Nil {
		response.Sender = &dto.UserSummary{
			ID:       msg.Sender.ID.String(),
			Name:     msg.Sender.Name,
			Username: msg.Sender.Username,
			Image:    msg.Sender.Image,
		}
	}

	return response
}

func (s *ChatService) sendMessageNotification(message *models.Message, recipient *models.User) {
	var sender models.User
	s.database.Where("id = ?", message.SenderID).First(&sender)

	msgResp := s.toMessageResponse(message)

	if s.hub != nil {
		notification := models.Notification{
			ID:      uuid.New(),
			Title:   "New Message",
			Content: sender.Name + " sent you a message",
			Type:    models.NewMessage,
			OwnerID: recipient.ID,
			IsRead:  false,
			Data: map[string]interface{}{
				"event_type":      "new_message",
				"conversation_id": message.ConversationID.String(),
				"message_id":      message.ID.String(),
				"message":         msgResp,
				"sender_id":       sender.ID.String(),
				"sender_name":     sender.Name,
			},
		}
		s.hub.SendToUser(recipient.ID.String(), notification)
	}

	if s.notificationService != nil {
		err := s.notificationService.SendRealTimeNotification(
			recipient.ID.String(),
			"New Message",
			sender.Name+" sent you a message",
			models.NewMessage,
			map[string]interface{}{
				"conversation_id": message.ConversationID.String(),
				"message_id":      message.ID.String(),
				"sender_id":       sender.ID.String(),
				"sender_name":     sender.Name,
			},
		)
		if err != nil {
			log.Printf("Failed to send message notification: %v", err)
		}
	}
}

func (s *ChatService) sendReadReceipt(conversationID, readerID, recipientID string) {
	if s.hub != nil {
		notification := models.Notification{
			ID:      uuid.New(),
			Title:   "Messages Read",
			Content: "Your messages have been read",
			Type:    models.System,
			OwnerID: uuid.Must(uuid.Parse(recipientID)),
			IsRead:  false,
			Data: map[string]interface{}{
				"event_type":      "messages_read",
				"conversation_id": conversationID,
				"reader_id":       readerID,
			},
		}
		s.hub.SendToUser(recipientID, notification)
	}
}

func (s *ChatService) HandleWebSocketMessage(senderID string, payload map[string]interface{}) (interface{}, error) {
	action, ok := payload["action"].(string)
	if !ok {
		return nil, errors.New("action is required")
	}

	switch action {
	case "send_message":
		return s.handleSendMessage(senderID, payload)
	case "typing":
		return s.handleTyping(senderID, payload, true)
	case "stop_typing":
		return s.handleTyping(senderID, payload, false)
	case "mark_messages_read":
		return s.handleMarkMessagesRead(senderID, payload)
	default:
		return nil, errors.New("unknown action: " + action)
	}
}

func (s *ChatService) handleSendMessage(senderID string, payload map[string]interface{}) (interface{}, error) {
	recipientID, ok := payload["recipient_id"].(string)
	if !ok || recipientID == "" {
		return nil, errors.New("recipient_id is required")
	}

	content, _ := payload["content"].(string)

	var mediaList []dto.MediaDto
	if mediaRaw, ok := payload["media"].([]interface{}); ok {
		for _, m := range mediaRaw {
			if mediaMap, ok := m.(map[string]interface{}); ok {
				media := dto.MediaDto{
					Type: models.MediaType(getStringValue(mediaMap, "type")),
					URL:  getStringValue(mediaMap, "url"),
				}
				if media.Type == "" || media.URL == "" {
					continue
				}
				media.ID = getStringValue(mediaMap, "id")
				media.FileName = getStringValue(mediaMap, "file_name")
				media.MimeType = getStringValue(mediaMap, "mime_type")
				media.Thumbnail = getStringValue(mediaMap, "thumbnail")
				if size, ok := mediaMap["file_size"].(float64); ok {
					media.FileSize = int64(size)
				}
				if width, ok := mediaMap["width"].(float64); ok {
					media.Width = int(width)
				}
				if height, ok := mediaMap["height"].(float64); ok {
					media.Height = int(height)
				}
				if duration, ok := mediaMap["duration"].(float64); ok {
					media.Duration = int(duration)
				}
				mediaList = append(mediaList, media)
			}
		}
	}

	if content == "" && len(mediaList) == 0 {
		return nil, errors.New("message must have content or media")
	}

	messageDto := dto.SendMessageDto{
		RecipientID: recipientID,
		Content:     content,
		Media:       mediaList,
	}

	message, err := s.SendMessage(senderID, messageDto)
	if err != nil {
		return nil, err
	}

	return s.toMessageResponse(message), nil
}

func (s *ChatService) handleTyping(senderID string, payload map[string]interface{}, isTyping bool) (interface{}, error) {
	recipientID, ok := payload["recipient_id"].(string)
	if !ok || recipientID == "" {
		return nil, errors.New("recipient_id is required")
	}

	conversationID, _ := payload["conversation_id"].(string)

	eventType := "typing"
	if !isTyping {
		eventType = "stop_typing"
	}

	if s.hub != nil {
		notification := models.Notification{
			ID:      uuid.New(),
			Title:   "Typing",
			Content: "",
			Type:    models.System,
			OwnerID: uuid.Must(uuid.Parse(recipientID)),
			IsRead:  false,
			Data: map[string]interface{}{
				"event_type":      eventType,
				"conversation_id": conversationID,
				"user_id":         senderID,
			},
		}
		s.hub.SendToUser(recipientID, notification)
	}

	return map[string]interface{}{"sent": true}, nil
}

func (s *ChatService) handleMarkMessagesRead(senderID string, payload map[string]interface{}) (interface{}, error) {
	conversationID, ok := payload["conversation_id"].(string)
	if !ok || conversationID == "" {
		return nil, errors.New("conversation_id is required")
	}

	err := s.MarkMessagesAsRead(senderID, conversationID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"marked_read": true}, nil
}

func getStringValue(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
