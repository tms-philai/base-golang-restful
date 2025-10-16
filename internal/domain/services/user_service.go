package services

import (
	"base-gin/internal/domain/models"
	"base-gin/internal/domain/repository"
	"base-gin/internal/pkg/utils"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Create(req models.UserCreateRequest) (*models.User, error) {
	ctx := context.Background()

	// Validate password strength
	if errs := utils.ValidatePasswordStrength(req.Password); len(errs) > 0 {
		return nil, fmt.Errorf("password validation failed: %s", strings.Join(errs, ", "))
	}

	// Check if email already exists
	_, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	// Save to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetByID(id string) (*models.User, error) {
	ctx := context.Background()

	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	var user models.User
	if err := s.userRepo.FindByID(ctx, userID, &user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByIDWithRoles gets user with roles and permissions
func (s *UserService) GetByIDWithRoles(id string) (*models.User, error) {
	ctx := context.Background()

	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.FindByIDWithRoles(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {
	ctx := context.Background()

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *UserService) List(page, pageSize int) ([]*models.User, models.PaginationMetadata, error) {
	ctx := context.Background()

	users, total, err := s.userRepo.FindAllWithPagination(ctx, page, pageSize)
	if err != nil {
		return nil, models.PaginationMetadata{}, fmt.Errorf("failed to list users: %w", err)
	}

	// Convert []models.User to []*models.User
	userPtrs := make([]*models.User, len(users))
	for i := range users {
		userPtrs[i] = &users[i]
	}

	pagination := models.NewPaginationMetadata(page, pageSize, total)

	return userPtrs, pagination, nil
}

func (s *UserService) Update(id string, req models.UserUpdateRequest) (*models.User, error) {
	ctx := context.Background()

	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Get existing user
	var user models.User
	if err := s.userRepo.FindByID(ctx, userID, &user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check email uniqueness if email is being updated
	if req.Email != nil && *req.Email != user.Email {
		_, err := s.userRepo.FindByEmail(ctx, *req.Email)
		if err == nil {
			return nil, errors.New("email already exists")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to check email: %w", err)
		}
		user.Email = *req.Email
	}

	// Update fields
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		user.LastName = *req.LastName
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Save updates
	if err := s.userRepo.Update(ctx, &user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}

func (s *UserService) Delete(id string) error {
	ctx := context.Background()

	userID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Check if user exists
	var user models.User
	if err := s.userRepo.FindByID(ctx, userID, &user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Delete user
	if err := s.userRepo.Delete(ctx, userID, &user); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *UserService) ValidateCredentials(email, password string) (*models.User, error) {
	user, err := s.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := utils.CheckPassword(password, user.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	return user, nil
}

func (s *UserService) ChangePassword(userID, oldPassword, newPassword string) error {
	ctx := context.Background()

	uid, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Get user
	var user models.User
	if err := s.userRepo.FindByID(ctx, uid, &user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify old password
	if err := utils.CheckPassword(oldPassword, user.Password); err != nil {
		return errors.New("invalid old password")
	}

	// Validate new password
	if errs := utils.ValidatePasswordStrength(newPassword); len(errs) > 0 {
		return fmt.Errorf("password validation failed: %s", strings.Join(errs, ", "))
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.Password = hashedPassword
	if err := s.userRepo.Update(ctx, &user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
