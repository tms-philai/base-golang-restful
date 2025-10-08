package config

import (
	"testing"
	"time"
)

func TestJWTConfig_GetExpiration(t *testing.T) {
	config := &JWTConfig{
		ExpirationHours: 24,
	}

	expected := 24 * time.Hour
	if config.GetExpiration() != expected {
		t.Errorf("Expected %v, got %v", expected, config.GetExpiration())
	}
}

func TestJWTConfig_GetRefreshExpiration(t *testing.T) {
	config := &JWTConfig{
		RefreshExpHours: 168,
	}

	expected := 168 * time.Hour
	if config.GetRefreshExpiration() != expected {
		t.Errorf("Expected %v, got %v", expected, config.GetRefreshExpiration())
	}
}

func TestLoadJWTConfig(t *testing.T) {
	config := loadJWTConfig()

	if config.Secret == "" {
		t.Error("Expected JWT secret to be set")
	}

	if config.ExpirationHours <= 0 {
		t.Error("Expected expiration hours to be positive")
	}
}
