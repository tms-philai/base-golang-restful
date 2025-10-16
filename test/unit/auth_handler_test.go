package unit

import (
	"base-gin/internal/app/handlers"
	"base-gin/internal/domain/models"
	"base-gin/internal/pkg/auth"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of UserService
type MockUserServiceInterface struct {
	mock.Mock
}

func (m *MockUserServiceInterface) Create(req models.UserCreateRequest) (*models.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserServiceInterface) GetByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserServiceInterface) GetByIDWithRoles(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserServiceInterface) ValidateCredentials(email, password string) (*models.User, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserServiceInterface) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserServiceInterface) List(page, pageSize int) ([]*models.User, models.PaginationMetadata, error) {
	args := m.Called(page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(models.PaginationMetadata), args.Error(2)
	}
	return args.Get(0).([]*models.User), args.Get(1).(models.PaginationMetadata), args.Error(2)
}

func (m *MockUserServiceInterface) Update(id string, req models.UserUpdateRequest) (*models.User, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserServiceInterface) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserServiceInterface) ChangePassword(userID, oldPassword, newPassword string) error {
	args := m.Called(userID, oldPassword, newPassword)
	return args.Error(0)
}

// MockRoleService is a mock implementation of RoleService
type MockRoleServiceInterface struct {
	mock.Mock
}

func (m *MockRoleServiceInterface) AssignDefaultRole(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

// MockJWTManager is a mock implementation of JWTManager
type MockJWTManagerInterface struct {
	mock.Mock
}

func (m *MockJWTManagerInterface) GenerateAccessToken(userID uuid.UUID, email string) (string, error) {
	args := m.Called(userID, email)
	return args.String(0), args.Error(1)
}

func (m *MockJWTManagerInterface) GenerateRefreshToken(userID uuid.UUID, email string) (string, error) {
	args := m.Called(userID, email)
	return args.String(0), args.Error(1)
}

func (m *MockJWTManagerInterface) ValidateToken(token string, expectedType auth.TokenType) (*auth.JWTClaims, error) {
	args := m.Called(token, expectedType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.JWTClaims), args.Error(1)
}

func (m *MockJWTManagerInterface) GetTokenDuration(tokenType auth.TokenType) time.Duration {
	args := m.Called(tokenType)
	return args.Get(0).(time.Duration)
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		request        models.RegisterRequest
		mockSetup      func(*MockUserServiceInterface, *MockRoleServiceInterface, *MockJWTManagerInterface)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful registration",
			request: models.RegisterRequest{
				Username:  "testuser",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			mockSetup: func(userService *MockUserServiceInterface, roleService *MockRoleServiceInterface, jwtManager *MockJWTManagerInterface) {
				user := &models.User{
					ID:        uuid.New(),
					Email:     "test@example.com",
					FirstName: "Test",
					LastName:  "User",
					IsActive:  true,
				}
				userWithRoles := &models.User{
					ID:        user.ID,
					Email:     "test@example.com",
					FirstName: "Test",
					LastName:  "User",
					IsActive:  true,
				}

				userService.On("Create", mock.AnythingOfType("models.UserCreateRequest")).Return(user, nil)
				roleService.On("AssignDefaultRole", user.ID).Return(nil)
				userService.On("GetByIDWithRoles", user.ID.String()).Return(userWithRoles, nil)
				jwtManager.On("GenerateAccessToken", user.ID, user.Email).Return("access_token", nil)
				jwtManager.On("GenerateRefreshToken", user.ID, user.Email).Return("refresh_token", nil)
				jwtManager.On("GetTokenDuration", auth.AccessToken).Return(15 * time.Minute)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "validation error - missing email",
			request: models.RegisterRequest{
				Username:  "testuser",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			mockSetup:      func(*MockUserServiceInterface, *MockRoleServiceInterface, *MockJWTManagerInterface) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name: "user creation failed - email already exists",
			request: models.RegisterRequest{
				Username:  "testuser",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			mockSetup: func(userService *MockUserServiceInterface, roleService *MockRoleServiceInterface, jwtManager *MockJWTManagerInterface) {
				userService.On("Create", mock.AnythingOfType("models.UserCreateRequest")).Return((*models.User)(nil), errors.New("email already exists"))
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "registration_failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockUserService := new(MockUserServiceInterface)
			mockRoleService := new(MockRoleServiceInterface)
			mockJWTManager := new(MockJWTManagerInterface)

			tt.mockSetup(mockUserService, mockRoleService, mockJWTManager)
			// Create handler
			handler := handlers.NewAuthHandler(
				mockUserService,
				mockRoleService,
				mockJWTManager,
			)

			// Setup router
			router := gin.New()
			router.POST("/register", handler.Register)

			// Create request
			jsonBody, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response.Error)
			}

			// Verify mock expectations
			mockUserService.AssertExpectations(t)
			mockRoleService.AssertExpectations(t)
			mockJWTManager.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		request        models.LoginRequest
		mockSetup      func(*MockUserServiceInterface, *MockJWTManagerInterface)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful login",
			request: models.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(userService *MockUserServiceInterface, jwtManager *MockJWTManagerInterface) {
				user := &models.User{
					ID:       uuid.New(),
					Email:    "test@example.com",
					IsActive: true,
				}
				userWithRoles := &models.User{
					ID:       user.ID,
					Email:    "test@example.com",
					IsActive: true,
				}

				userService.On("ValidateCredentials", "test@example.com", "password123").Return(user, nil)
				userService.On("GetByIDWithRoles", user.ID.String()).Return(userWithRoles, nil)
				jwtManager.On("GenerateAccessToken", user.ID, user.Email).Return("access_token", nil)
				jwtManager.On("GenerateRefreshToken", user.ID, user.Email).Return("refresh_token", nil)
				jwtManager.On("GetTokenDuration", auth.AccessToken).Return(15 * time.Minute)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid credentials",
			request: models.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(userService *MockUserServiceInterface, jwtManager *MockJWTManagerInterface) {
				userService.On("ValidateCredentials", "test@example.com", "wrongpassword").Return((*models.User)(nil), errors.New("invalid credentials"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "authentication_failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockUserService := new(MockUserServiceInterface)
			mockJWTManager := new(MockJWTManagerInterface)

			tt.mockSetup(mockUserService, mockJWTManager)

			// Create handler
			handler := handlers.NewAuthHandler(mockUserService, nil, mockJWTManager)

			// Setup router
			router := gin.New()
			router.POST("/login", handler.Login)

			// Create request
			jsonBody, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response.Error)
			}

			// Verify mock expectations
			mockUserService.AssertExpectations(t)
			mockJWTManager.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		request        models.RefreshTokenRequest
		mockSetup      func(*MockUserServiceInterface, *MockJWTManagerInterface)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful token refresh",
			request: models.RefreshTokenRequest{
				RefreshToken: "valid_refresh_token",
			},
			mockSetup: func(userService *MockUserServiceInterface, jwtManager *MockJWTManagerInterface) {
				userID := uuid.New()
				claims := &auth.JWTClaims{
					UserID: userID,
					Email:  "test@example.com",
				}
				user := &models.User{
					ID:       userID,
					Email:    "test@example.com",
					IsActive: true,
				}

				jwtManager.On("ValidateToken", "valid_refresh_token", auth.RefreshToken).Return(claims, nil)
				userService.On("GetByIDWithRoles", userID.String()).Return(user, nil)
				jwtManager.On("GenerateAccessToken", userID, user.Email).Return("new_access_token", nil)
				jwtManager.On("GenerateRefreshToken", userID, user.Email).Return("new_refresh_token", nil)
				jwtManager.On("GetTokenDuration", auth.AccessToken).Return(15 * time.Minute)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid refresh token",
			request: models.RefreshTokenRequest{
				RefreshToken: "invalid_token",
			},
			mockSetup: func(userService *MockUserServiceInterface, jwtManager *MockJWTManagerInterface) {
				jwtManager.On("ValidateToken", "invalid_token", auth.RefreshToken).Return((*auth.JWTClaims)(nil), errors.New("invalid token"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "invalid_token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockUserService := new(MockUserServiceInterface)
			mockJWTManager := new(MockJWTManagerInterface)

			tt.mockSetup(mockUserService, mockJWTManager)

			// Create handler
			handler := handlers.NewAuthHandler(mockUserService, nil, mockJWTManager)

			// Setup router
			router := gin.New()
			router.POST("/refresh", handler.RefreshToken)

			// Create request
			jsonBody, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response.Error)
			}

			// Verify mock expectations
			mockUserService.AssertExpectations(t)
			mockJWTManager.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_ChangePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		request        models.ChangePasswordRequest
		mockSetup      func(*MockUserServiceInterface)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successful password change",
			userID: "user-123",
			request: models.ChangePasswordRequest{
				CurrentPassword: "oldpassword",
				NewPassword:     "newpassword",
			},
			mockSetup: func(userService *MockUserServiceInterface) {
				userService.On("ChangePassword", "user-123", "oldpassword", "newpassword").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "invalid current password",
			userID: "user-123",
			request: models.ChangePasswordRequest{
				CurrentPassword: "wrongpassword",
				NewPassword:     "newpassword",
			},
			mockSetup: func(userService *MockUserServiceInterface) {
				userService.On("ChangePassword", "user-123", "wrongpassword", "newpassword").Return(errors.New("current password is incorrect"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "password_change_failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockUserService := new(MockUserServiceInterface)

			tt.mockSetup(mockUserService)

			// Create handler
			handler := handlers.NewAuthHandler(mockUserService, nil, nil)

			// Setup router with middleware mock
			router := gin.New()
			router.POST("/change-password", func(c *gin.Context) {
				// Mock middleware setting user ID
				c.Set("user_id", tt.userID)
				handler.ChangePassword(c)
			})

			// Create request
			jsonBody, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/change-password", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response.Error)
			}

			// Verify mock expectations
			mockUserService.AssertExpectations(t)
		})
	}
}
