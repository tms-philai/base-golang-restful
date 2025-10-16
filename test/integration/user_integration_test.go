package integration

import (
	"base-gin/internal/app/handlers"
	"base-gin/internal/app/middleware"
	"base-gin/internal/app/routes"
	"base-gin/internal/domain/models"
	"base-gin/internal/domain/repository"
	"base-gin/internal/domain/services"
	"base-gin/internal/pkg/auth"
	"base-gin/test/config"
	"base-gin/test/helpers"
	"base-gin/test/setup"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserIntegrationTestSuite is the test suite for user integration tests
type UserIntegrationTestSuite struct {
	suite.Suite
	router         *gin.Engine
	dbSetup        *setup.TestDatabaseSetup
	userService    *services.UserService
	roleService    *services.RoleService
	jwtManager     *auth.JWTManager
	userHandler    *handlers.UserHandler
	authMiddleware *middleware.AuthMiddleware
	adminToken     string
	userToken      string
	adminUser      *models.User
	regularUser    *models.User
}

// T returns the testing.T instance for the current test
func (suite *UserIntegrationTestSuite) T() *testing.T {
	return suite.Suite.T()
}

// SetupSuite sets up the test suite
func (suite *UserIntegrationTestSuite) SetupSuite() {
	// Setup test database
	suite.dbSetup = setup.SetupTestEnvironment(suite.T())

	// Load test configuration
	testConfig, err := config.LoadTestConfig()
	suite.Require().NoError(err)

	// Initialize repositories
	userRepo := repository.NewUserRepository(setup.GetTestDB())
	roleRepo := repository.NewRoleRepository(setup.GetTestDB())
	permissionRepo := repository.NewPermissionRepository(setup.GetTestDB())

	// Initialize JWT manager
	suite.jwtManager = auth.NewJWTManager(auth.JWTConfig{
		SecretKey:            testConfig.JWT.SecretKey,
		AccessTokenDuration:  testConfig.JWT.AccessTokenDuration,
		RefreshTokenDuration: testConfig.JWT.RefreshTokenDuration,
		Issuer:               "base-golang-restful-test",
	})

	// Initialize services
	suite.userService = services.NewUserService(userRepo)
	suite.roleService = services.NewRoleService(setup.GetTestDB(), roleRepo, userRepo, permissionRepo)

	// Initialize middleware and handlers
	suite.authMiddleware = middleware.NewAuthMiddleware(suite.jwtManager, suite.userService)
	suite.userHandler = handlers.NewUserHandler(suite.userService)

	// Setup router
	suite.router = routes.SetupRouter(routes.RouterConfig{
		AuthHandler:    nil, // We'll set this later
		AuthMiddleware: suite.authMiddleware,
		UserHandler:    suite.userHandler,
	})
}

// TearDownSuite tears down the test suite
func (suite *UserIntegrationTestSuite) TearDownSuite() {
	setup.CleanupTestEnvironment(suite.T(), suite.dbSetup)
}

// SetupTest sets up each test
func (suite *UserIntegrationTestSuite) SetupTest() {
	// Clean up test data before each test
	suite.dbSetup.TruncateTestTables()

	// Create test users
	suite.createTestUsers()
}

// createTestUsers creates test users for the suite
func (suite *UserIntegrationTestSuite) createTestUsers() {
	// Create admin user
	adminUser, err := suite.userService.Create(models.UserCreateRequest{
		Username:  "admin",
		Email:     "admin@example.com",
		Password:  "admin123",
		FirstName: "Admin",
		LastName:  "User",
	})
	suite.Require().NoError(err)

	// Assign admin role
	// Get admin role ID and assign it
	adminRole, err := suite.roleService.GetRoleByName("admin")
	suite.Require().NoError(err)
	err = suite.roleService.AssignRole(adminUser.ID, adminRole.(models.Role).ID)
	suite.Require().NoError(err)

	// Create regular user
	regularUser, err := suite.userService.Create(models.UserCreateRequest{
		Username:  "user",
		Email:     "user@example.com",
		Password:  "user123",
		FirstName: "Regular",
		LastName:  "User",
	})
	suite.Require().NoError(err)

	// Assign user role
	// Get user role ID and assign it
	userRole, err := suite.roleService.GetRoleByName("user")
	suite.Require().NoError(err)
	err = suite.roleService.AssignRole(regularUser.ID, userRole.(models.Role).ID)
	suite.Require().NoError(err)

	// Generate tokens
	suite.adminToken, err = suite.jwtManager.GenerateAccessToken(adminUser.ID, adminUser.Email)
	suite.Require().NoError(err)

	suite.userToken, err = suite.jwtManager.GenerateAccessToken(regularUser.ID, regularUser.Email)
	suite.Require().NoError(err)

	suite.adminUser = adminUser
	suite.regularUser = regularUser
}

// TestCreateUser tests user creation (admin only)
func (suite *UserIntegrationTestSuite) TestCreateUser() {
	tests := []struct {
		name           string
		token          string
		request        models.UserCreateRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name:  "successful user creation by admin",
			token: suite.adminToken,
			request: models.UserCreateRequest{
				Username:  "newuser",
				Email:     "newuser@example.com",
				Password:  "password123",
				FirstName: "New",
				LastName:  "User",
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name:  "unauthorized - regular user trying to create user",
			token: suite.userToken,
			request: models.UserCreateRequest{
				Username:  "newuser",
				Email:     "newuser@example.com",
				Password:  "password123",
				FirstName: "New",
				LastName:  "User",
			},
			expectedStatus: http.StatusForbidden,
			expectError:    true,
		},
		{
			name:           "unauthorized - no token",
			token:          "",
			request:        models.UserCreateRequest{},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:  "validation error - missing email",
			token: suite.adminToken,
			request: models.UserCreateRequest{
				Username:  "newuser",
				Password:  "password123",
				FirstName: "New",
				LastName:  "User",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Create request
			req, err := helpers.CreateJSONRequest("POST", "/api/v1/users", tt.request)
			suite.Require().NoError(err)

			// Add authorization header if token provided
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			suite.router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if !tt.expectError {
				helpers.AssertUserResponse(suite.T(), w, tt.expectedStatus)
			} else {
				if tt.expectedStatus == http.StatusUnauthorized {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "unauthorized")
				} else if tt.expectedStatus == http.StatusForbidden {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "forbidden")
				} else {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "validation_error")
				}
			}
		})
	}
}

// TestGetUser tests getting a user by ID
func (suite *UserIntegrationTestSuite) TestGetUser() {
	tests := []struct {
		name           string
		userID         string
		token          string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "successful get user by ID",
			userID:         suite.regularUser.ID.String(),
			token:          suite.userToken,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "successful get user by ID without token (optional auth)",
			userID:         suite.regularUser.ID.String(),
			token:          "",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "user not found",
			userID:         "00000000-0000-0000-0000-000000000000",
			token:          suite.userToken,
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:           "invalid user ID format",
			userID:         "invalid-id",
			token:          suite.userToken,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Create request
			req, err := http.NewRequest("GET", "/api/v1/users/"+tt.userID, nil)
			suite.Require().NoError(err)

			// Add authorization header if token provided
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			suite.router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if !tt.expectError {
				helpers.AssertUserResponse(suite.T(), w, tt.expectedStatus)
			} else {
				if tt.expectedStatus == http.StatusNotFound {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "user_not_found")
				} else {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "validation_error")
				}
			}
		})
	}
}

// TestUpdateUser tests updating a user
func (suite *UserIntegrationTestSuite) TestUpdateUser() {
	tests := []struct {
		name           string
		userID         string
		token          string
		request        models.UserUpdateRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name:   "successful update by admin",
			userID: suite.regularUser.ID.String(),
			token:  suite.adminToken,
			request: models.UserUpdateRequest{
				FirstName: helpers.StringPtr("Updated"),
				LastName:  helpers.StringPtr("Name"),
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:   "successful update by user of their own profile",
			userID: suite.regularUser.ID.String(),
			token:  suite.userToken,
			request: models.UserUpdateRequest{
				FirstName: helpers.StringPtr("Updated"),
				LastName:  helpers.StringPtr("Name"),
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:   "forbidden - user trying to update another user",
			userID: suite.adminUser.ID.String(),
			token:  suite.userToken,
			request: models.UserUpdateRequest{
				FirstName: helpers.StringPtr("Updated"),
				LastName:  helpers.StringPtr("Name"),
			},
			expectedStatus: http.StatusForbidden,
			expectError:    true,
		},
		{
			name:           "unauthorized - no token",
			userID:         suite.regularUser.ID.String(),
			token:          "",
			request:        models.UserUpdateRequest{},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:   "user not found",
			userID: "00000000-0000-0000-0000-000000000000",
			token:  suite.adminToken,
			request: models.UserUpdateRequest{
				FirstName: helpers.StringPtr("Updated"),
				LastName:  helpers.StringPtr("Name"),
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Create request
			req, err := helpers.CreateJSONRequest("PUT", "/api/v1/users/"+tt.userID, tt.request)
			suite.Require().NoError(err)

			// Add authorization header if token provided
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			suite.router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if !tt.expectError {
				helpers.AssertUserResponse(suite.T(), w, tt.expectedStatus)
			} else {
				if tt.expectedStatus == http.StatusUnauthorized {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "unauthorized")
				} else if tt.expectedStatus == http.StatusForbidden {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "forbidden")
				} else if tt.expectedStatus == http.StatusNotFound {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "user_update_failed")
				} else {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "validation_error")
				}
			}
		})
	}
}

// TestDeleteUser tests deleting a user (admin only)
func (suite *UserIntegrationTestSuite) TestDeleteUser() {
	tests := []struct {
		name           string
		userID         string
		token          string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "successful delete by admin",
			userID:         suite.regularUser.ID.String(),
			token:          suite.adminToken,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "forbidden - regular user trying to delete user",
			userID:         suite.adminUser.ID.String(),
			token:          suite.userToken,
			expectedStatus: http.StatusForbidden,
			expectError:    true,
		},
		{
			name:           "unauthorized - no token",
			userID:         suite.regularUser.ID.String(),
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:           "user not found",
			userID:         "00000000-0000-0000-0000-000000000000",
			token:          suite.adminToken,
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Create request
			req, err := http.NewRequest("DELETE", "/api/v1/users/"+tt.userID, nil)
			suite.Require().NoError(err)

			// Add authorization header if token provided
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			suite.router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if !tt.expectError {
				helpers.AssertSuccessResponse(suite.T(), w, tt.expectedStatus)
			} else {
				if tt.expectedStatus == http.StatusUnauthorized {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "unauthorized")
				} else if tt.expectedStatus == http.StatusForbidden {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "forbidden")
				} else if tt.expectedStatus == http.StatusNotFound {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "user_deletion_failed")
				} else {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "validation_error")
				}
			}
		})
	}
}

// TestListUsers tests listing users (admin only)
func (suite *UserIntegrationTestSuite) TestListUsers() {
	tests := []struct {
		name           string
		token          string
		queryParams    string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "successful list users by admin",
			token:          suite.adminToken,
			queryParams:    "?page=1&limit=10",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "successful list users with custom pagination",
			token:          suite.adminToken,
			queryParams:    "?page=2&limit=5",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "forbidden - regular user trying to list users",
			token:          suite.userToken,
			queryParams:    "?page=1&limit=10",
			expectedStatus: http.StatusForbidden,
			expectError:    true,
		},
		{
			name:           "unauthorized - no token",
			token:          "",
			queryParams:    "?page=1&limit=10",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Create request
			req, err := http.NewRequest("GET", "/api/v1/users"+tt.queryParams, nil)
			suite.Require().NoError(err)

			// Add authorization header if token provided
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			suite.router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if !tt.expectError {
				helpers.AssertPaginationResponse(suite.T(), w, tt.expectedStatus)

				// Verify response structure
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				suite.Require().NoError(err)
				assert.Contains(suite.T(), response, "users")
				assert.Contains(suite.T(), response, "pagination")
			} else {
				if tt.expectedStatus == http.StatusUnauthorized {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "unauthorized")
				} else if tt.expectedStatus == http.StatusForbidden {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "forbidden")
				}
			}
		})
	}
}

// TestUserIntegrationSuite runs the user integration test suite
func TestUserIntegrationSuite(t *testing.T) {
	suite.Run(t, new(UserIntegrationTestSuite))
}
