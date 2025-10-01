package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	Content   string    `gorm:"not null" json:"content"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	OwnerID   uuid.UUID `json:"owner_id" gorm:"type:uuid;not null;index"`
	Owner     User      `json:"owner" gorm:"foreignKey:OwnerID;references:ID;constraint:OnDelete:CASCADE"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) {
	now := time.Now()
	n.CreatedAt = now
	n.UpdatedAt = now
}

func (n *Notification) BeforeUpdate(tx *gorm.DB) {
	now := time.Now()
	n.UpdatedAt = now
}
