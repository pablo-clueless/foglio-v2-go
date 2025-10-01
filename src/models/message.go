package models

import "github.com/google/uuid"

type Message struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
}
