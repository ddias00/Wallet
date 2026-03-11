package service

import (
	"testing"

	"github.com/wallet/pkg/validator"
)

func TestValidateCPF(t *testing.T) {
	tests := []struct {
		cpf   string
		valid bool
	}{
		{"12345678909", true},
		{"00000000000", false},
		{"123", false},
	}

	for _, tt := range tests {
		result := validator.ValidateCPF(tt.cpf)
		if result != tt.valid {
			t.Errorf("ValidateCPF(%s) = %v, want %v", tt.cpf, result, tt.valid)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"invalid-email", false},
		{"test@", false},
	}

	for _, tt := range tests {
		result := validator.ValidateEmail(tt.email)
		if result != tt.valid {
			t.Errorf("ValidateEmail(%s) = %v, want %v", tt.email, result, tt.valid)
		}
	}
}

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		phone string
		valid bool
	}{
		{"11987654321", true},
		{"1198765432", true},
		{"123", false},
	}

	for _, tt := range tests {
		result := validator.ValidatePhone(tt.phone)
		if result != tt.valid {
			t.Errorf("ValidatePhone(%s) = %v, want %v", tt.phone, result, tt.valid)
		}
	}
}
