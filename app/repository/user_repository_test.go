package repository

import (
	"testing"

	"gorm.io/gorm"
)

func TestNewUserRepository(t *testing.T) {
	var db *gorm.DB
	repo := NewUserRepository(db)

	if repo == nil {
		t.Error("Expected NewUserRepository to return non-nil repository")
	}
}