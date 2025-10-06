package services

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"base-golang-restful-app/models"
	"base-golang-restful-app/utils"
)

// UserService handles user-related business logic
type UserService struct {
	users map[string]*models.User // In-memory storage for demo
	mutex sync.RWMutex
}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]*models.User),
		mutex: sync.RWMutex{},
	}
}

// Create creates a new user
func (s *UserService) Create(req models.UserCreateRequest) (*models.User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Validate password strength
	if errors := utils.ValidatePasswordStrength(req.Password); len(errors) > 0 {
		return nil, fmt.Errorf("password validation failed: %s", strings.Join(errors, ", "))
	}

	// Check if username already exists
	for _, user := range s.users {
		if user.Username == req.Username {
			return nil, errors.New("username already exists")
		}
		if user.Email == req.Email {
			return nil, errors.New("email already exists")
		}
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user
	user := models.NewUser(req)
	user.Password = hashedPassword

	// Store user
	s.users[user.ID] = user

	return user, nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id string) (*models.User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetByUsername retrieves a user by username
func (s *UserService) GetByUsername(username string) (*models.User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, user := range s.users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

// GetByEmail retrieves a user by email
func (s *UserService) GetByEmail(email string) (*models.User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, user := range s.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

// Update updates a user
func (s *UserService) Update(id string, req models.UserUpdateRequest) (*models.User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	// Check for username/email conflicts
	if req.Username != nil {
		for _, existingUser := range s.users {
			if existingUser.ID != id && existingUser.Username == *req.Username {
				return nil, errors.New("username already exists")
			}
		}
	}

	if req.Email != nil {
		for _, existingUser := range s.users {
			if existingUser.ID != id && existingUser.Email == *req.Email {
				return nil, errors.New("email already exists")
			}
		}
	}

	// Apply updates
	user.Update(req)

	return user, nil
}

// Delete deletes a user (soft delete by setting IsActive to false)
func (s *UserService) Delete(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, exists := s.users[id]
	if !exists {
		return errors.New("user not found")
	}

	user.IsActive = false
	return nil
}

// List retrieves all users with pagination
func (s *UserService) List(page, limit int) ([]models.User, models.PaginationMetadata, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Convert map to slice
	var allUsers []models.User
	for _, user := range s.users {
		if user.IsActive { // Only return active users
			allUsers = append(allUsers, *user)
		}
	}

	total := int64(len(allUsers))
	
	// Calculate pagination
	offset := (page - 1) * limit
	end := offset + limit

	if offset > len(allUsers) {
		return []models.User{}, models.NewPaginationMetadata(page, limit, total), nil
	}

	if end > len(allUsers) {
		end = len(allUsers)
	}

	users := allUsers[offset:end]
	pagination := models.NewPaginationMetadata(page, limit, total)

	return users, pagination, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(userID, currentPassword, newPassword string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, exists := s.users[userID]
	if !exists {
		return errors.New("user not found")
	}

	// Verify current password
	if err := utils.CheckPassword(currentPassword, user.Password); err != nil {
		return errors.New("current password is incorrect")
	}

	// Validate new password strength
	if errors := utils.ValidatePasswordStrength(newPassword); len(errors) > 0 {
		return fmt.Errorf("password validation failed: %s", strings.Join(errors, ", "))
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword
	return nil
}

// ValidateCredentials validates user credentials for login
func (s *UserService) ValidateCredentials(username, password string) (*models.User, error) {
	user, err := s.GetByUsername(username)
	if err != nil {
		// Try to find by email if username lookup fails
		user, err = s.GetByEmail(username)
		if err != nil {
			return nil, errors.New("invalid credentials")
		}
	}

	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	if err := utils.CheckPassword(password, user.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
