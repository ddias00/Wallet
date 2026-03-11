package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FromUserID  *uuid.UUID `gorm:"type:uuid;index" json:"from_user_id"`
	ToUserID    *uuid.UUID `gorm:"type:uuid;index" json:"to_user_id"`
	Amount      float64    `gorm:"type:decimal(15,2);not null" json:"amount"`
	Type        string     `gorm:"not null;index" json:"type"`
	Status      string     `gorm:"default:'pending';index" json:"status"`
	Description string     `json:"description"`
	Metadata    *string    `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt   time.Time  `gorm:"index" json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
