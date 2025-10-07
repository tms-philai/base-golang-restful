package mocks

import (
	"base-golang-restful-app/models"
	"base-golang-restful-app/tests/interfaces"

	"github.com/stretchr/testify/mock"
)

// MockJWTService is a mock implementation of JWTService
type MockJWTService struct {
	mock.Mock
}

// Ensure MockJWTService implements JWTServiceInterface
var _ interfaces.JWTServiceInterface = (*MockJWTService)(nil)

// GenerateTokenPair mocks the GenerateTokenPair method
func (m *MockJWTService) GenerateTokenPair(user *models.User) (*models.TokenPair, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TokenPair), args.Error(1)
}

// ValidateToken mocks the ValidateToken method
func (m *MockJWTService) ValidateToken(tokenString string) (*models.JWTClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JWTClaims), args.Error(1)
}

// RefreshToken mocks the RefreshToken method
func (m *MockJWTService) RefreshToken(refreshTokenString string, user *models.User) (*models.TokenPair, error) {
	args := m.Called(refreshTokenString, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TokenPair), args.Error(1)
}

// ExtractTokenFromHeader mocks the ExtractTokenFromHeader method
func (m *MockJWTService) ExtractTokenFromHeader(authHeader string) (string, error) {
	args := m.Called(authHeader)
	return args.String(0), args.Error(1)
}
