package services

import (
	"base-gin/internal/domain/models"
	"base-gin/internal/domain/repository"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleService struct {
	*BaseService
	roleRepo       repository.RoleRepository
	userRepo       repository.UserRepository
	permissionRepo repository.PermissionRepository
}

func (s *RoleService) GetRoleByName(param1 string) (any, error) {
	panic("unimplemented")
}

func NewRoleService(db *gorm.DB, roleRepo repository.RoleRepository, userRepo repository.UserRepository, permissionRepo repository.PermissionRepository) *RoleService {
	return &RoleService{
		BaseService:    NewBaseService(db),
		roleRepo:       roleRepo,
		userRepo:       userRepo,
		permissionRepo: permissionRepo,
	}
}

// AssignDefaultRole assigns default "user" role to new user
func (s *RoleService) AssignDefaultRole(userID uuid.UUID) error {
	ctx := context.Background()

	// Get default "user" role
	role, err := s.roleRepo.FindByName(ctx, "user")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("default user role not found")
		}
		return err
	}

	// Assign role to user
	if err := s.roleRepo.AssignRoleToUser(ctx, userID, role.ID); err != nil {
		return err
	}

	return nil
}

// GetUserRoles gets all roles for a user
func (s *RoleService) GetUserRoles(userID uuid.UUID) ([]models.Role, error) {
	ctx := context.Background()
	return s.roleRepo.FindByUserID(ctx, userID)
}

// AssignRole assigns a role to a user
func (s *RoleService) AssignRole(userID, roleID uuid.UUID) error {
	ctx := context.Background()
	return s.roleRepo.AssignRoleToUser(ctx, userID, roleID)
}

// RemoveRole removes a role from a user
func (s *RoleService) RemoveRole(userID, roleID uuid.UUID) error {
	ctx := context.Background()
	return s.roleRepo.RemoveRoleFromUser(ctx, userID, roleID)
}

func (s *RoleService) GetRoleNamesByUserID(ctx context.Context, userID uuid.UUID) ([]string, error) {
	roles, err := s.roleRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	roleNames := make([]string, 0, len(roles))
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}
	return roleNames, nil

}
