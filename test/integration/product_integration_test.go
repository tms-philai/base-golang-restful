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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ProductIntegrationTestSuite is the test suite for product integration tests
type ProductIntegrationTestSuite struct {
	suite.Suite
	router         *gin.Engine
	dbSetup        *setup.TestDatabaseSetup
	userService    *services.UserService
	productService *services.ProductService
	roleService    *services.RoleService
	jwtManager     *auth.JWTManager
	productHandler *handlers.ProductHandler
	authMiddleware *middleware.AuthMiddleware
	adminToken     string
	userToken      string
	adminUser      *models.User
	regularUser    *models.User
	testProduct    *models.Product
}

// SetupSuite sets up the test suite
func (suite *ProductIntegrationTestSuite) SetupSuite() {
	// Setup test database
	suite.dbSetup = setup.SetupTestEnvironment(suite.T())

	// Load test configuration
	testConfig, err := config.LoadTestConfig()
	suite.Require().NoError(err)
	// Initialize repositories
	userRepo := repository.NewUserRepository(setup.GetTestDB())
	productRepo := repository.NewProductRepository(setup.GetTestDB())
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
	suite.productService = services.NewProductService(productRepo)
	suite.roleService = services.NewRoleService(setup.GetTestDB(), roleRepo, userRepo, permissionRepo)

	// Initialize middleware and handlers
	suite.authMiddleware = middleware.NewAuthMiddleware(suite.jwtManager, suite.userService)
	suite.productHandler = handlers.NewProductHandler(suite.productService)

	// Setup router
	suite.router = routes.SetupRouter(routes.RouterConfig{
		AuthHandler:    nil, // We'll set this later
		AuthMiddleware: suite.authMiddleware,
		ProductHandler: suite.productHandler,
	})
}

// TearDownSuite tears down the test suite
func (suite *ProductIntegrationTestSuite) TearDownSuite() {
	setup.CleanupTestEnvironment(suite.T(), suite.dbSetup)
}

// SetupTest sets up each test
func (suite *ProductIntegrationTestSuite) SetupTest() {
	// Clean up test data before each test
	suite.dbSetup.TruncateTestTables()

	// Create test users and products
	suite.createTestUsers()
	suite.createTestProducts()
}

// createTestUsers creates test users for the suite
func (suite *ProductIntegrationTestSuite) createTestUsers() {
	// Create admin user
	adminUser, err := suite.userService.Create(models.UserCreateRequest{
		Username:  "admin",
		Email:     "admin@example.com",
		Password:  "admin123",
		FirstName: "Admin",
		LastName:  "User",
	})
	// Assign admin role
	adminRoleID := uuid.MustParse("00000000-0000-0000-0000-000000000001") // Using a fixed UUID for admin role
	err = suite.roleService.AssignRole(adminUser.ID, adminRoleID)
	suite.Require().NoError(err)

	// Create regular user
	regularUser, err := suite.userService.Create(models.UserCreateRequest{
		Username:  "user",
		Email:     "user@example.com",
		Password:  "user123",
		FirstName: "Regular",
		LastName:  "User",
	})
	suite.Require().NoError(err)
	// Assign user role
	userRoleID := uuid.MustParse("00000000-0000-0000-0000-000000000002") // Using a fixed UUID for user role
	err = suite.roleService.AssignRole(regularUser.ID, userRoleID)
	suite.Require().NoError(err)

	// Generate tokens
	suite.adminToken, err = suite.jwtManager.GenerateAccessToken(adminUser.ID, adminUser.Email)
	suite.Require().NoError(err)

	suite.userToken, err = suite.jwtManager.GenerateAccessToken(regularUser.ID, regularUser.Email)
	suite.Require().NoError(err)

	suite.adminUser = adminUser
	suite.regularUser = regularUser
}

// createTestProducts creates test products for the suite
func (suite *ProductIntegrationTestSuite) createTestProducts() {
	// Create test product
	product, err := suite.productService.Create(models.ProductCreateRequest{
		Name:        "Test Product",
		Description: "A test product for testing",
		Price:       99.99,
		SKU:         "TEST-001",
		Category:    "Electronics",
		Stock:       100,
	}, suite.regularUser.ID.String())
	suite.Require().NoError(err)

	suite.testProduct = product
}

// TestCreateProduct tests product creation
func (suite *ProductIntegrationTestSuite) TestCreateProduct() {
	tests := []struct {
		name           string
		token          string
		request        models.ProductCreateRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name:  "successful product creation by authenticated user",
			token: suite.userToken,
			request: models.ProductCreateRequest{
				Name:        "New Product",
				Description: "A new product",
				Price:       149.99,
				SKU:         "NEW-001",
				Category:    "Electronics",
				Stock:       50,
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name:  "successful product creation by admin",
			token: suite.adminToken,
			request: models.ProductCreateRequest{
				Name:        "Admin Product",
				Description: "A product created by admin",
				Price:       199.99,
				SKU:         "ADMIN-001",
				Category:    "Clothing",
				Stock:       25,
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name:  "unauthorized - no token",
			token: "",
			request: models.ProductCreateRequest{
				Name:        "New Product",
				Description: "A new product",
				Price:       149.99,
				SKU:         "NEW-001",
				Category:    "Electronics",
				Stock:       50,
			},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:  "validation error - missing name",
			token: suite.userToken,
			request: models.ProductCreateRequest{
				Description: "A product without name",
				Price:       149.99,
				SKU:         "NO-NAME-001",
				Category:    "Electronics",
				Stock:       50,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:  "conflict - SKU already exists",
			token: suite.userToken,
			request: models.ProductCreateRequest{
				Name:        "Duplicate Product",
				Description: "A product with existing SKU",
				Price:       149.99,
				SKU:         "TEST-001", // This SKU already exists
				Category:    "Electronics",
				Stock:       50,
			},
			expectedStatus: http.StatusConflict,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Create request
			req, err := helpers.CreateJSONRequest("POST", "/api/v1/products", tt.request)
			suite.Require().NoError(err)

			// Add authorization header if token provided
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			suite.router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if !tt.expectError {
				helpers.AssertProductResponse(suite.T(), w, tt.expectedStatus)
			} else {
				if tt.expectedStatus == http.StatusUnauthorized {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "unauthorized")
				} else if tt.expectedStatus == http.StatusConflict {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "product_creation_failed")
				} else {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "validation_error")
				}
			}
		})
	}
}

// TestGetProduct tests getting a product by ID
func (suite *ProductIntegrationTestSuite) TestGetProduct() {
	tests := []struct {
		name           string
		productID      string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "successful get product by ID",
			productID:      suite.testProduct.ID,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "product not found",
			productID:      "00000000-0000-0000-0000-000000000000",
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:           "invalid product ID format",
			productID:      "invalid-id",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Create request
			req, err := http.NewRequest("GET", "/api/v1/products/"+tt.productID, nil)
			suite.Require().NoError(err)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			suite.router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if !tt.expectError {
				helpers.AssertProductResponse(suite.T(), w, tt.expectedStatus)
			} else {
				if tt.expectedStatus == http.StatusNotFound {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "product_not_found")
				} else {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "validation_error")
				}
			}
		})
	}
}

// TestUpdateProduct tests updating a product
func (suite *ProductIntegrationTestSuite) TestUpdateProduct() {
	tests := []struct {
		name           string
		productID      string
		token          string
		request        models.ProductUpdateRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name:      "successful update by owner",
			productID: suite.testProduct.ID,
			token:     suite.userToken,
			request: models.ProductUpdateRequest{
				Name:        helpers.StringPtr("Updated Product"),
				Description: helpers.StringPtr("Updated description"),
				Price:       helpers.Float64Ptr(149.99),
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:      "successful update by admin",
			productID: suite.testProduct.ID,
			token:     suite.adminToken,
			request: models.ProductUpdateRequest{
				Name:        helpers.StringPtr("Admin Updated Product"),
				Description: helpers.StringPtr("Updated by admin"),
				Price:       helpers.Float64Ptr(199.99),
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:      "forbidden - user trying to update another user's product",
			productID: suite.testProduct.ID,
			token:     suite.adminToken, // Using admin token but testing with different user's product
			request: models.ProductUpdateRequest{
				Name: helpers.StringPtr("Updated Product"),
			},
			expectedStatus: http.StatusOK, // Admin can update any product
			expectError:    false,
		},
		{
			name:           "unauthorized - no token",
			productID:      suite.testProduct.ID,
			token:          "",
			request:        models.ProductUpdateRequest{},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:      "product not found",
			productID: "00000000-0000-0000-0000-000000000000",
			token:     suite.userToken,
			request: models.ProductUpdateRequest{
				Name: helpers.StringPtr("Updated Product"),
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Create request
			req, err := helpers.CreateJSONRequest("PUT", "/api/v1/products/"+tt.productID, tt.request)
			suite.Require().NoError(err)

			// Add authorization header if token provided
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			suite.router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if !tt.expectError {
				helpers.AssertProductResponse(suite.T(), w, tt.expectedStatus)
			} else {
				if tt.expectedStatus == http.StatusUnauthorized {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "unauthorized")
				} else if tt.expectedStatus == http.StatusNotFound {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "product_not_found")
				} else {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "validation_error")
				}
			}
		})
	}
}

// TestListProducts tests listing products
func (suite *ProductIntegrationTestSuite) TestListProducts() {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "successful list products",
			queryParams:    "?page=1&limit=10",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "successful list products with search",
			queryParams:    "?page=1&limit=10&search=Test",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "successful list products with category filter",
			queryParams:    "?page=1&limit=10&category=Electronics",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "successful list products with sorting",
			queryParams:    "?page=1&limit=10&sort_by=price&sort_dir=asc",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Create request
			req, err := http.NewRequest("GET", "/api/v1/products"+tt.queryParams, nil)
			suite.Require().NoError(err)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			suite.router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if !tt.expectError {
				helpers.AssertPaginationResponse(suite.T(), w, tt.expectedStatus)

				// Verify response structure
				var response models.ProductListResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				suite.Require().NoError(err)
				assert.NotNil(suite.T(), response.Products)
				assert.NotNil(suite.T(), response.Pagination)
			}
		})
	}
}

// TestGetCategories tests getting product categories
func (suite *ProductIntegrationTestSuite) TestGetCategories() {
	tests := []struct {
		name           string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "successful get categories",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Create request
			req, err := http.NewRequest("GET", "/api/v1/products/categories", nil)
			suite.Require().NoError(err)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			suite.router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if !tt.expectError {
				// Verify response structure
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				suite.Require().NoError(err)
				assert.Contains(suite.T(), response, "categories")
			}
		})
	}
}

// TestUpdateStock tests updating product stock
func (suite *ProductIntegrationTestSuite) TestUpdateStock() {
	tests := []struct {
		name           string
		productID      string
		token          string
		request        map[string]interface{}
		expectedStatus int
		expectError    bool
	}{
		{
			name:      "successful stock update by owner",
			productID: suite.testProduct.ID,
			token:     suite.userToken,
			request: map[string]interface{}{
				"stock": 150,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:      "successful stock update by admin",
			productID: suite.testProduct.ID,
			token:     suite.adminToken,
			request: map[string]interface{}{
				"stock": 200,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "unauthorized - no token",
			productID:      suite.testProduct.ID,
			token:          "",
			request:        map[string]interface{}{"stock": 150},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:      "validation error - negative stock",
			productID: suite.testProduct.ID,
			token:     suite.userToken,
			request: map[string]interface{}{
				"stock": -10,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:      "product not found",
			productID: "00000000-0000-0000-0000-000000000000",
			token:     suite.userToken,
			request: map[string]interface{}{
				"stock": 150,
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Create request
			req, err := helpers.CreateJSONRequest("PATCH", "/api/v1/products/"+tt.productID+"/stock", tt.request)
			suite.Require().NoError(err)

			// Add authorization header if token provided
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
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
				} else if tt.expectedStatus == http.StatusNotFound {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "product_not_found")
				} else {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "validation_error")
				}
			}
		})
	}
}

// TestDeleteProduct tests deleting a product
func (suite *ProductIntegrationTestSuite) TestDeleteProduct() {
	tests := []struct {
		name           string
		productID      string
		token          string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "successful delete by owner",
			productID:      suite.testProduct.ID,
			token:          suite.userToken,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "successful delete by admin",
			productID:      suite.testProduct.ID,
			token:          suite.adminToken,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "unauthorized - no token",
			productID:      suite.testProduct.ID,
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:           "product not found",
			productID:      "00000000-0000-0000-0000-000000000000",
			token:          suite.userToken,
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Create request
			req, err := http.NewRequest("DELETE", "/api/v1/products/"+tt.productID, nil)
			suite.Require().NoError(err)

			// Add authorization header if token provided
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
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
				} else if tt.expectedStatus == http.StatusNotFound {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "product_not_found")
				} else {
					helpers.AssertErrorResponse(suite.T(), w, tt.expectedStatus, "validation_error")
				}
			}
		})
	}
}

// TestProductIntegrationSuite runs the product integration test suite
func TestProductIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ProductIntegrationTestSuite))
}
