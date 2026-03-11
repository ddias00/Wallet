package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Boleto struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	TransactionID *uuid.UUID `gorm:"type:uuid;index" json:"transaction_id,omitempty"`
	Amount        float64    `gorm:"type:decimal(15,2);not null" json:"amount"`
	Code          string     `gorm:"unique;not null;index" json:"code"`
	Barcode       string     `gorm:"unique;not null" json:"barcode"`
	Status        string     `gorm:"default:'pending';index" json:"status"`
	DueDate       time.Time  `gorm:"not null;index" json:"due_date"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
	CreatedAt     time.Time  `gorm:"index" json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (b *Boleto) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
