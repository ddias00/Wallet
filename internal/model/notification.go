package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	EventType      string     `gorm:"not null;index" json:"event_type"`
	Channel        string     `gorm:"not null;index" json:"channel"`
	Message        string     `gorm:"not null" json:"message"`
	DeliveryStatus string     `gorm:"default:'pending';index" json:"delivery_status"`
	SentAt         *time.Time `json:"sent_at,omitempty"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
	FailedAt       *time.Time `json:"failed_at,omitempty"`
	ErrorMessage   string     `json:"error_message,omitempty"`
	Metadata       string     `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt      time.Time  `gorm:"index" json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}
