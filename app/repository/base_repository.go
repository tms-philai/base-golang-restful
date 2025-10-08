package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseRepository interface {
	Create(ctx context.Context, entity interface{}) error
	FindByID(ctx context.Context, id uuid.UUID, entity interface{}) error
	FindAll(ctx context.Context, entities interface{}) error
	Update(ctx context.Context, entity interface{}) error
	Delete(ctx context.Context, id uuid.UUID, entity interface{}) error
	Count(ctx context.Context, entity interface{}) (int64, error)
	Exists(ctx context.Context, id uuid.UUID, entity interface{}) (bool, error)
}

type baseRepository struct {
	db *gorm.DB
}

func NewBaseRepository(db *gorm.DB) BaseRepository {
	return &baseRepository{db: db}
}

func (r *baseRepository) Create(ctx context.Context, entity interface{}) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *baseRepository) FindByID(ctx context.Context, id uuid.UUID, entity interface{}) error {
	return r.db.WithContext(ctx).First(entity, "id = ?", id).Error
}

func (r *baseRepository) FindAll(ctx context.Context, entities interface{}) error {
	return r.db.WithContext(ctx).Find(entities).Error
}

func (r *baseRepository) Update(ctx context.Context, entity interface{}) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *baseRepository) Delete(ctx context.Context, id uuid.UUID, entity interface{}) error {
	return r.db.WithContext(ctx).Delete(entity, "id = ?", id).Error
}

func (r *baseRepository) Count(ctx context.Context, entity interface{}) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(entity).Count(&count).Error
	return count, err
}

func (r *baseRepository) Exists(ctx context.Context, id uuid.UUID, entity interface{}) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(entity).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
