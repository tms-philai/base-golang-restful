package handlers

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"base-golang-restful-app/handlers"
	"base-golang-restful-app/middleware"
	"base-golang-restful-app/models"
	"base-golang-restful-app/tests/mocks"
	"base-golang-restful-app/tests/testhelpers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// UserHandlerTestSuite defines the test suite for UserHandler
type UserHandlerTestSuite struct {
	suite.Suite
	helper         *testhelpers.TestHelper
	mockUserSvc    *mocks.MockUserService
	mockJWTSvc     *mocks.MockJWTService
	userHandler    *handlers.UserHandler
	router         *gin.Engine
	testUser       *models.User
	testAdmin      *models.User
	testToken      string
	adminToken     string
}

// SetupTest sets up the test environment before each test
func (suite *UserHandlerTestSuite) SetupTest() {
	suite.helper = testhelpers.NewTestHelper()
	suite.mockUserSvc = new(mocks.MockUserService)
	suite.mockJWTSvc = new(mocks.MockJWTService)
	suite.userHandler = handlers.NewUserHandler(suite.mockUserSvc)

	// Create test users
	suite.testUser = &models.User{
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

	suite.testAdmin = &models.User{
		ID:        "admin-123",
		Username:  "admin",
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "User",
		Role:      "admin",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	suite.testToken = "valid-user-token"
	suite.adminToken = "valid-admin-token"

	// Setup router with middleware
	suite.router = gin.New()
	
	// Mock middleware that sets user context
	authMiddleware := func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "Bearer "+suite.testToken {
			c.Set("user_id", suite.testUser.ID)
			c.Set("user", suite.testUser)
		} else if authHeader == "Bearer "+suite.adminToken {
			c.Set("user_id", suite.testAdmin.ID)
			c.Set("user", suite.testAdmin)
		}
		c.Next()
	}

	adminOnlyMiddleware := func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		userModel := user.(*models.User)
		if userModel.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}

	users := suite.router.Group("/users")
	{
		users.GET("/:id", suite.userHandler.GetUser)
		users.GET("", authMiddleware, adminOnlyMiddleware, suite.userHandler.ListUsers)
		users.POST("", authMiddleware, adminOnlyMiddleware, suite.userHandler.CreateUser)
		users.PUT("/:id", authMiddleware, suite.userHandler.UpdateUser)
		users.DELETE("/:id", authMiddleware, adminOnlyMiddleware, suite.userHandler.DeleteUser)
	}

	suite.helper.Router = suite.router
}

// TearDownTest cleans up after each test
func (suite *UserHandlerTestSuite) TearDownTest() {
	suite.mockUserSvc.AssertExpectations(suite.T())
}

// TestGetUserSuccess tests successful user retrieval
func (suite *UserHandlerTestSuite) TestGetUserSuccess() {
	// Arrange
	userID := "user-123"
	suite.mockUserSvc.On("GetByID", userID).Return(suite.testUser, nil)

	// Act
	w := suite.helper.MakeRequest("GET", "/users/"+userID, nil, nil)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.UserResponse
	err := suite.helper.ParseJSONResponse(w, &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.testUser.Username, response.Username)
	assert.Equal(suite.T(), suite.testUser.Email, response.Email)
}

// TestGetUserNotFound tests user retrieval when user doesn't exist
func (suite *UserHandlerTestSuite) TestGetUserNotFound() {
	// Arrange
	userID := "nonexistent-user"
	suite.mockUserSvc.On("GetByID", userID).Return((*models.User)(nil), errors.New("user not found"))

	// Act
	w := suite.helper.MakeRequest("GET", "/users/"+userID, nil, nil)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	suite.helper.AssertErrorResponse(suite.T(), w, http.StatusNotFound, "user_not_found")
}

// TestCreateUserSuccess tests successful user creation by admin
func (suite *UserHandlerTestSuite) TestCreateUserSuccess() {
	// Arrange
	createReq := models.UserCreateRequest{
		Username:  "newuser",
		Email:     "newuser@example.com",
		Password:  "password123",
		FirstName: "New",
		LastName:  "User",
	}

	expectedUser := &models.User{
		ID:        "new-user-123",
		Username:  "newuser",
		Email:     "newuser@example.com",
		FirstName: "New",
		LastName:  "User",
		Role:      "user",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	suite.mockUserSvc.On("Create", mock.AnythingOfType("models.UserCreateRequest")).Return(expectedUser, nil)

	// Act
	w := suite.helper.MakeAuthenticatedRequest("POST", "/users", createReq, suite.adminToken)

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response models.UserResponse
	err := suite.helper.ParseJSONResponse(w, &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedUser.Username, response.Username)
}

// TestCreateUserForbidden tests user creation by non-admin user
func (suite *UserHandlerTestSuite) TestCreateUserForbidden() {
	// Arrange
	createReq := models.UserCreateRequest{
		Username:  "newuser",
		Email:     "newuser@example.com",
		Password:  "password123",
		FirstName: "New",
		LastName:  "User",
	}

	// Act
	w := suite.helper.MakeAuthenticatedRequest("POST", "/users", createReq, suite.testToken)

	// Assert
	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
}

// TestUpdateUserSuccess tests successful user update
func (suite *UserHandlerTestSuite) TestUpdateUserSuccess() {
	// Arrange
	userID := suite.testUser.ID
	updateReq := models.UserUpdateRequest{
		FirstName: stringPtr("Updated"),
		LastName:  stringPtr("Name"),
	}

	updatedUser := *suite.testUser
	updatedUser.FirstName = "Updated"
	updatedUser.LastName = "Name"

	suite.mockUserSvc.On("Update", userID, mock.AnythingOfType("models.UserUpdateRequest")).Return(&updatedUser, nil)

	// Act
	w := suite.helper.MakeAuthenticatedRequest("PUT", "/users/"+userID, updateReq, suite.testToken)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.UserResponse
	err := suite.helper.ParseJSONResponse(w, &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated", response.FirstName)
	assert.Equal(suite.T(), "Name", response.LastName)
}

// TestUpdateUserForbidden tests updating another user's profile
func (suite *UserHandlerTestSuite) TestUpdateUserForbidden() {
	// Arrange
	otherUserID := "other-user-123"
	updateReq := models.UserUpdateRequest{
		FirstName: stringPtr("Updated"),
	}

	// Act
	w := suite.helper.MakeAuthenticatedRequest("PUT", "/users/"+otherUserID, updateReq, suite.testToken)

	// Assert
	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
	suite.helper.AssertErrorResponse(suite.T(), w, http.StatusForbidden, "forbidden")
}

// TestDeleteUserSuccess tests successful user deletion by admin
func (suite *UserHandlerTestSuite) TestDeleteUserSuccess() {
	// Arrange
	userID := "user-to-delete"
	suite.mockUserSvc.On("Delete", userID).Return(nil)

	// Act
	w := suite.helper.MakeAuthenticatedRequest("DELETE", "/users/"+userID, nil, suite.adminToken)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.helper.AssertSuccessResponse(suite.T(), w, http.StatusOK, "User deleted successfully")
}

// TestDeleteUserForbidden tests user deletion by non-admin
func (suite *UserHandlerTestSuite) TestDeleteUserForbidden() {
	// Arrange
	userID := "user-to-delete"

	// Act
	w := suite.helper.MakeAuthenticatedRequest("DELETE", "/users/"+userID, nil, suite.testToken)

	// Assert
	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
}

// TestListUsersSuccess tests successful user listing by admin
func (suite *UserHandlerTestSuite) TestListUsersSuccess() {
	// Arrange
	users := []models.User{*suite.testUser, *suite.testAdmin}
	pagination := models.PaginationMetadata{
		Page:       1,
		Limit:      10,
		Total:      2,
		TotalPages: 1,
		HasNext:    false,
		HasPrev:    false,
	}

	suite.mockUserSvc.On("List", 1, 10).Return(users, pagination, nil)

	// Act
	w := suite.helper.MakeAuthenticatedRequest("GET", "/users?page=1&limit=10", nil, suite.adminToken)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response struct {
		Users      []models.UserResponse         `json:"users"`
		Pagination models.PaginationMetadata `json:"pagination"`
	}
	err := suite.helper.ParseJSONResponse(w, &response)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), response.Users, 2)
	assert.Equal(suite.T(), pagination.Total, response.Pagination.Total)
}

// TestListUsersForbidden tests user listing by non-admin
func (suite *UserHandlerTestSuite) TestListUsersForbidden() {
	// Act
	w := suite.helper.MakeAuthenticatedRequest("GET", "/users", nil, suite.testToken)

	// Assert
	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

// TestUserHandlerTestSuite runs the test suite
func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
