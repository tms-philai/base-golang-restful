package mocks

import (
	"base-golang-restful-app/models"
	"base-golang-restful-app/tests/interfaces"

	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

// Ensure MockUserService implements UserServiceInterface
var _ interfaces.UserServiceInterface = (*MockUserService)(nil)

// Create mocks the Create method
func (m *MockUserService) Create(req models.UserCreateRequest) (*models.User, error) {
	args := m.Called(req)
	return args.Get(0).(*models.User), args.Error(1)
}

// GetByID mocks the GetByID method
func (m *MockUserService) GetByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// GetByUsername mocks the GetByUsername method
func (m *MockUserService) GetByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// GetByEmail mocks the GetByEmail method
func (m *MockUserService) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// Update mocks the Update method
func (m *MockUserService) Update(id string, req models.UserUpdateRequest) (*models.User, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// Delete mocks the Delete method
func (m *MockUserService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// List mocks the List method
func (m *MockUserService) List(page, limit int) ([]models.User, models.PaginationMetadata, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]models.User), args.Get(1).(models.PaginationMetadata), args.Error(2)
}

// ChangePassword mocks the ChangePassword method
func (m *MockUserService) ChangePassword(userID, currentPassword, newPassword string) error {
	args := m.Called(userID, currentPassword, newPassword)
	return args.Error(0)
}

// ValidateCredentials mocks the ValidateCredentials method
func (m *MockUserService) ValidateCredentials(username, password string) (*models.User, error) {
	args := m.Called(username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
