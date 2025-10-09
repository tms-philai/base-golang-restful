package repository

import (
	"testing"

	"gorm.io/gorm"
)

func TestNewBaseRepository(t *testing.T) {
	var db *gorm.DB
	repo := NewBaseRepository(db)

	if repo == nil {
		t.Error("Expected NewBaseRepository to return non-nil repository")
	}
}
