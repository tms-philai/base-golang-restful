package services

import (
	"testing"

	"gorm.io/gorm"
)

func TestNewBaseService(t *testing.T) {
	var db *gorm.DB
	service := NewBaseService(db)

	if service == nil {
		t.Error("Expected NewBaseService to return non-nil service")
	}
}

func TestBaseService_GetDB(t *testing.T) {
	var db *gorm.DB
	service := NewBaseService(db)

	if service.GetDB() != db {
		t.Error("Expected GetDB to return the same db instance")
	}
}

func TestBaseService_GetTransactionManager(t *testing.T) {
	var db *gorm.DB
	service := NewBaseService(db)

	tm := service.GetTransactionManager()
	if tm == nil {
		t.Error("Expected GetTransactionManager to return non-nil transaction manager")
	}
}

func TestBaseService_WithTransaction(t *testing.T) {
	t.Skip("Skipping test that requires database connection")
}

func TestBaseService_BeginTransaction(t *testing.T) {
	t.Skip("Skipping test that requires database connection")
}

func TestBaseService_CommitTransaction(t *testing.T) {
	t.Skip("Skipping test that requires database connection")
}

func TestBaseService_RollbackTransaction(t *testing.T) {
	t.Skip("Skipping test that requires database connection")
}
