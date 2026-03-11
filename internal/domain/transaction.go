package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransactionStatus string

const (
	TransactionPending   TransactionStatus = "pending"
	TransactionConfirmed TransactionStatus = "confirmed"
	TransactionCancelled TransactionStatus = "cancelled"
	TransactionFailed    TransactionStatus = "failed"
)

type TransactionType string

const (
	TransactionCredit   TransactionType = "credit"
	TransactionDebit    TransactionType = "debit"
	TransactionTransfer TransactionType = "transfer"
	TransactionPayment  TransactionType = "payment"
	TransactionRefund   TransactionType = "refund"
)

type Transaction struct {
	ID          uuid.UUID         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FromUserID  *uuid.UUID        `json:"from_user_id" gorm:"type:uuid;index"`
	ToUserID    *uuid.UUID        `json:"to_user_id" gorm:"type:uuid;index"`
	Amount      float64           `json:"amount" gorm:"type:decimal(15,2);not null"`
	Type        TransactionType   `json:"type" gorm:"not null;index"`
	Status      TransactionStatus `json:"status" gorm:"default:'pending';index"`
	Description string            `json:"description"`
	Metadata    string            `json:"metadata" gorm:"type:jsonb"`
	CreatedAt   time.Time         `json:"created_at" gorm:"index"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
