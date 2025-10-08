package database

import (
	"testing"

	"gorm.io/gorm"
)

func TestNewTransactionManager(t *testing.T) {
	var db *gorm.DB
	tm := NewTransactionManager(db)

	if tm == nil {
		t.Error("Expected NewTransactionManager to return non-nil transaction manager")
	}
}

func TestTransactionManager_Methods(t *testing.T) {
	var db *gorm.DB
	tm := NewTransactionManager(db)

	if tm == nil {
		t.Fatal("Transaction manager should not be nil")
	}
}

func TestWithNestedTransaction(t *testing.T) {
	t.Skip("Skipping test that requires database connection")
}

func TestNewSavepointManager(t *testing.T) {
	var tx *gorm.DB
	sm := NewSavepointManager(tx)

	if sm == nil {
		t.Error("Expected NewSavepointManager to return non-nil savepoint manager")
	}

	if sm.tx != tx {
		t.Error("Expected savepoint manager to have the same tx")
	}
}

func TestSavepointManager_Methods(t *testing.T) {
	t.Skip("Skipping test that requires database connection")
}

func TestTransactionManager_WithTransactionTimeout(t *testing.T) {
	t.Skip("Skipping test that requires database connection")
}

func TestTransactionManager_WithTransaction_Rollback(t *testing.T) {
	t.Skip("Skipping test that requires database connection")
}

func TestIsolationLevels(t *testing.T) {
	levels := []IsolationLevel{
		ReadUncommitted,
		ReadCommitted,
		RepeatableRead,
		Serializable,
	}

	expectedValues := []string{
		"READ UNCOMMITTED",
		"READ COMMITTED",
		"REPEATABLE READ",
		"SERIALIZABLE",
	}

	for i, level := range levels {
		if string(level) != expectedValues[i] {
			t.Errorf("Expected isolation level '%s', got '%s'", expectedValues[i], string(level))
		}
	}
}
