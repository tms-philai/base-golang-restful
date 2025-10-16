package services

import (
	"base-gin/internal/pkg/database"
	"context"

	"gorm.io/gorm"
)

type BaseService struct {
	db *gorm.DB
	tm database.TransactionManager
}

func NewBaseService(db *gorm.DB) *BaseService {
	return &BaseService{
		db: db,
		tm: database.NewTransactionManager(db),
	}
}

func (s *BaseService) GetDB() *gorm.DB {
	return s.db
}

func (s *BaseService) GetTransactionManager() database.TransactionManager {
	return s.tm
}

func (s *BaseService) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return s.tm.WithTransaction(ctx, fn)
}

func (s *BaseService) BeginTransaction(ctx context.Context) (*gorm.DB, error) {
	return s.tm.BeginTransaction(ctx)
}

func (s *BaseService) CommitTransaction(tx *gorm.DB) error {
	return s.tm.CommitTransaction(tx)
}

func (s *BaseService) RollbackTransaction(tx *gorm.DB) error {
	return s.tm.RollbackTransaction(tx)
}
