package middleware

import (
	"net/http"
	"testing"
	"time"

	"base-golang-restful-app/middleware"
	"base-golang-restful-app/models"
	"base-golang-restful-app/tests/mocks"
	"base-golang-restful-app/tests/testhelpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// AuthMiddlewareTestSuite defines the test suite for AuthMiddleware
type AuthMiddlewareTestSuite struct {
	suite.Suite
	helper       *testhelpers.TestHelper
	mockUserSvc  *mocks.MockUserService
	mockJWTSvc   *mocks.MockJWTService
	router       *gin.Engine
	testUser     *models.User
	testAdmin    *models.User
	validToken   string
	invalidToken string
}

// SetupTest sets up the test environment before each test
func (suite *AuthMiddlewareTestSuite) SetupTest() {
	suite.helper = testhelpers.NewTestHelper()
	suite.mockUserSvc = new(mocks.MockUserService)
	suite.mockJWTSvc = new(mocks.MockJWTService)

	// Create test users
	suite.testUser = &models.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	suite.testAdmin = &models.User{
		ID:        uuid.New(),
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "User",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	suite.validToken = "valid-jwt-token"
	suite.invalidToken = "invalid-jwt-token"

	// Setup router with middleware
	suite.router = gin.New()

	// Protected route
	protected := suite.router.Group("/protected")
	protected.Use(middleware.AuthMiddleware(suite.mockJWTSvc, suite.mockUserSvc))
	{
		protected.GET("/user", func(c *gin.Context) {
			user, _ := c.Get("user")
			c.JSON(http.StatusOK, gin.H{"user": user})
		})
	}

	// Admin only route
	admin := suite.router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(suite.mockJWTSvc, suite.mockUserSvc))
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.GET("/users", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
		})
	}

	// Optional auth route
	optional := suite.router.Group("/optional")
	optional.Use(middleware.OptionalAuthMiddleware(suite.mockJWTSvc, suite.mockUserSvc))
	{
		optional.GET("/public", func(c *gin.Context) {
			user, exists := c.Get("user")
			if exists {
				c.JSON(http.StatusOK, gin.H{"authenticated": true, "user": user})
			} else {
				c.JSON(http.StatusOK, gin.H{"authenticated": false})
			}
		})
	}

	suite.helper.Router = suite.router
}

// TearDownTest cleans up after each test
func (suite *AuthMiddlewareTestSuite) TearDownTest() {
	suite.mockUserSvc.AssertExpectations(suite.T())
	suite.mockJWTSvc.AssertExpectations(suite.T())
}

// TestAuthMiddlewareSuccess tests successful authentication
func (suite *AuthMiddlewareTestSuite) TestAuthMiddlewareSuccess() {
	// Arrange
	expectedClaims := &models.JWTClaims{
		UserID:   suite.testUser.ID,
		Username: suite.testUser.Username,
		Email:    suite.testUser.Email,
		Role:     suite.testUser.Role,
	}

	suite.mockJWTSvc.On("ExtractTokenFromHeader", "Bearer "+suite.validToken).Return(suite.validToken, nil)
	suite.mockJWTSvc.On("ValidateToken", suite.validToken).Return(expectedClaims, nil)
	suite.mockUserSvc.On("GetByID", suite.testUser.ID).Return(suite.testUser, nil)

	// Act
	w := suite.helper.MakeAuthenticatedRequest("GET", "/protected/user", nil, suite.validToken)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response struct {
		User models.User `json:"user"`
	}
	err := suite.helper.ParseJSONResponse(w, &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.testUser.Username, response.User.Username)
}

// TestAuthMiddlewareNoToken tests authentication without token
func (suite *AuthMiddlewareTestSuite) TestAuthMiddlewareNoToken() {
	// Arrange
	suite.mockJWTSvc.On("ExtractTokenFromHeader", "").Return("", assert.AnError)

	// Act
	w := suite.helper.MakeRequest("GET", "/protected/user", nil, nil)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	suite.helper.AssertErrorResponse(suite.T(), w, http.StatusUnauthorized, "unauthorized")
}

// TestAuthMiddlewareInvalidToken tests authentication with invalid token
func (suite *AuthMiddlewareTestSuite) TestAuthMiddlewareInvalidToken() {
	// Arrange
	suite.mockJWTSvc.On("ExtractTokenFromHeader", "Bearer "+suite.invalidToken).Return(suite.invalidToken, nil)
	suite.mockJWTSvc.On("ValidateToken", suite.invalidToken).Return((*models.JWTClaims)(nil), assert.AnError)

	// Act
	w := suite.helper.MakeAuthenticatedRequest("GET", "/protected/user", nil, suite.invalidToken)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	suite.helper.AssertErrorResponse(suite.T(), w, http.StatusUnauthorized, "unauthorized")
}

// TestAuthMiddlewareUserNotFound tests authentication when user doesn't exist
func (suite *AuthMiddlewareTestSuite) TestAuthMiddlewareUserNotFound() {
	// Arrange
	expectedClaims := &models.JWTClaims{
		UserID:   "nonexistent-user",
		Username: "nonexistent",
		Email:    "nonexistent@example.com",
		Role:     "user",
	}

	suite.mockJWTSvc.On("ExtractTokenFromHeader", "Bearer "+suite.validToken).Return(suite.validToken, nil)
	suite.mockJWTSvc.On("ValidateToken", suite.validToken).Return(expectedClaims, nil)
	suite.mockUserSvc.On("GetByID", "nonexistent-user").Return((*models.User)(nil), assert.AnError)

	// Act
	w := suite.helper.MakeAuthenticatedRequest("GET", "/protected/user", nil, suite.validToken)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	suite.helper.AssertErrorResponse(suite.T(), w, http.StatusUnauthorized, "unauthorized")
}

// TestAuthMiddlewareInactiveUser tests authentication with inactive user
func (suite *AuthMiddlewareTestSuite) TestAuthMiddlewareInactiveUser() {
	// Arrange
	inactiveUser := *suite.testUser
	inactiveUser.IsActive = false

	expectedClaims := &models.JWTClaims{
		UserID:   suite.testUser.ID,
		Username: suite.testUser.Username,
		Email:    suite.testUser.Email,
		Role:     suite.testUser.Role,
	}

	suite.mockJWTSvc.On("ExtractTokenFromHeader", "Bearer "+suite.validToken).Return(suite.validToken, nil)
	suite.mockJWTSvc.On("ValidateToken", suite.validToken).Return(expectedClaims, nil)
	suite.mockUserSvc.On("GetByID", suite.testUser.ID).Return(&inactiveUser, nil)

	// Act
	w := suite.helper.MakeAuthenticatedRequest("GET", "/protected/user", nil, suite.validToken)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	suite.helper.AssertErrorResponse(suite.T(), w, http.StatusUnauthorized, "unauthorized")
}

// TestRequireRoleSuccess tests successful role-based authorization
func (suite *AuthMiddlewareTestSuite) TestRequireRoleSuccess() {
	// Arrange
	expectedClaims := &models.JWTClaims{
		UserID:   suite.testAdmin.ID,
		Username: suite.testAdmin.Username,
		Email:    suite.testAdmin.Email,
		Role:     suite.testAdmin.Role,
	}

	suite.mockJWTSvc.On("ExtractTokenFromHeader", "Bearer "+suite.validToken).Return(suite.validToken, nil)
	suite.mockJWTSvc.On("ValidateToken", suite.validToken).Return(expectedClaims, nil)
	suite.mockUserSvc.On("GetByID", suite.testAdmin.ID).Return(suite.testAdmin, nil)

	// Act
	w := suite.helper.MakeAuthenticatedRequest("GET", "/admin/users", nil, suite.validToken)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response struct {
		Message string `json:"message"`
	}
	err := suite.helper.ParseJSONResponse(w, &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "admin access granted", response.Message)
}

// TestRequireRoleForbidden tests role-based authorization failure
func (suite *AuthMiddlewareTestSuite) TestRequireRoleForbidden() {
	// Arrange
	expectedClaims := &models.JWTClaims{
		UserID:   suite.testUser.ID,
		Username: suite.testUser.Username,
		Email:    suite.testUser.Email,
		Role:     suite.testUser.Role,
	}

	suite.mockJWTSvc.On("ExtractTokenFromHeader", "Bearer "+suite.validToken).Return(suite.validToken, nil)
	suite.mockJWTSvc.On("ValidateToken", suite.validToken).Return(expectedClaims, nil)
	suite.mockUserSvc.On("GetByID", suite.testUser.ID).Return(suite.testUser, nil)

	// Act
	w := suite.helper.MakeAuthenticatedRequest("GET", "/admin/users", nil, suite.validToken)

	// Assert
	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
	suite.helper.AssertErrorResponse(suite.T(), w, http.StatusForbidden, "forbidden")
}

// TestOptionalAuthMiddlewareWithToken tests optional auth with valid token
func (suite *AuthMiddlewareTestSuite) TestOptionalAuthMiddlewareWithToken() {
	// Arrange
	expectedClaims := &models.JWTClaims{
		UserID:   suite.testUser.ID,
		Username: suite.testUser.Username,
		Email:    suite.testUser.Email,
		Role:     suite.testUser.Role,
	}

	suite.mockJWTSvc.On("ExtractTokenFromHeader", "Bearer "+suite.validToken).Return(suite.validToken, nil)
	suite.mockJWTSvc.On("ValidateToken", suite.validToken).Return(expectedClaims, nil)
	suite.mockUserSvc.On("GetByID", suite.testUser.ID).Return(suite.testUser, nil)

	// Act
	w := suite.helper.MakeAuthenticatedRequest("GET", "/optional/public", nil, suite.validToken)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response struct {
		Authenticated bool        `json:"authenticated"`
		User          models.User `json:"user"`
	}
	err := suite.helper.ParseJSONResponse(w, &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Authenticated)
	assert.Equal(suite.T(), suite.testUser.Username, response.User.Username)
}

// TestOptionalAuthMiddlewareWithoutToken tests optional auth without token
func (suite *AuthMiddlewareTestSuite) TestOptionalAuthMiddlewareWithoutToken() {
	// Act
	w := suite.helper.MakeRequest("GET", "/optional/public", nil, nil)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response struct {
		Authenticated bool `json:"authenticated"`
	}
	err := suite.helper.ParseJSONResponse(w, &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Authenticated)
}

// TestOptionalAuthMiddlewareWithInvalidToken tests optional auth with invalid token
func (suite *AuthMiddlewareTestSuite) TestOptionalAuthMiddlewareWithInvalidToken() {
	// Arrange
	suite.mockJWTSvc.On("ExtractTokenFromHeader", "Bearer "+suite.invalidToken).Return(suite.invalidToken, nil)
	suite.mockJWTSvc.On("ValidateToken", suite.invalidToken).Return((*models.JWTClaims)(nil), assert.AnError)

	// Act
	w := suite.helper.MakeAuthenticatedRequest("GET", "/optional/public", nil, suite.invalidToken)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response struct {
		Authenticated bool `json:"authenticated"`
	}
	err := suite.helper.ParseJSONResponse(w, &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Authenticated)
}

// TestAuthMiddlewareTestSuite runs the test suite
func TestAuthMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}
