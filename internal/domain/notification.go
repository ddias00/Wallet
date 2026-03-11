package domain

import (
	"time"

	"github.com/google/uuid"
)

type NotificationEventType string

const (
	NotificationTransactionCreated   NotificationEventType = "transaction_created"
	NotificationTransactionConfirmed NotificationEventType = "transaction_confirmed"
	NotificationTransactionCancelled NotificationEventType = "transaction_cancelled"
	NotificationPaymentReceived      NotificationEventType = "payment_received"
	NotificationBoletoGenerated      NotificationEventType = "boleto_generated"
	NotificationPixKeyRegistered     NotificationEventType = "pix_key_registered"
	NotificationAccountBlocked       NotificationEventType = "account_blocked"
	NotificationAccountUnblocked     NotificationEventType = "account_unblocked"
)

type NotificationDeliveryStatus string

const (
	NotificationPending   NotificationDeliveryStatus = "pending"
	NotificationSent      NotificationDeliveryStatus = "sent"
	NotificationDelivered NotificationDeliveryStatus = "delivered"
	NotificationFailed    NotificationDeliveryStatus = "failed"
)

type NotificationChannel string

const (
	NotificationEmail NotificationChannel = "email"
	NotificationSMS   NotificationChannel = "sms"
	NotificationPush  NotificationChannel = "push"
	NotificationInApp NotificationChannel = "in_app"
)

type Notification struct {
	ID             uuid.UUID                  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID         uuid.UUID                  `json:"user_id" gorm:"type:uuid;not null;index"`
	EventType      NotificationEventType      `json:"event_type" gorm:"not null;index"`
	Channel        NotificationChannel        `json:"channel" gorm:"not null;index"`
	Message        string                     `json:"message" gorm:"not null"`
	DeliveryStatus NotificationDeliveryStatus `json:"delivery_status" gorm:"default:'pending';index"`
	SentAt         *time.Time                 `json:"sent_at"`
	DeliveredAt    *time.Time                 `json:"delivered_at"`
	FailedAt       *time.Time                 `json:"failed_at"`
	ErrorMessage   string                     `json:"error_message"`
	Metadata       string                     `json:"metadata" gorm:"type:jsonb"`
	CreatedAt      time.Time                  `json:"created_at" gorm:"index"`
	UpdatedAt      time.Time                  `json:"updated_at"`
}
