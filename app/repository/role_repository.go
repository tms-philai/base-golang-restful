package repository

import (
	"context"

	"base-golang-restful-app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindByName(ctx context.Context, name string) (*models.Role, error)
	FindByNameWithPermissions(ctx context.Context, name string) (*models.Role, error)
	FindWithPermissions(ctx context.Context, id uuid.UUID) (*models.Role, error)
	AssignPermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
	RemovePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
	GetRoleUsers(ctx context.Context, roleID uuid.UUID) ([]models.User, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindByNameWithPermissions(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		Where("name = ?", name).
		First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindWithPermissions(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		Where("id = ?", id).
		First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) AssignPermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	var role models.Role
	if err := r.db.WithContext(ctx).First(&role, roleID).Error; err != nil {
		return err
	}

	var permissions []models.Permission
	if err := r.db.WithContext(ctx).Find(&permissions, permissionIDs).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(&role).Association("Permissions").Append(permissions)
}

func (r *roleRepository) RemovePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	var role models.Role
	if err := r.db.WithContext(ctx).First(&role, roleID).Error; err != nil {
		return err
	}

	var permissions []models.Permission
	if err := r.db.WithContext(ctx).Find(&permissions, permissionIDs).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(&role).Association("Permissions").Delete(permissions)
}

func (r *roleRepository) GetRoleUsers(ctx context.Context, roleID uuid.UUID) ([]models.User, error) {
	var role models.Role
	if err := r.db.WithContext(ctx).Preload("Users").First(&role, roleID).Error; err != nil {
		return nil, err
	}
	return role.Users, nil
}
