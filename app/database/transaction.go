package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	WithTransactionTimeout(ctx context.Context, timeout time.Duration, fn func(tx *gorm.DB) error) error
	BeginTransaction(ctx context.Context) (*gorm.DB, error)
	CommitTransaction(tx *gorm.DB) error
	RollbackTransaction(tx *gorm.DB) error
}

type transactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &transactionManager{db: db}
}

func (tm *transactionManager) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := tm.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback().Error; rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %w", err, rbErr)
		}
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (tm *transactionManager) WithTransactionTimeout(ctx context.Context, timeout time.Duration, fn func(tx *gorm.DB) error) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		errChan <- tm.WithTransaction(ctx, fn)
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return fmt.Errorf("transaction timeout: %w", ctx.Err())
	}
}

func (tm *transactionManager) BeginTransaction(ctx context.Context) (*gorm.DB, error) {
	tx := tm.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	return tx, nil
}

func (tm *transactionManager) CommitTransaction(tx *gorm.DB) error {
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (tm *transactionManager) RollbackTransaction(tx *gorm.DB) error {
	if err := tx.Rollback().Error; err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	return nil
}

type SavepointManager struct {
	tx *gorm.DB
}

func NewSavepointManager(tx *gorm.DB) *SavepointManager {
	return &SavepointManager{tx: tx}
}

func (sm *SavepointManager) CreateSavepoint(name string) error {
	return sm.tx.Exec(fmt.Sprintf("SAVEPOINT %s", name)).Error
}

func (sm *SavepointManager) RollbackToSavepoint(name string) error {
	return sm.tx.Exec(fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", name)).Error
}

func (sm *SavepointManager) ReleaseSavepoint(name string) error {
	return sm.tx.Exec(fmt.Sprintf("RELEASE SAVEPOINT %s", name)).Error
}

func WithNestedTransaction(tx *gorm.DB, savepointName string, fn func(tx *gorm.DB) error) error {
	sm := NewSavepointManager(tx)

	if err := sm.CreateSavepoint(savepointName); err != nil {
		return fmt.Errorf("failed to create savepoint: %w", err)
	}

	if err := fn(tx); err != nil {
		if rbErr := sm.RollbackToSavepoint(savepointName); rbErr != nil {
			return fmt.Errorf("nested transaction error: %v, rollback error: %w", err, rbErr)
		}
		return err
	}

	if err := sm.ReleaseSavepoint(savepointName); err != nil {
		return fmt.Errorf("failed to release savepoint: %w", err)
	}

	return nil
}

type IsolationLevel string

const (
	ReadUncommitted IsolationLevel = "READ UNCOMMITTED"
	ReadCommitted   IsolationLevel = "READ COMMITTED"
	RepeatableRead  IsolationLevel = "REPEATABLE READ"
	Serializable    IsolationLevel = "SERIALIZABLE"
)

func (tm *transactionManager) WithTransactionIsolation(ctx context.Context, level IsolationLevel, fn func(tx *gorm.DB) error) error {
	tx := tm.db.WithContext(ctx)

	if err := tx.Exec(fmt.Sprintf("SET TRANSACTION ISOLATION LEVEL %s", level)).Error; err != nil {
		return fmt.Errorf("failed to set isolation level: %w", err)
	}

	return tm.WithTransaction(ctx, fn)
}
