package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PixKey struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Key       string    `gorm:"unique;not null;index" json:"key"`
	Type      string    `gorm:"not null;index" json:"type"`
	Status    string    `gorm:"default:'active';index" json:"status"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *PixKey) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
