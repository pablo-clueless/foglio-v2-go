package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "SENT"
	MessageStatusDelivered MessageStatus = "DELIVERED"
	MessageStatusRead      MessageStatus = "READ"
)

type MediaType string

const (
	MediaTypeImage    MediaType = "IMAGE"
	MediaTypeVideo    MediaType = "VIDEO"
	MediaTypeAudio    MediaType = "AUDIO"
	MediaTypeDocument MediaType = "DOCUMENT"
	MediaTypeFile     MediaType = "FILE"
)

// MessageMedia represents a media attachment in a message
type MessageMedia struct {
	ID        string    `json:"id"`
	Type      MediaType `json:"type"`
	URL       string    `json:"url"`
	FileName  string    `json:"file_name,omitempty"`
	FileSize  int64     `json:"file_size,omitempty"`
	MimeType  string    `json:"mime_type,omitempty"`
	Width     int       `json:"width,omitempty"`
	Height    int       `json:"height,omitempty"`
	Duration  int       `json:"duration,omitempty"` // For audio/video in seconds
	Thumbnail string    `json:"thumbnail,omitempty"`
}

// MessageMediaList is a slice of MessageMedia that implements SQL Value/Scan
type MessageMediaList []MessageMedia

func (m MessageMediaList) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (m *MessageMediaList) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, m)
}

type Conversation struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Participant1 uuid.UUID      `gorm:"type:uuid;not null;index" json:"participant_1"`
	Participant2 uuid.UUID      `gorm:"type:uuid;not null;index" json:"participant_2"`
	User1        User           `gorm:"foreignKey:Participant1;references:ID;constraint:OnDelete:CASCADE" json:"user_1,omitempty"`
	User2        User           `gorm:"foreignKey:Participant2;references:ID;constraint:OnDelete:CASCADE" json:"user_2,omitempty"`
	LastMessage  *Message       `gorm:"-" json:"last_message,omitempty"`
	Messages     []Message      `gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE" json:"messages,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Message struct {
	ID             uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ConversationID uuid.UUID        `gorm:"type:uuid;not null;index" json:"conversation_id"`
	Conversation   Conversation     `gorm:"foreignKey:ConversationID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	SenderID       uuid.UUID        `gorm:"type:uuid;not null;index" json:"sender_id"`
	Sender         User             `gorm:"foreignKey:SenderID;references:ID;constraint:OnDelete:CASCADE" json:"sender,omitempty"`
	RecipientID    uuid.UUID        `gorm:"type:uuid;not null;index" json:"recipient_id"`
	Recipient      User             `gorm:"foreignKey:RecipientID;references:ID;constraint:OnDelete:CASCADE" json:"recipient,omitempty"`
	Content        string           `gorm:"type:text" json:"content"`
	Media          MessageMediaList `gorm:"type:jsonb" json:"media,omitempty"`
	Status         MessageStatus    `gorm:"not null;default:'SENT'" json:"status"`
	ReadAt         *time.Time       `json:"read_at,omitempty"`
	DeliveredAt    *time.Time       `json:"delivered_at,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	DeletedAt      gorm.DeletedAt   `gorm:"index" json:"-"`
}

func (c *Conversation) IsParticipant(userID uuid.UUID) bool {
	return c.Participant1 == userID || c.Participant2 == userID
}

func (c *Conversation) GetOtherParticipant(userID uuid.UUID) uuid.UUID {
	if c.Participant1 == userID {
		return c.Participant2
	}
	return c.Participant1
}
