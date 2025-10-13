package services

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"base-golang-restful-app/models"
	"base-golang-restful-app/utils"

	"github.com/google/uuid"
)

type UserService struct {
	users map[string]*models.User
	mutex sync.RWMutex
}

func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]*models.User),
		mutex: sync.RWMutex{},
	}
}

func (s *UserService) Create(req models.UserCreateRequest) (*models.User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if errors := utils.ValidatePasswordStrength(req.Password); len(errors) > 0 {
		return nil, fmt.Errorf("password validation failed: %s", strings.Join(errors, ", "))
	}

	for _, user := range s.users {
		if user.Email == req.Email {
			return nil, errors.New("email already exists")
		}
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	s.users[user.ID.String()] = user
	return user, nil
}

func (s *UserService) GetByID(id string) (*models.User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

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

func (s *UserService) List(page, pageSize int) ([]*models.User, int, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	users := make([]*models.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}

	total := len(users)
	start := (page - 1) * pageSize
	if start > total {
		return []*models.User{}, total, nil
	}

	end := start + pageSize
	if end > total {
		end = total
	}

	return users[start:end], total, nil
}

func (s *UserService) Update(id string, req models.UserUpdateRequest) (*models.User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	if req.Email != nil {
		for _, existingUser := range s.users {
			if existingUser.ID.String() != id && existingUser.Email == *req.Email {
				return nil, errors.New("email already exists")
			}
		}
		user.Email = *req.Email
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		user.LastName = *req.LastName
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	return user, nil
}

func (s *UserService) Delete(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(s.users, id)
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
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, exists := s.users[userID]
	if !exists {
		return errors.New("user not found")
	}

	if err := utils.CheckPassword(oldPassword, user.Password); err != nil {
		return errors.New("invalid old password")
	}

	if errors := utils.ValidatePasswordStrength(newPassword); len(errors) > 0 {
		return fmt.Errorf("password validation failed: %s", strings.Join(errors, ", "))
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword
	return nil
}
