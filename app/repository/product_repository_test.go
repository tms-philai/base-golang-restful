package repository

import (
	"testing"

	"gorm.io/gorm"
)

func TestNewProductRepository(t *testing.T) {
	var db *gorm.DB
	repo := NewProductRepository(db)

	if repo == nil {
		t.Error("Expected NewProductRepository to return non-nil repository")
	}
}
