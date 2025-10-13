package handlers

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"base-golang-restful-app/handlers"
	"base-golang-restful-app/models"
	"base-golang-restful-app/tests/mocks"
	"base-golang-restful-app/tests/testhelpers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// AuthHandlerTestSuite defines the test suite for AuthHandler
type AuthHandlerTestSuite struct {
	suite.Suite
	helper      *testhelpers.TestHelper
	mockUserSvc *mocks.MockUserService
	mockJWTSvc  *mocks.MockJWTService
	authHandler *handlers.AuthHandler
	router      *gin.Engine
}

// SetupTest sets up the test environment before each test
func (suite *AuthHandlerTestSuite) SetupTest() {
	suite.helper = testhelpers.NewTestHelper()
	suite.mockUserSvc = new(mocks.MockUserService)
	suite.mockJWTSvc = new(mocks.MockJWTService)
	suite.authHandler = handlers.NewAuthHandler(suite.mockUserSvc, suite.mockJWTSvc)

	// Setup router
	suite.router = gin.New()
	auth := suite.router.Group("/auth")
	{
		auth.POST("/register", suite.authHandler.Register)
		auth.POST("/login", suite.authHandler.Login)
		auth.POST("/refresh", suite.authHandler.RefreshToken)
	}
	suite.helper.Router = suite.router
}

// TearDownTest cleans up after each test
func (suite *AuthHandlerTestSuite) TearDownTest() {
	suite.mockUserSvc.AssertExpectations(suite.T())
	suite.mockJWTSvc.AssertExpectations(suite.T())
}

// TestRegisterSuccess tests successful user registration
func (suite *AuthHandlerTestSuite) TestRegisterSuccess() {
	// Arrange
	registerReq := models.RegisterRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	expectedUser := &models.User{
		ID:        "user-123",
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Role:      "user",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	expectedTokenPair := &models.TokenPair{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    3600,
	}

	// Setup mocks
	suite.mockUserSvc.On("Create", mock.AnythingOfType("models.UserCreateRequest")).Return(expectedUser, nil)
	suite.mockJWTSvc.On("GenerateTokenPair", expectedUser).Return(expectedTokenPair, nil)

	// Act
	w := suite.helper.MakeRequest("POST", "/auth/register", registerReq, nil)

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response models.AuthResponse
	err := suite.helper.ParseJSONResponse(w, &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedUser.Username, response.User.Username)
	assert.Equal(suite.T(), expectedUser.Email, response.User.Email)
	assert.Equal(suite.T(), expectedTokenPair.AccessToken, response.AccessToken)
	assert.Equal(suite.T(), "Bearer", response.TokenType)
}

// TestRegisterValidationError tests registration with validation errors
func (suite *AuthHandlerTestSuite) TestRegisterValidationError() {
	// Arrange
	invalidReq := models.RegisterRequest{
		Username: "ab", // Too short
		Email:    "invalid-email",
		Password: "123", // Too short
	}

	// Act
	w := suite.helper.MakeRequest("POST", "/auth/register", invalidReq, nil)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	suite.helper.AssertErrorResponse(suite.T(), w, http.StatusBadRequest, "validation_error")
}

// TestRegisterUserExists tests registration when user already exists
func (suite *AuthHandlerTestSuite) TestRegisterUserExists() {
	// Arrange
	registerReq := models.RegisterRequest{
		Username:  "existinguser",
		Email:     "existing@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	// Setup mocks
	suite.mockUserSvc.On("Create", mock.AnythingOfType("models.UserCreateRequest")).Return((*models.User)(nil), errors.New("username already exists"))

	// Act
	w := suite.helper.MakeRequest("POST", "/auth/register", registerReq, nil)

	// Assert
	assert.Equal(suite.T(), http.StatusConflict, w.Code)
	suite.helper.AssertErrorResponse(suite.T(), w, http.StatusConflict, "registration_failed")
}

// TestLoginSuccess tests successful user login
func (suite *AuthHandlerTestSuite) TestLoginSuccess() {
	// Arrange
	loginReq := models.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	expectedUser := &models.User{
		ID:       "user-123",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
		IsActive: true,
	}

	expectedTokenPair := &models.TokenPair{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    3600,
	}

	// Setup mocks
	suite.mockUserSvc.On("ValidateCredentials", "testuser", "password123").Return(expectedUser, nil)
	suite.mockJWTSvc.On("GenerateTokenPair", expectedUser).Return(expectedTokenPair, nil)

	// Act
	w := suite.helper.MakeRequest("POST", "/auth/login", loginReq, nil)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.AuthResponse
	err := suite.helper.ParseJSONResponse(w, &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedUser.Username, response.User.Username)
	assert.Equal(suite.T(), expectedTokenPair.AccessToken, response.AccessToken)
}

// TestLoginInvalidCredentials tests login with invalid credentials
func (suite *AuthHandlerTestSuite) TestLoginInvalidCredentials() {
	// Arrange
	loginReq := models.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	// Setup mocks
	suite.mockUserSvc.On("ValidateCredentials", "testuser", "wrongpassword").Return((*models.User)(nil), errors.New("invalid credentials"))

	// Act
	w := suite.helper.MakeRequest("POST", "/auth/login", loginReq, nil)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	suite.helper.AssertErrorResponse(suite.T(), w, http.StatusUnauthorized, "authentication_failed")
}

// TestRefreshTokenSuccess tests successful token refresh
func (suite *AuthHandlerTestSuite) TestRefreshTokenSuccess() {
	// Arrange
	refreshReq := models.RefreshTokenRequest{
		RefreshToken: "valid-refresh-token",
	}

	expectedClaims := &models.JWTClaims{
		UserID:   "user-123",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
	}

	expectedUser := &models.User{
		ID:       "user-123",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
		IsActive: true,
	}

	expectedTokenPair := &models.TokenPair{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		ExpiresIn:    3600,
	}

	// Setup mocks
	suite.mockJWTSvc.On("ValidateToken", "valid-refresh-token").Return(expectedClaims, nil)
	suite.mockUserSvc.On("GetByID", "user-123").Return(expectedUser, nil)
	suite.mockJWTSvc.On("RefreshToken", "valid-refresh-token", expectedUser).Return(expectedTokenPair, nil)

	// Act
	w := suite.helper.MakeRequest("POST", "/auth/refresh", refreshReq, nil)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.AuthResponse
	err := suite.helper.ParseJSONResponse(w, &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedTokenPair.AccessToken, response.AccessToken)
}

// TestRefreshTokenInvalid tests refresh with invalid token
func (suite *AuthHandlerTestSuite) TestRefreshTokenInvalid() {
	// Arrange
	refreshReq := models.RefreshTokenRequest{
		RefreshToken: "invalid-refresh-token",
	}

	// Setup mocks
	suite.mockJWTSvc.On("ValidateToken", "invalid-refresh-token").Return((*models.JWTClaims)(nil), errors.New("invalid token"))

	// Act
	w := suite.helper.MakeRequest("POST", "/auth/refresh", refreshReq, nil)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	suite.helper.AssertErrorResponse(suite.T(), w, http.StatusUnauthorized, "invalid_token")
}

// TestAuthHandlerTestSuite runs the test suite
func TestAuthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}
