package config

import (
	"testing"
	"time"
)

func TestLoadJWTExpiry(t *testing.T) {
	for _, tc := range []struct {
		name, access, refresh string
		wantErr               bool
	}{
		{"valid", "15m", "168h", false},
		{"unsupported day unit", "15m", "7d", true},
		{"invalid access duration", "invalid", "168h", true},
		{"zero access duration", "0s", "168h", true},
		{"negative refresh duration", "15m", "-1h", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("JWT_ACCESS_EXPIRY", tc.access)
			t.Setenv("JWT_REFRESH_EXPIRY", tc.refresh)
			cfg, err := Load()
			if (err != nil) != tc.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if err == nil && (cfg.JWT.AccessExpiry != 15*time.Minute || cfg.JWT.RefreshExpiry != 168*time.Hour) {
				t.Fatal("unexpected token expiry durations")
			}
		})
	}
}
