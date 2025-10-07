package tests

import (
	"net/http"
	"testing"

	"base-golang-restful-app/handlers"
	"base-golang-restful-app/models"
	"base-golang-restful-app/services"
	"base-golang-restful-app/tests/testhelpers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestSimpleHealthCheck tests a simple health check endpoint
func TestSimpleHealthCheck(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	helper := testhelpers.NewTestHelper()
	
	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	helper.Router = router

	// Test
	w := helper.MakeRequest("GET", "/health", nil, nil)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]string
	err := helper.ParseJSONResponse(w, &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
}

// TestUserServiceBasic tests basic user service functionality
func TestUserServiceBasic(t *testing.T) {
	// Setup
	userService := services.NewUserService()

	// Test user creation
	req := models.UserCreateRequest{
		Username:  "testuser",
		Email:     "test@example.com", 
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	user, err := userService.Create(req)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
}

// TestAuthHandlerBasic tests basic auth handler functionality
func TestAuthHandlerBasic(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	helper := testhelpers.NewTestHelper()
	
	userService := services.NewUserService()
	jwtService := helper.JWTService
	authHandler := handlers.NewAuthHandler(userService, jwtService)

	router := gin.New()
	router.POST("/register", authHandler.Register)
	helper.Router = router

	// Test registration
	registerReq := models.RegisterRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	w := helper.MakeRequest("POST", "/register", registerReq, nil)
	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.AuthResponse
	err := helper.ParseJSONResponse(w, &response)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", response.User.Username)
	assert.NotEmpty(t, response.AccessToken)
}
