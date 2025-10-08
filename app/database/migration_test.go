package database

import (
	"testing"

	"gorm.io/gorm"
)

type TestModel struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"size:100"`
}

func TestNewMigrator(t *testing.T) {
	var db *gorm.DB
	migrator := NewMigrator(db)

	if migrator == nil {
		t.Error("Expected NewMigrator to return non-nil migrator")
	}

	if migrator.db != db {
		t.Error("Expected migrator.db to equal input db")
	}
}