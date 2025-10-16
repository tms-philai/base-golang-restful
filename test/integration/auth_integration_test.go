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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// AuthIntegrationTestSuite is the test suite for auth integration tests
type AuthIntegrationTestSuite struct {
	suite.Suite
	router         *gin.Engine
	dbSetup        *setup.TestDatabaseSetup
	userService    *services.UserService
	roleService    *services.RoleService
	jwtManager     *auth.JWTManager
	authHandler    *handlers.AuthHandler
	authMiddleware *middleware.AuthMiddleware
}

// T returns the testing.T instance for the current test
func (suite *AuthIntegrationTestSuite) T() *testing.T {
	return suite.Suite.T()
}

// SetupSuite sets up the test suite
func (suite *AuthIntegrationTestSuite) SetupSuite() {
	// Setup test database
	suite.dbSetup = setup.SetupTestEnvironment(suite.T())

	// Load test configuration
	testConfig, err := config.LoadTestConfig()
	assert.NoError(suite.T(), err)

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
	suite.authHandler = handlers.NewAuthHandler(suite.userService, suite.roleService, suite.jwtManager)

	// Setup router
	suite.router = routes.SetupRouter(routes.RouterConfig{
		AuthHandler:    suite.authHandler,
		AuthMiddleware: suite.authMiddleware,
	})
}

// TearDownSuite tears down the test suite
func (suite *AuthIntegrationTestSuite) TearDownSuite() {
	setup.CleanupTestEnvironment(suite.T(), suite.dbSetup)
}

// SetupTest sets up each test
func (suite *AuthIntegrationTestSuite) SetupTest() {
	// Clean up test data before each test
	suite.dbSetup.TruncateTestTables()
}

// TestRegister tests user registration
func (suite *AuthIntegrationTestSuite) TestRegister() {
	tests := []struct {
		name           string
		request        helpers.TestUser
		expectedStatus int
		expectError    bool
	}{
		{
			name: "successful registration",
			request: helpers.TestUser{
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
			name: "registration with existing email",
			request: helpers.TestUser{
				Username:  "existinguser",
				Email:     "existing@example.com",
				Password:  "password123",
				FirstName: "Existing",
				LastName:  "User",
			},
			expectedStatus: http.StatusConflict,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Create request
			req, err := helpers.CreateJSONRequest("POST", "/api/v1/auth/register", tt.request)
			assert.NoError(suite.T(), err)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			suite.router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if !tt.expectError {
				helpers.AssertAuthResponse(suite.T(), w, tt.expectedStatus)
			} else {
				helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "registration_failed")
			}
		})
	}
}

func (suite *AuthIntegrationTestSuite) Run(name string, f func()) {
	panic("unimplemented")
}

// TestLogin tests user login
func (suite *AuthIntegrationTestSuite) TestLogin() {
	// First, create a test user
	testUser := helpers.CreateTestUser()
	user, err := suite.userService.Create(models.UserCreateRequest{
		Email:     testUser.Email,
		Password:  testUser.Password,
		FirstName: testUser.FirstName,
		LastName:  testUser.LastName,
	})
	assert.NoError(suite.T(), err)

	// Assign default role
	err = suite.roleService.AssignDefaultRole(user.ID)
	assert.NoError(suite.T(), err)

	tests := []struct {
		name           string
		email          string
		password       string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "successful login",
			email:          testUser.Email,
			password:       testUser.Password,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "invalid credentials",
			email:          testUser.Email,
			password:       "wrongpassword",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:           "non-existent user",
			email:          "nonexistent@example.com",
			password:       "password123",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			loginReq := struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Email:    tt.email,
				Password: tt.password,
			}

			// Create request
			req, err := helpers.CreateJSONRequest("POST", "/api/v1/auth/login", loginReq)
			assert.NoError(suite.T(), err)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			suite.router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if !tt.expectError {
				helpers.AssertAuthResponse(suite.T(), w, tt.expectedStatus)
			} else {
				helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "authentication_failed")
			}
		})
	}
}

// TestRefreshToken tests token refresh
func (suite *AuthIntegrationTestSuite) TestRefreshToken() {
	// First, create a test user and get tokens
	testUser := helpers.CreateTestUser()
	user, err := suite.userService.Create(models.UserCreateRequest{
		Email:     testUser.Email,
		Password:  testUser.Password,
		FirstName: testUser.FirstName,
		LastName:  testUser.LastName,
	})
	assert.NoError(suite.T(), err)

	// Generate refresh token
	refreshToken, err := suite.jwtManager.GenerateRefreshToken(user.ID, user.Email)
	assert.NoError(suite.T(), err)

	tests := []struct {
		name           string
		refreshToken   string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "successful token refresh",
			refreshToken:   refreshToken,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "invalid refresh token",
			refreshToken:   "invalid-token",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			refreshReq := struct {
				RefreshToken string `json:"refresh_token"`
			}{
				RefreshToken: tt.refreshToken,
			}

			// Create request
			req, err := helpers.CreateJSONRequest("POST", "/api/v1/auth/refresh", refreshReq)
			assert.NoError(suite.T(), err)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			suite.router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if !tt.expectError {
				helpers.AssertAuthResponse(suite.T(), w, tt.expectedStatus)
			} else {
				helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "invalid_token")
			}
		})
	}
}

// TestGetProfile tests getting user profile
func (suite *AuthIntegrationTestSuite) TestGetProfile() {
	// First, create a test user and get tokens
	testUser := helpers.CreateTestUser()
	user, err := suite.userService.Create(models.UserCreateRequest{
		Email:     testUser.Email,
		Password:  testUser.Password,
		FirstName: testUser.FirstName,
		LastName:  testUser.LastName,
	})
	assert.NoError(suite.T(), err)

	// Generate access token
	accessToken, err := suite.jwtManager.GenerateAccessToken(user.ID, user.Email)
	assert.NoError(suite.T(), err)

	tests := []struct {
		name           string
		accessToken    string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "successful get profile",
			accessToken:    accessToken,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "unauthorized - no token",
			accessToken:    "",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:           "unauthorized - invalid token",
			accessToken:    "invalid-token",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Create request
			req, err := http.NewRequest("GET", "/api/v1/auth/profile", nil)
			assert.NoError(suite.T(), err)

			// Add authorization header if token provided
			if tt.accessToken != "" {
				req.Header.Set("Authorization", "Bearer "+tt.accessToken)
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
				helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "unauthorized")
			}
		})
	}
}

// TestChangePassword tests password change
func (suite *AuthIntegrationTestSuite) TestChangePassword() {
	// First, create a test user and get tokens
	testUser := helpers.CreateTestUser()
	user, err := suite.userService.Create(models.UserCreateRequest{
		Email:     testUser.Email,
		Password:  testUser.Password,
		FirstName: testUser.FirstName,
		LastName:  testUser.LastName,
	})
	assert.NoError(suite.T(), err)

	// Generate access token
	accessToken, err := suite.jwtManager.GenerateAccessToken(user.ID, user.Email)
	assert.NoError(suite.T(), err)

	tests := []struct {
		name            string
		accessToken     string
		currentPassword string
		newPassword     string
		expectedStatus  int
		expectError     bool
	}{
		{
			name:            "successful password change",
			accessToken:     accessToken,
			currentPassword: testUser.Password,
			newPassword:     "newpassword123",
			expectedStatus:  http.StatusOK,
			expectError:     false,
		},
		{
			name:            "invalid current password",
			accessToken:     accessToken,
			currentPassword: "wrongpassword",
			newPassword:     "newpassword123",
			expectedStatus:  http.StatusBadRequest,
			expectError:     true,
		},
		{
			name:            "unauthorized - no token",
			accessToken:     "",
			currentPassword: testUser.Password,
			newPassword:     "newpassword123",
			expectedStatus:  http.StatusUnauthorized,
			expectError:     true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			changePasswordReq := struct {
				CurrentPassword string `json:"current_password"`
				NewPassword     string `json:"new_password"`
			}{
				CurrentPassword: tt.currentPassword,
				NewPassword:     tt.newPassword,
			}

			// Create request
			req, err := helpers.CreateJSONRequest("POST", "/api/v1/auth/change-password", changePasswordReq)
			assert.NoError(suite.T(), err)

			// Add authorization header if token provided
			if tt.accessToken != "" {
				req.Header.Set("Authorization", "Bearer "+tt.accessToken)
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
				} else {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "password_change_failed")
				}
			}
		})
	}
}

// TestAuthIntegrationSuite runs the auth integration test suite
func TestAuthIntegrationSuite(t *testing.T) {
	suite.Run(t, new(AuthIntegrationTestSuite))
}
