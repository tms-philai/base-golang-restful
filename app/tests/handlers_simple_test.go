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

// TestAuthHandlerRegister tests user registration endpoint
func TestAuthHandlerRegister(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	helper := testhelpers.NewTestHelper()
	
	userService := services.NewUserService()
	jwtService := helper.JWTService
	authHandler := handlers.NewAuthHandler(userService, jwtService)

	router := gin.New()
	router.POST("/register", authHandler.Register)
	helper.Router = router

	// Test successful registration
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
	assert.Equal(t, "Bearer", response.TokenType)
}

// TestAuthHandlerLogin tests user login endpoint
func TestAuthHandlerLogin(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	helper := testhelpers.NewTestHelper()
	
	userService := services.NewUserService()
	jwtService := helper.JWTService
	authHandler := handlers.NewAuthHandler(userService, jwtService)

	router := gin.New()
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	helper.Router = router

	// First register a user
	registerReq := models.RegisterRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	helper.MakeRequest("POST", "/register", registerReq, nil)

	// Test login
	loginReq := models.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	w := helper.MakeRequest("POST", "/login", loginReq, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.AuthResponse
	err := helper.ParseJSONResponse(w, &response)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", response.User.Username)
	assert.NotEmpty(t, response.AccessToken)
}

// TestAuthHandlerLoginInvalidCredentials tests login with wrong password
func TestAuthHandlerLoginInvalidCredentials(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	helper := testhelpers.NewTestHelper()
	
	userService := services.NewUserService()
	jwtService := helper.JWTService
	authHandler := handlers.NewAuthHandler(userService, jwtService)

	router := gin.New()
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	helper.Router = router

	// First register a user
	registerReq := models.RegisterRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	helper.MakeRequest("POST", "/register", registerReq, nil)

	// Test login with wrong password
	loginReq := models.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	w := helper.MakeRequest("POST", "/login", loginReq, nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	helper.AssertErrorResponse(t, w, http.StatusUnauthorized, "authentication_failed")
}

// TestUserHandlerBasic tests basic user handler functionality
func TestUserHandlerBasic(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	helper := testhelpers.NewTestHelper()
	
	userService := services.NewUserService()
	userHandler := handlers.NewUserHandler(userService)

	router := gin.New()
	router.GET("/users/:id", userHandler.GetUser)
	helper.Router = router

	// Create a test user first
	req := models.UserCreateRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	user, err := userService.Create(req)
	assert.NoError(t, err)

	// Test getting user by ID
	w := helper.MakeRequest("GET", "/users/"+user.ID, nil, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.UserResponse
	err = helper.ParseJSONResponse(w, &response)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", response.Username)
	assert.Equal(t, "test@example.com", response.Email)
}

// TestUserHandlerNotFound tests getting non-existent user
func TestUserHandlerNotFound(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	helper := testhelpers.NewTestHelper()
	
	userService := services.NewUserService()
	userHandler := handlers.NewUserHandler(userService)

	router := gin.New()
	router.GET("/users/:id", userHandler.GetUser)
	helper.Router = router

	// Test getting non-existent user
	w := helper.MakeRequest("GET", "/users/nonexistent-id", nil, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)

	helper.AssertErrorResponse(t, w, http.StatusNotFound, "user_not_found")
}

// TestProductHandlerBasic tests basic product handler functionality
func TestProductHandlerBasic(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	helper := testhelpers.NewTestHelper()
	
	productService := services.NewProductService()
	productHandler := handlers.NewProductHandler(productService)

	router := gin.New()
	router.GET("/products/categories", productHandler.GetCategories)
	router.GET("/products", productHandler.ListProducts)
	helper.Router = router

	// Test getting categories
	w := helper.MakeRequest("GET", "/products/categories", nil, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var categoriesResponse struct {
		Categories []string `json:"categories"`
	}
	err := helper.ParseJSONResponse(w, &categoriesResponse)
	assert.NoError(t, err)
	assert.IsType(t, []string{}, categoriesResponse.Categories)

	// Test listing products (should be empty initially)
	w2 := helper.MakeRequest("GET", "/products", nil, nil)
	assert.Equal(t, http.StatusOK, w2.Code)

	var listResponse models.ProductListResponse
	err = helper.ParseJSONResponse(w2, &listResponse)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), listResponse.Pagination.Total)
}
