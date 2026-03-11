package domain

import (
	"time"

	"github.com/google/uuid"
)

type PixKeyType string

const (
	PixKeyCPF   PixKeyType = "cpf"
	PixKeyEmail PixKeyType = "email"
	PixKeyPhone PixKeyType = "phone"
	PixKeyRandom PixKeyType = "random"
)

type PixKeyStatus string

const (
	PixKeyActive   PixKeyStatus = "active"
	PixKeyInactive PixKeyStatus = "inactive"
	PixKeyDeleted  PixKeyStatus = "deleted"
)

type PixKey struct {
	ID        uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID    `json:"user_id" gorm:"type:uuid;not null;index"`
	Key       string       `json:"key" gorm:"unique;not null;index"`
	Type      PixKeyType   `json:"type" gorm:"not null;index"`
	Status    PixKeyStatus `json:"status" gorm:"default:'active';index"`
	CreatedAt time.Time    `json:"created_at" gorm:"index"`
	UpdatedAt time.Time    `json:"updated_at"`
}
