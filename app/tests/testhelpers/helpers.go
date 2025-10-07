package testhelpers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"base-golang-restful-app/config"
	"base-golang-restful-app/models"
	"base-golang-restful-app/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestHelper provides common testing utilities
type TestHelper struct {
	UserService    *services.UserService
	ProductService *services.ProductService
	JWTService     *services.JWTService
	Router         *gin.Engine
}

// NewTestHelper creates a new test helper instance
func NewTestHelper() *TestHelper {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create test configuration
	cfg := &config.JWTConfig{
		SecretKey:            "test-secret-key",
		AccessTokenDuration:  3600,
		RefreshTokenDuration: 86400,
		Issuer:               "test-issuer",
	}

	// Initialize services
	userService := services.NewUserService()
	productService := services.NewProductService()
	jwtService := services.NewJWTService(cfg)

	return &TestHelper{
		UserService:    userService,
		ProductService: productService,
		JWTService:     jwtService,
		Router:         gin.New(),
	}
}

// CreateTestUser creates a test user and returns it
func (h *TestHelper) CreateTestUser(username, email string) (*models.User, error) {
	req := models.UserCreateRequest{
		Username:  username,
		Email:     email,
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	return h.UserService.Create(req)
}

// CreateTestAdmin creates a test admin user and returns it
func (h *TestHelper) CreateTestAdmin(username, email string) (*models.User, error) {
	user, err := h.CreateTestUser(username, email)
	if err != nil {
		return nil, err
	}
	user.Role = "admin"
	return user, nil
}

// CreateTestProduct creates a test product and returns it
func (h *TestHelper) CreateTestProduct(name, sku, createdBy string) (*models.Product, error) {
	req := models.ProductCreateRequest{
		Name:        name,
		Description: "Test product description",
		Price:       99.99,
		Category:    "Test Category",
		SKU:         sku,
		Stock:       100,
	}
	return h.ProductService.Create(req, createdBy)
}

// GenerateTestToken generates a JWT token for testing
func (h *TestHelper) GenerateTestToken(user *models.User) (string, error) {
	tokenPair, err := h.JWTService.GenerateTokenPair(user)
	if err != nil {
		return "", err
	}
	return tokenPair.AccessToken, nil
}

// MakeRequest makes an HTTP request for testing
func (h *TestHelper) MakeRequest(method, url string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody bytes.Buffer
	if body != nil {
		json.NewEncoder(&reqBody).Encode(body)
	}

	req, _ := http.NewRequest(method, url, &reqBody)
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	h.Router.ServeHTTP(w, req)
	return w
}

// MakeAuthenticatedRequest makes an authenticated HTTP request
func (h *TestHelper) MakeAuthenticatedRequest(method, url string, body interface{}, token string) *httptest.ResponseRecorder {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}
	return h.MakeRequest(method, url, body, headers)
}

// AssertJSONResponse asserts that the response contains expected JSON
func (h *TestHelper) AssertJSONResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedBody interface{}) {
	assert.Equal(t, expectedStatus, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	if expectedBody != nil {
		var actualBody interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		assert.NoError(t, err)

		expectedJSON, _ := json.Marshal(expectedBody)
		var expectedBodyParsed interface{}
		json.Unmarshal(expectedJSON, &expectedBodyParsed)

		assert.Equal(t, expectedBodyParsed, actualBody)
	}
}

// AssertErrorResponse asserts that the response is an error with expected message
func (h *TestHelper) AssertErrorResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedError string) {
	assert.Equal(t, expectedStatus, w.Code)

	var errorResp models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, expectedError, errorResp.Error)
}

// AssertSuccessResponse asserts that the response is a success with expected message
func (h *TestHelper) AssertSuccessResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedMessage string) {
	assert.Equal(t, expectedStatus, w.Code)

	var successResp models.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &successResp)
	assert.NoError(t, err)
	assert.Equal(t, expectedMessage, successResp.Message)
}

// ParseJSONResponse parses JSON response into provided struct
func (h *TestHelper) ParseJSONResponse(w *httptest.ResponseRecorder, target interface{}) error {
	return json.Unmarshal(w.Body.Bytes(), target)
}

// CleanupTestData cleans up test data (for in-memory storage)
func (h *TestHelper) CleanupTestData() {
	// Reset services to clean state
	h.UserService = services.NewUserService()
	h.ProductService = services.NewProductService()
}
