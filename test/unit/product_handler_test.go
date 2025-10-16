package unit

import (
	"base-gin/internal/app/handlers"
	"base-gin/internal/domain/models"
	"base-gin/test/helpers"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductService is a mock implementation of ProductService
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) Create(req models.ProductCreateRequest, userID string) (*models.Product, error) {
	args := m.Called(req, userID)
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) GetByID(id string) (*models.Product, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) Update(id string, req models.ProductUpdateRequest) (*models.Product, error) {
	args := m.Called(id, req)
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductService) List(query models.ProductListQuery) ([]models.Product, models.PaginationMetadata, error) {
	args := m.Called(query)
	return args.Get(0).([]models.Product), args.Get(1).(models.PaginationMetadata), args.Error(2)
}

func (m *MockProductService) GetCategories() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockProductService) UpdateStock(id string, stock int) error {
	args := m.Called(id, stock)
	return args.Error(0)
}

func TestProductHandler_CreateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		request        models.ProductCreateRequest
		mockSetup      func(*MockProductService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successful product creation",
			userID: "user-123",
			request: models.ProductCreateRequest{
				Name:        "Test Product",
				Description: "A test product",
				Price:       99.99,
				SKU:         "TEST-001",
				Category:    "Electronics",
				Stock:       100,
			},
			mockSetup: func(productService *MockProductService) {
				product := &models.Product{
					ID:          uuid.New().String(),
					Name:        "Test Product",
					Description: "A test product",
					Price:       99.99,
					SKU:         "TEST-001",
					Category:    "Electronics",
					Stock:       100,
					CreatedBy:   "user-123",
				}
				productService.On("Create", mock.AnythingOfType("models.ProductCreateRequest"), "user-123").Return(product, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "validation error - missing name",
			userID: "user-123",
			request: models.ProductCreateRequest{
				Description: "A test product",
				Price:       99.99,
				SKU:         "TEST-001",
				Category:    "Electronics",
				Stock:       100,
			},
			mockSetup:      func(*MockProductService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name:   "product creation failed - SKU already exists",
			userID: "user-123",
			request: models.ProductCreateRequest{
				Name:        "Test Product",
				Description: "A test product",
				Price:       99.99,
				SKU:         "EXISTING-SKU",
				Category:    "Electronics",
				Stock:       100,
			},
			mockSetup: func(productService *MockProductService) {
				productService.On("Create", mock.AnythingOfType("models.ProductCreateRequest"), "user-123").Return((*models.Product)(nil), errors.New("SKU already exists"))
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "product_creation_failed",
		},
		{
			name:           "unauthorized - missing user ID",
			userID:         "",
			request:        models.ProductCreateRequest{},
			mockSetup:      func(*MockProductService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockProductService := new(MockProductService)
			tt.mockSetup(mockProductService)

			// Create handler
			handler := handlers.NewProductHandler(mockProductService)

			// Setup router with middleware mock
			router := gin.New()
			router.POST("/products", func(c *gin.Context) {
				// Mock middleware setting user ID
				if tt.userID != "" {
					c.Set("user_id", tt.userID)
				}
				handler.CreateProduct(c)
			})

			// Create request
			jsonBody, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response.Error)
			}

			// Verify mock expectations
			mockProductService.AssertExpectations(t)
		})
	}
}

func TestProductHandler_GetProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		productID      string
		mockSetup      func(*MockProductService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:      "successful get product",
			productID: "product-123",
			mockSetup: func(productService *MockProductService) {
				product := &models.Product{
					ID:          uuid.New().String(),
					Name:        "Test Product",
					Description: "A test product",
					Price:       99.99,
					SKU:         "TEST-001",
					Category:    "Electronics",
					Stock:       100,
				}
				productService.On("GetByID", "product-123").Return(product, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "product not found",
			productID: "nonexistent",
			mockSetup: func(productService *MockProductService) {
				productService.On("GetByID", "nonexistent").Return((*models.Product)(nil), errors.New("product not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "product_not_found",
		},
		{
			name:           "missing product ID",
			productID:      "",
			mockSetup:      func(*MockProductService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockProductService := new(MockProductService)
			tt.mockSetup(mockProductService)

			// Create handler
			handler := handlers.NewProductHandler(mockProductService)

			// Setup router
			router := gin.New()
			router.GET("/products/:id", handler.GetProduct)

			// Create request
			url := "/products/" + tt.productID
			req, _ := http.NewRequest("GET", url, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response.Error)
			}

			// Verify mock expectations
			mockProductService.AssertExpectations(t)
		})
	}
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		productID      string
		currentUserID  string
		request        models.ProductUpdateRequest
		mockSetup      func(*MockProductService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:          "successful update by owner",
			productID:     "product-123",
			currentUserID: "user-123",
			request: models.ProductUpdateRequest{
				Name:        helpers.StringPtr("Updated Product"),
				Description: helpers.StringPtr("Updated description"),
				Price:       helpers.Float64Ptr(149.99),
			},
			mockSetup: func(productService *MockProductService) {
				existingProduct := &models.Product{
					ID:        uuid.New().String(),
					Name:      "Test Product",
					CreatedBy: "user-123",
				}
				updatedProduct := &models.Product{
					ID:          uuid.New().String(),
					Name:        *helpers.StringPtr("Updated Product"),
					Description: *helpers.StringPtr("Updated description"),
					Price:       *helpers.Float64Ptr(149.99),
					CreatedBy:   "user-123",
				}

				productService.On("GetByID", "product-123").Return(existingProduct, nil)
				productService.On("Update", "product-123", mock.AnythingOfType("models.ProductUpdateRequest")).Return(updatedProduct, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "successful update by admin",
			productID:     "product-123",
			currentUserID: "admin-456",
			request: models.ProductUpdateRequest{
				Name:        helpers.StringPtr("Updated Product"),
				Description: helpers.StringPtr("Updated description"),
				Price:       helpers.Float64Ptr(149.99),
			},
			mockSetup: func(productService *MockProductService) {
				existingProduct := &models.Product{
					ID:        uuid.New().String(),
					Name:      "Test Product",
					CreatedBy: "user-123",
				}
				updatedProduct := &models.Product{
					ID:          uuid.New().String(),
					Name:        *helpers.StringPtr("Updated Product"),
					Description: *helpers.StringPtr("Updated description"),
					Price:       *helpers.Float64Ptr(149.99),
					CreatedBy:   "user-123",
				}

				productService.On("GetByID", "product-123").Return(existingProduct, nil)
				productService.On("Update", "product-123", mock.AnythingOfType("models.ProductUpdateRequest")).Return(updatedProduct, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "forbidden - user trying to update another user's product",
			productID:     "product-123",
			currentUserID: "user-456",
			request: models.ProductUpdateRequest{
				Name: helpers.StringPtr("Updated Product"),
			},
			mockSetup: func(productService *MockProductService) {
				existingProduct := &models.Product{
					ID:        uuid.New().String(),
					Name:      "Test Product",
					CreatedBy: "user-123",
				}
				productService.On("GetByID", "product-123").Return(existingProduct, nil)
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  "forbidden",
		},
		{
			name:          "product not found",
			productID:     "nonexistent",
			currentUserID: "user-123",
			request: models.ProductUpdateRequest{
				Name: helpers.StringPtr("Updated Product"),
			},
			mockSetup: func(productService *MockProductService) {
				productService.On("GetByID", "nonexistent").Return((*models.Product)(nil), errors.New("product not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "product_not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockProductService := new(MockProductService)
			tt.mockSetup(mockProductService)

			// Create handler
			handler := handlers.NewProductHandler(mockProductService)

			// Setup router with middleware mock
			router := gin.New()
			router.PUT("/products/:id", func(c *gin.Context) {
				// Mock middleware setting current user
				user := &models.User{
					ID: uuid.MustParse(tt.currentUserID),
				}
				// Add admin role if current user is admin
				if tt.currentUserID == "admin-456" {
					user.Roles = []models.Role{{Name: "admin"}}
				}
				c.Set("current_user", user)
				handler.UpdateProduct(c)
			})

			// Create request
			jsonBody, _ := json.Marshal(tt.request)
			url := "/products/" + tt.productID
			req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response.Error)
			}

			// Verify mock expectations
			mockProductService.AssertExpectations(t)
		})
	}
}

func TestProductHandler_ListProducts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*MockProductService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:        "successful list products",
			queryParams: "?page=1&limit=10",
			mockSetup: func(productService *MockProductService) {
				products := []models.Product{
					{
						ID:       uuid.New().String(),
						Name:     "Product 1",
						Price:    99.99,
						SKU:      "PROD-001",
						Category: "Electronics",
						Stock:    100,
					},
					{
						ID:       uuid.New().String(),
						Name:     "Product 2",
						Price:    149.99,
						SKU:      "PROD-002",
						Category: "Clothing",
						Stock:    50,
					},
				}
				pagination := models.PaginationMetadata{
					Page:       1,
					Limit:      10,
					Total:      2,
					TotalPages: 1,
				}
				productService.On("List", mock.AnythingOfType("models.ProductListQuery")).Return(products, pagination, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "list products with search and category filter",
			queryParams: "?page=1&limit=10&search=test&category=Electronics",
			mockSetup: func(productService *MockProductService) {
				products := []models.Product{}
				pagination := models.PaginationMetadata{
					Page:       1,
					Limit:      10,
					Total:      0,
					TotalPages: 0,
				}
				productService.On("List", mock.AnythingOfType("models.ProductListQuery")).Return(products, pagination, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "list products service error",
			queryParams: "?page=1&limit=10",
			mockSetup: func(productService *MockProductService) {
				productService.On("List", mock.AnythingOfType("models.ProductListQuery")).Return(([]models.Product)(nil), models.PaginationMetadata{}, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "list_products_failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockProductService := new(MockProductService)
			tt.mockSetup(mockProductService)

			// Create handler
			handler := handlers.NewProductHandler(mockProductService)

			// Setup router
			router := gin.New()
			router.GET("/products", handler.ListProducts)

			// Create request
			url := "/products" + tt.queryParams
			req, _ := http.NewRequest("GET", url, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response.Error)
			}

			// Verify mock expectations
			mockProductService.AssertExpectations(t)
		})
	}
}

func TestProductHandler_GetCategories(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockSetup      func(*MockProductService)
		expectedStatus int
	}{
		{
			name: "successful get categories",
			mockSetup: func(productService *MockProductService) {
				categories := []string{"Electronics", "Clothing", "Books", "Home & Garden"}
				productService.On("GetCategories").Return(categories)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockProductService := new(MockProductService)
			tt.mockSetup(mockProductService)

			// Create handler
			handler := handlers.NewProductHandler(mockProductService)

			// Setup router
			router := gin.New()
			router.GET("/products/categories", handler.GetCategories)

			// Create request
			req, _ := http.NewRequest("GET", "/products/categories", nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verify mock expectations
			mockProductService.AssertExpectations(t)
		})
	}
}

func TestProductHandler_UpdateStock(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		productID      string
		currentUserID  string
		request        map[string]interface{}
		mockSetup      func(*MockProductService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:          "successful stock update by owner",
			productID:     "product-123",
			currentUserID: "user-123",
			request: map[string]interface{}{
				"stock": 150,
			},
			mockSetup: func(productService *MockProductService) {
				existingProduct := &models.Product{
					ID:        uuid.New().String(),
					Name:      "Test Product",
					CreatedBy: "user-123",
				}
				productService.On("GetByID", "product-123").Return(existingProduct, nil)
				productService.On("UpdateStock", "product-123", 150).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "forbidden - user trying to update another user's product stock",
			productID:     "product-123",
			currentUserID: "user-456",
			request: map[string]interface{}{
				"stock": 150,
			},
			mockSetup: func(productService *MockProductService) {
				existingProduct := &models.Product{
					ID:        uuid.New().String(),
					Name:      "Test Product",
					CreatedBy: "user-123",
				}
				productService.On("GetByID", "product-123").Return(existingProduct, nil)
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  "forbidden",
		},
		{
			name:          "validation error - negative stock",
			productID:     "product-123",
			currentUserID: "user-123",
			request: map[string]interface{}{
				"stock": -10,
			},
			mockSetup: func(productService *MockProductService) {
				existingProduct := &models.Product{
					ID:        uuid.New().String(),
					Name:      "Test Product",
					CreatedBy: "user-123",
				}
				productService.On("GetByID", "product-123").Return(existingProduct, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockProductService := new(MockProductService)
			tt.mockSetup(mockProductService)

			// Create handler
			handler := handlers.NewProductHandler(mockProductService)

			// Setup router with middleware mock
			router := gin.New()
			router.PATCH("/products/:id/stock", func(c *gin.Context) {
				// Mock middleware setting current user
				user := &models.User{
					ID: uuid.MustParse(tt.currentUserID),
				}
				c.Set("current_user", user)
				handler.UpdateStock(c)
			})

			// Create request
			jsonBody, _ := json.Marshal(tt.request)
			url := "/products/" + tt.productID + "/stock"
			req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response.Error)
			}

			// Verify mock expectations
			mockProductService.AssertExpectations(t)
		})
	}
}
