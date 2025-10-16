package repository

import (
	"base-gin/internal/domain/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	FindByName(ctx context.Context, name string) (*models.Permission, error)
	FindByResource(ctx context.Context, resource string) ([]models.Permission, error)
	FindByResourceAndAction(ctx context.Context, resource, action string) (*models.Permission, error)
	FindByNames(ctx context.Context, names []string) ([]models.Permission, error)
	GetPermissionRoles(ctx context.Context, permissionID uuid.UUID) ([]models.Role, error)
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{
		db: db,
	}
}

func (r *permissionRepository) FindByName(ctx context.Context, name string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) FindByResource(ctx context.Context, resource string) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.WithContext(ctx).Where("resource = ?", resource).Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *permissionRepository) FindByResourceAndAction(ctx context.Context, resource, action string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.WithContext(ctx).
		Where("resource = ? AND action = ?", resource, action).
		First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) FindByNames(ctx context.Context, names []string) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.WithContext(ctx).Where("name IN ?", names).Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *permissionRepository) GetPermissionRoles(ctx context.Context, permissionID uuid.UUID) ([]models.Role, error) {
	var permission models.Permission
	if err := r.db.WithContext(ctx).Preload("Roles").First(&permission, permissionID).Error; err != nil {
		return nil, err
	}
	return permission.Roles, nil
}
