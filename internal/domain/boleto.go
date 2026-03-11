package domain

import (
	"time"

	"github.com/google/uuid"
)

type BoletoStatus string

const (
	BoletoPending   BoletoStatus = "pending"
	BoletoPaid      BoletoStatus = "paid"
	BoletoExpired   BoletoStatus = "expired"
	BoletoCancelled BoletoStatus = "cancelled"
)

type Boleto struct {
	ID            uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        uuid.UUID    `json:"user_id" gorm:"type:uuid;not null;index"`
	TransactionID *uuid.UUID   `json:"transaction_id" gorm:"type:uuid;index"`
	Amount        float64      `json:"amount" gorm:"type:decimal(15,2);not null"`
	Code          string       `json:"code" gorm:"unique;not null;index"`
	Barcode       string       `json:"barcode" gorm:"unique;not null"`
	Status        BoletoStatus `json:"status" gorm:"default:'pending';index"`
	DueDate       time.Time    `json:"due_date" gorm:"not null;index"`
	PaidAt        *time.Time   `json:"paid_at"`
	CreatedAt     time.Time    `json:"created_at" gorm:"index"`
	UpdatedAt     time.Time    `json:"updated_at"`
}
