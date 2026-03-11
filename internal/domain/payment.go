package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethod string

const (
	PaymentBoleto PaymentMethod = "boleto"
	PaymentPIX    PaymentMethod = "pix"
	PaymentQRCode PaymentMethod = "qrcode"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentCompleted PaymentStatus = "completed"
	PaymentExpired   PaymentStatus = "expired"
	PaymentCancelled PaymentStatus = "cancelled"
)

type Payment struct {
	ID            uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        uuid.UUID     `json:"user_id" gorm:"type:uuid;not null;index"`
	TransactionID *uuid.UUID    `json:"transaction_id" gorm:"type:uuid;index"`
	Amount        float64       `json:"amount" gorm:"type:decimal(15,2);not null"`
	Method        PaymentMethod `json:"method" gorm:"not null;index"`
	Status        PaymentStatus `json:"status" gorm:"default:'pending';index"`
	Code          string        `json:"code" gorm:"unique;index"`
	QRCode        string        `json:"qr_code"`
	ExpiresAt     *time.Time    `json:"expires_at"`
	CreatedAt     time.Time     `json:"created_at" gorm:"index"`
	UpdatedAt     time.Time     `json:"updated_at"`
}
