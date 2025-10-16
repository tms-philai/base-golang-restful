package unit

import (
	"base-gin/internal/app/handlers"
	"base-gin/internal/domain/models"
	"base-gin/test/helpers"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(req models.UserCreateRequest) (*models.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetByIDWithRoles(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) List(page, pageSize int) ([]*models.User, models.PaginationMetadata, error) {
	args := m.Called(page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(models.PaginationMetadata), args.Error(2)
	}
	return args.Get(0).([]*models.User), args.Get(1).(models.PaginationMetadata), args.Error(2)
}

func (m *MockUserService) Update(id string, req models.UserUpdateRequest) (*models.User, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) ValidateCredentials(email, password string) (*models.User, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) ChangePassword(userID, oldPassword, newPassword string) error {
	args := m.Called(userID, oldPassword, newPassword)
	return args.Error(0)
}

func TestUserHandler_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		request        models.UserCreateRequest
		mockSetup      func(*MockUserService)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful user creation",
			request: models.UserCreateRequest{
				Email:     "newuser@example.com",
				Password:  "password123",
				FirstName: "New",
				LastName:  "User",
			},
			mockSetup: func(userService *MockUserService) {
				user := &models.User{
					ID:        uuid.New(),
					Email:     "newuser@example.com",
					FirstName: "New",
					LastName:  "User",
					IsActive:  true,
				}
				userService.On("Create", mock.AnythingOfType("models.UserCreateRequest")).Return(user, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "validation error - missing email",
			request: models.UserCreateRequest{
				Username:  "newuser",
				Password:  "password123",
				FirstName: "New",
				LastName:  "User",
			},
			mockSetup:      func(*MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name: "user creation failed - email already exists",
			request: models.UserCreateRequest{
				Username:  "newuser",
				Email:     "existing@example.com",
				Password:  "password123",
				FirstName: "New",
				LastName:  "User",
			},
			mockSetup: func(userService *MockUserService) {
				userService.On("Create", mock.AnythingOfType("models.UserCreateRequest")).Return((*models.User)(nil), errors.New("email already exists"))
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "user_creation_failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockUserService := new(MockUserService)
			tt.mockSetup(mockUserService)

			// Create handler
			handler := handlers.NewUserHandler(mockUserService)

			// Setup router
			router := gin.New()
			router.POST("/users", handler.CreateUser)

			// Create request
			jsonBody, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
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

func TestUserHandler_GetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*MockUserService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successful get user",
			userID: "user-123",
			mockSetup: func(userService *MockUserService) {
				user := &models.User{
					ID:        uuid.New(),
					Email:     "test@example.com",
					FirstName: "Test",
					LastName:  "User",
					IsActive:  true,
				}
				userService.On("GetByID", "user-123").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "user not found",
			userID: "nonexistent",
			mockSetup: func(userService *MockUserService) {
				userService.On("GetByID", "nonexistent").Return((*models.User)(nil), errors.New("user not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "user_not_found",
		},
		{
			name:           "missing user ID",
			userID:         "",
			mockSetup:      func(*MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockUserService := new(MockUserService)
			tt.mockSetup(mockUserService)

			// Create handler
			handler := handlers.NewUserHandler(mockUserService)

			// Setup router
			router := gin.New()
			router.GET("/users/:id", handler.GetUser)

			// Create request
			url := "/users/" + tt.userID
			req, _ := http.NewRequest("GET", url, nil)

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

func TestUserHandler_UpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		currentUserID  string
		request        models.UserUpdateRequest
		mockSetup      func(*MockUserService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:          "successful update by admin",
			userID:        "user-123",
			currentUserID: "admin-456",
			request: models.UserUpdateRequest{
				FirstName: helpers.StringPtr("Updated"),
				LastName:  helpers.StringPtr("Name"),
			},
			mockSetup: func(userService *MockUserService) {
				user := &models.User{
					ID:        uuid.New(),
					Email:     "test@example.com",
					FirstName: "Updated",
					LastName:  "Name",
					IsActive:  true,
				}
				userService.On("Update", "user-123", mock.AnythingOfType("models.UserUpdateRequest")).Return(user, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "user not found",
			userID:        "nonexistent",
			currentUserID: "admin-456",
			request: models.UserUpdateRequest{
				FirstName: helpers.StringPtr("Updated"),
				LastName:  helpers.StringPtr("Name"),
			},
			mockSetup: func(userService *MockUserService) {
				userService.On("Update", "user-123", mock.AnythingOfType("models.UserUpdateRequest")).Return((*models.User)(nil), errors.New("user not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "user_not_found",
		},
		{
			name:           "missing user ID",
			userID:         "",
			currentUserID:  "admin-456",
			request:        models.UserUpdateRequest{},
			mockSetup:      func(*MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name:           "missing current user ID",
			userID:         "user-123",
			currentUserID:  "",
			request:        models.UserUpdateRequest{},
			mockSetup:      func(*MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name:           "missing request",
			userID:         "user-123",
			currentUserID:  "admin-456",
			request:        models.UserUpdateRequest{},
			mockSetup:      func(*MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name:           "missing request fields",
			userID:         "user-123",
			currentUserID:  "admin-456",
			request:        models.UserUpdateRequest{},
			mockSetup:      func(*MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name:           "invalid user ID",
			userID:         "invalid-id",
			currentUserID:  "admin-456",
			request:        models.UserUpdateRequest{},
			mockSetup:      func(*MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockUserService := new(MockUserService)
			tt.mockSetup(mockUserService)

			// Create handler
			handler := handlers.NewUserHandler(mockUserService)

			// Setup router with middleware mock
			router := gin.New()
			router.PUT("/users/:id", func(c *gin.Context) {
				// Mock middleware setting current user
				user := &models.User{
					ID: uuid.MustParse(tt.currentUserID),
				}
				// Add admin role if current user is admin
				if tt.currentUserID == "admin-456" {
					user.Roles = []models.Role{{Name: "admin"}}
				}
				c.Set("current_user", user)
				handler.UpdateUser(c)
			})

			// Create request
			jsonBody, _ := json.Marshal(tt.request)
			url := "/users/" + tt.userID
			req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
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

func TestUserHandler_DeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*MockUserService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successful delete user",
			userID: "user-123",
			mockSetup: func(userService *MockUserService) {
				userService.On("Delete", "user-123").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "user not found",
			userID: "nonexistent",
			mockSetup: func(userService *MockUserService) {
				userService.On("Delete", "nonexistent").Return(errors.New("user not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "user_deletion_failed",
		},
		{
			name:           "missing user ID",
			userID:         "",
			mockSetup:      func(*MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockUserService := new(MockUserService)
			tt.mockSetup(mockUserService)

			// Create handler
			handler := handlers.NewUserHandler(mockUserService)

			// Setup router
			router := gin.New()
			router.DELETE("/users/:id", handler.DeleteUser)

			// Create request
			url := "/users/" + tt.userID
			req, _ := http.NewRequest("DELETE", url, nil)

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

func TestUserHandler_ListUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*MockUserService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:        "successful list users",
			queryParams: "?page=1&limit=10",
			mockSetup: func(userService *MockUserService) {
				users := []models.User{
					{
						ID:        uuid.New(),
						Email:     "user1@example.com",
						FirstName: "User",
						LastName:  "One",
						IsActive:  true,
					},
					{
						ID:        uuid.New(),
						Email:     "user2@example.com",
						FirstName: "User",
						LastName:  "Two",
						IsActive:  true,
					},
				}
				pagination := models.PaginationMetadata{
					Page:       1,
					Limit:      10,
					Total:      2,
					TotalPages: 1,
				}
				userService.On("List", 1, 10).Return(users, pagination, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "list users with custom pagination",
			queryParams: "?page=2&limit=5",
			mockSetup: func(userService *MockUserService) {
				users := []models.User{}
				pagination := models.PaginationMetadata{
					Page:       2,
					Limit:      5,
					Total:      0,
					TotalPages: 0,
				}
				userService.On("List", 2, 5).Return(users, pagination, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "list users service error",
			queryParams: "?page=1&limit=10",
			mockSetup: func(userService *MockUserService) {
				userService.On("List", 1, 10).Return(([]models.User)(nil), models.PaginationMetadata{}, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "list_users_failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockUserService := new(MockUserService)
			tt.mockSetup(mockUserService)

			// Create handler
			handler := handlers.NewUserHandler(mockUserService)

			// Setup router
			router := gin.New()
			router.GET("/users", handler.ListUsers)

			// Create request
			url := "/users" + tt.queryParams
			req, _ := http.NewRequest("GET", url, nil)

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
