package database

import (
	"testing"

	"gorm.io/gorm"
)

func TestNewSeeder(t *testing.T) {
	var db *gorm.DB
	seeder := NewSeeder(db)

	if seeder == nil {
		t.Error("Expected NewSeeder to return non-nil seeder")
	}

	if seeder.db != db {
		t.Error("Expected seeder.db to equal input db")
	}
}