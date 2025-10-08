package database

import (
	"context"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	config := Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		DBName:   "testdb",
		SSLMode:  "disable",
		TimeZone: "UTC",
	}

	if config.Host != "localhost" {
		t.Errorf("Expected Host to be localhost, got %s", config.Host)
	}
	if config.Port != "5432" {
		t.Errorf("Expected Port to be 5432, got %s", config.Port)
	}
}

func TestGetDB(t *testing.T) {
	DB = nil
	db := GetDB()
	if db != nil {
		t.Error("Expected GetDB to return nil when DB is not initialized")
	}
}

func TestHealthCheck_NilDB(t *testing.T) {
	DB = nil
	ctx := context.Background()
	err := HealthCheck(ctx)
	if err == nil {
		t.Error("Expected HealthCheck to return error when DB is nil")
	}
}

func TestClose_NilDB(t *testing.T) {
	DB = nil
	err := Close()
	if err != nil {
		t.Errorf("Expected Close to return nil when DB is nil, got %v", err)
	}
}

func TestConnectWithRetry_InvalidConfig(t *testing.T) {
	config := Config{
		Host:     "invalid-host",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		DBName:   "testdb",
		SSLMode:  "disable",
		TimeZone: "UTC",
	}

	err := ConnectWithRetry(config, 2, 100*time.Millisecond)
	if err == nil {
		t.Error("Expected ConnectWithRetry to return error with invalid config")
	}
}
