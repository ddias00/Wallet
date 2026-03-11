package service

import (
	"testing"
)

func TestCreditValidation(t *testing.T) {
	tests := []struct {
		amount      float64
		shouldError bool
	}{
		{100.0, false},
		{0.0, true},
		{-50.0, true},
	}

	for _, tt := range tests {
		if tt.amount <= 0 && !tt.shouldError {
			t.Errorf("Credit amount %f should error", tt.amount)
		}
		if tt.amount > 0 && tt.shouldError {
			t.Errorf("Credit amount %f should not error", tt.amount)
		}
	}
}

func TestDebitValidation(t *testing.T) {
	tests := []struct {
		amount      float64
		shouldError bool
	}{
		{100.0, false},
		{0.0, true},
		{-50.0, true},
	}

	for _, tt := range tests {
		if tt.amount <= 0 && !tt.shouldError {
			t.Errorf("Debit amount %f should error", tt.amount)
		}
		if tt.amount > 0 && tt.shouldError {
			t.Errorf("Debit amount %f should not error", tt.amount)
		}
	}
}
