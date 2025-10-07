package services

import (
	"testing"

	"base-golang-restful-app/models"
	"base-golang-restful-app/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserServiceTestSuite defines the test suite for UserService
type UserServiceTestSuite struct {
	suite.Suite
	userService *services.UserService
}

// SetupTest sets up the test environment before each test
func (suite *UserServiceTestSuite) SetupTest() {
	suite.userService = services.NewUserService()
}

// TestCreateUserSuccess tests successful user creation
func (suite *UserServiceTestSuite) TestCreateUserSuccess() {
	// Arrange
	req := models.UserCreateRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	// Act
	user, err := suite.userService.Create(req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), req.Username, user.Username)
	assert.Equal(suite.T(), req.Email, user.Email)
	assert.Equal(suite.T(), req.FirstName, user.FirstName)
	assert.Equal(suite.T(), req.LastName, user.LastName)
	assert.Equal(suite.T(), "user", user.Role)
	assert.True(suite.T(), user.IsActive)
	assert.NotEmpty(suite.T(), user.ID)
	assert.NotEqual(suite.T(), req.Password, user.Password) // Password should be hashed
}

// TestCreateUserDuplicateUsername tests user creation with duplicate username
func (suite *UserServiceTestSuite) TestCreateUserDuplicateUsername() {
	// Arrange
	req1 := models.UserCreateRequest{
		Username:  "testuser",
		Email:     "test1@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	req2 := models.UserCreateRequest{
		Username:  "testuser", // Same username
		Email:     "test2@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	// Act
	user1, err1 := suite.userService.Create(req1)
	user2, err2 := suite.userService.Create(req2)

	// Assert
	assert.NoError(suite.T(), err1)
	assert.NotNil(suite.T(), user1)
	assert.Error(suite.T(), err2)
	assert.Nil(suite.T(), user2)
	assert.Contains(suite.T(), err2.Error(), "username already exists")
}

// TestCreateUserDuplicateEmail tests user creation with duplicate email
func (suite *UserServiceTestSuite) TestCreateUserDuplicateEmail() {
	// Arrange
	req1 := models.UserCreateRequest{
		Username:  "testuser1",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	req2 := models.UserCreateRequest{
		Username:  "testuser2",
		Email:     "test@example.com", // Same email
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	// Act
	user1, err1 := suite.userService.Create(req1)
	user2, err2 := suite.userService.Create(req2)

	// Assert
	assert.NoError(suite.T(), err1)
	assert.NotNil(suite.T(), user1)
	assert.Error(suite.T(), err2)
	assert.Nil(suite.T(), user2)
	assert.Contains(suite.T(), err2.Error(), "email already exists")
}

// TestCreateUserWeakPassword tests user creation with weak password
func (suite *UserServiceTestSuite) TestCreateUserWeakPassword() {
	// Arrange
	req := models.UserCreateRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "123", // Weak password
		FirstName: "Test",
		LastName:  "User",
	}

	// Act
	user, err := suite.userService.Create(req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Contains(suite.T(), err.Error(), "password validation failed")
}

// TestGetUserByID tests user retrieval by ID
func (suite *UserServiceTestSuite) TestGetUserByID() {
	// Arrange
	req := models.UserCreateRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	createdUser, _ := suite.userService.Create(req)

	// Act
	foundUser, err := suite.userService.GetByID(createdUser.ID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), foundUser)
	assert.Equal(suite.T(), createdUser.ID, foundUser.ID)
	assert.Equal(suite.T(), createdUser.Username, foundUser.Username)
}

// TestGetUserByIDNotFound tests user retrieval with non-existent ID
func (suite *UserServiceTestSuite) TestGetUserByIDNotFound() {
	// Act
	user, err := suite.userService.GetByID("nonexistent-id")

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Contains(suite.T(), err.Error(), "user not found")
}

// TestGetUserByUsername tests user retrieval by username
func (suite *UserServiceTestSuite) TestGetUserByUsername() {
	// Arrange
	req := models.UserCreateRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	createdUser, _ := suite.userService.Create(req)

	// Act
	foundUser, err := suite.userService.GetByUsername(createdUser.Username)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), foundUser)
	assert.Equal(suite.T(), createdUser.Username, foundUser.Username)
}

// TestUpdateUserSuccess tests successful user update
func (suite *UserServiceTestSuite) TestUpdateUserSuccess() {
	// Arrange
	req := models.UserCreateRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	createdUser, _ := suite.userService.Create(req)

	updateReq := models.UserUpdateRequest{
		FirstName: stringPtr("Updated"),
		LastName:  stringPtr("Name"),
	}

	// Act
	updatedUser, err := suite.userService.Update(createdUser.ID, updateReq)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), updatedUser)
	assert.Equal(suite.T(), "Updated", updatedUser.FirstName)
	assert.Equal(suite.T(), "Name", updatedUser.LastName)
	assert.Equal(suite.T(), createdUser.Username, updatedUser.Username) // Unchanged
}

// TestDeleteUser tests user deletion (soft delete)
func (suite *UserServiceTestSuite) TestDeleteUser() {
	// Arrange
	req := models.UserCreateRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	createdUser, _ := suite.userService.Create(req)

	// Act
	err := suite.userService.Delete(createdUser.ID)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify user is soft deleted
	deletedUser, getErr := suite.userService.GetByID(createdUser.ID)
	assert.NoError(suite.T(), getErr)
	assert.False(suite.T(), deletedUser.IsActive)
}

// TestListUsers tests user listing with pagination
func (suite *UserServiceTestSuite) TestListUsers() {
	// Arrange
	for i := 0; i < 5; i++ {
		req := models.UserCreateRequest{
			Username:  "testuser" + string(rune(i+'0')),
			Email:     "test" + string(rune(i+'0')) + "@example.com",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
		}
		suite.userService.Create(req)
	}

	// Act
	users, pagination, err := suite.userService.List(1, 3)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 3)
	assert.Equal(suite.T(), 1, pagination.Page)
	assert.Equal(suite.T(), 3, pagination.Limit)
	assert.Equal(suite.T(), int64(5), pagination.Total)
	assert.Equal(suite.T(), 2, pagination.TotalPages)
	assert.True(suite.T(), pagination.HasNext)
	assert.False(suite.T(), pagination.HasPrev)
}

// TestValidateCredentialsSuccess tests successful credential validation
func (suite *UserServiceTestSuite) TestValidateCredentialsSuccess() {
	// Arrange
	req := models.UserCreateRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	createdUser, _ := suite.userService.Create(req)

	// Act
	validatedUser, err := suite.userService.ValidateCredentials("testuser", "password123")

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), validatedUser)
	assert.Equal(suite.T(), createdUser.ID, validatedUser.ID)
}

// TestValidateCredentialsInvalidPassword tests credential validation with wrong password
func (suite *UserServiceTestSuite) TestValidateCredentialsInvalidPassword() {
	// Arrange
	req := models.UserCreateRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	suite.userService.Create(req)

	// Act
	user, err := suite.userService.ValidateCredentials("testuser", "wrongpassword")

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Contains(suite.T(), err.Error(), "invalid credentials")
}

// TestChangePasswordSuccess tests successful password change
func (suite *UserServiceTestSuite) TestChangePasswordSuccess() {
	// Arrange
	req := models.UserCreateRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	createdUser, _ := suite.userService.Create(req)

	// Act
	err := suite.userService.ChangePassword(createdUser.ID, "password123", "newpassword123")

	// Assert
	assert.NoError(suite.T(), err)

	// Verify new password works
	validatedUser, validateErr := suite.userService.ValidateCredentials("testuser", "newpassword123")
	assert.NoError(suite.T(), validateErr)
	assert.NotNil(suite.T(), validatedUser)

	// Verify old password doesn't work
	_, oldPassErr := suite.userService.ValidateCredentials("testuser", "password123")
	assert.Error(suite.T(), oldPassErr)
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

// TestUserServiceTestSuite runs the test suite
func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
