package helpers

import (
	"base-gin/internal/domain/models"
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestUser represents a test user for testing purposes
type TestUser struct {
	ID        uuid.UUID
	Username  string
	Email     string
	Password  string
	FirstName string
	LastName  string
	IsActive  bool
	Roles     []models.Role
}

// TestProduct represents a test product for testing purposes
type TestProduct struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       float64
	SKU         string
	Category    string
	Stock       int
	CreatedBy   string
}

// CreateTestUser creates a test user with default values
func CreateTestUser() *TestUser {
	return &TestUser{
		ID:        uuid.New(),
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
		Roles: []models.Role{
			{Name: "user"},
		},
	}
}

// CreateTestAdmin creates a test admin user
func CreateTestAdmin() *TestUser {
	return &TestUser{
		ID:        uuid.New(),
		Username:  "admin",
		Email:     "admin@example.com",
		Password:  "admin123",
		FirstName: "Admin",
		LastName:  "User",
		IsActive:  true,
		Roles: []models.Role{
			{Name: "admin"},
		},
	}
}

// CreateTestProduct creates a test product with default values
func CreateTestProduct() *TestProduct {
	return &TestProduct{
		ID:          uuid.New(),
		Name:        "Test Product",
		Description: "A test product for testing",
		Price:       99.99,
		SKU:         "TEST-001",
		Category:    "Electronics",
		Stock:       100,
		CreatedBy:   uuid.New().String(),
	}
}

// ToUser converts TestUser to models.User
func (tu *TestUser) ToUser() *models.User {
	return &models.User{
		ID:        tu.ID,
		Email:     tu.Email,
		FirstName: tu.FirstName,
		LastName:  tu.LastName,
		IsActive:  tu.IsActive,
		Roles:     tu.Roles,
	}
}

// ToProduct converts TestProduct to models.Product
func (tp *TestProduct) ToProduct() *models.Product {
	return &models.Product{
		ID:          tp.ID.String(),
		Name:        tp.Name,
		Description: tp.Description,
		Price:       tp.Price,
		SKU:         tp.SKU,
		Category:    tp.Category,
		Stock:       tp.Stock,
		CreatedBy:   tp.CreatedBy,
	}
}

// CreateJSONRequest creates a JSON HTTP request for testing
func CreateJSONRequest(method, url string, body interface{}) (*http.Request, error) {
	var jsonBody []byte
	var err error

	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// CreateMultipartRequest creates a multipart form HTTP request for testing
func CreateMultipartRequest(method, url string, formData map[string]string, fileField, fileName, fileContent string) (*http.Request, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add form fields
	for key, value := range formData {
		writer.WriteField(key, value)
	}

	// Add file if provided
	if fileField != "" && fileName != "" && fileContent != "" {
		fileWriter, err := writer.CreateFormFile(fileField, fileName)
		if err != nil {
			return nil, err
		}
		fileWriter.Write([]byte(fileContent))
	}

	writer.Close()

	req, err := http.NewRequest(method, url, &body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

// AssertJSONResponse asserts that the response is valid JSON and matches expected structure
func AssertJSONResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) {
	assert.Equal(t, expectedStatus, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

// AssertErrorResponse asserts that the response is an error response with expected error code
func AssertErrorResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedError string) {
	AssertJSONResponse(t, w, expectedStatus)

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedError, response.Error)
}

// AssertSuccessResponse asserts that the response is a success response
func AssertSuccessResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) {
	AssertJSONResponse(t, w, expectedStatus)

	var response models.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Message)
}

// AssertAuthResponse asserts that the response is a valid auth response
func AssertAuthResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) {
	AssertJSONResponse(t, w, expectedStatus)

	var response models.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
	assert.Equal(t, "Bearer", response.TokenType)
	assert.NotNil(t, response.User)
}

// AssertUserResponse asserts that the response is a valid user response
func AssertUserResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) {
	AssertJSONResponse(t, w, expectedStatus)

	var response models.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.ID)
	assert.NotEmpty(t, response.Email)
}

// AssertProductResponse asserts that the response is a valid product response
func AssertProductResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) {
	AssertJSONResponse(t, w, expectedStatus)

	var response models.Product
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.ID)
	assert.NotEmpty(t, response.Name)
}

// AssertFileResponse asserts that the response is a valid file response
func AssertFileResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) {
	AssertJSONResponse(t, w, expectedStatus)

	var response models.File
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.ID)
	assert.NotEmpty(t, response.FileName)
}

// AssertPaginationResponse asserts that the response contains pagination metadata
func AssertPaginationResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) {
	AssertJSONResponse(t, w, expectedStatus)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "pagination")

	pagination, ok := response["pagination"].(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, pagination, "page")
	assert.Contains(t, pagination, "limit")
	assert.Contains(t, pagination, "total")
	assert.Contains(t, pagination, "total_pages")
}

// MockAuthMiddleware creates a mock authentication middleware for testing
func MockAuthMiddleware(user *models.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		if user != nil {
			c.Set("current_user", user)
			c.Set("user_id", user.ID.String())
		}
		c.Next()
	}
}

// MockOptionalAuthMiddleware creates a mock optional authentication middleware for testing
func MockOptionalAuthMiddleware(user *models.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		if user != nil {
			c.Set("current_user", user)
			c.Set("user_id", user.ID.String())
		}
		c.Next()
	}
}

// MockAdminMiddleware creates a mock admin middleware for testing
func MockAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mock admin role check
		c.Next()
	}
}

// SetupTestRouter creates a test router with mocked middleware
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// CleanupTestData cleans up test data (placeholder for database cleanup)
func CleanupTestData() {
	// This would contain cleanup logic for test data
	// For now, it's a placeholder
}

// GenerateTestJWT generates a test JWT token (placeholder)
func GenerateTestJWT(userID string) string {
	// This would generate a real JWT token for testing
	// For now, it's a placeholder
	return "test-jwt-token"
}

func StringPtr(s string) *string {
	return &s
}

func Float64Ptr(f float64) *float64 {
	return &f
}
