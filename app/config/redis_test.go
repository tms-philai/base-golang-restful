package config

import (
	"testing"
)

func TestRedisConfig_GetAddress(t *testing.T) {
	config := &RedisConfig{
		Host: "localhost",
		Port: "6379",
	}

	expected := "localhost:6379"
	if config.GetAddress() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, config.GetAddress())
	}
}

func TestLoadRedisConfig(t *testing.T) {
	config := loadRedisConfig()

	if config.Host == "" {
		t.Error("Expected Redis host to be set")
	}

	if config.Port == "" {
		t.Error("Expected Redis port to be set")
	}
}
