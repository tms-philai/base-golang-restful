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

// ProductHandlerTestSuite defines the test suite for ProductHandler
type ProductHandlerTestSuite struct {
	suite.Suite
	helper            *testhelpers.TestHelper
	mockProductSvc    *mocks.MockProductService
	productHandler    *handlers.ProductHandler
	router            *gin.Engine
	testUser          *models.User
	testAdmin         *models.User
	testProduct       *models.Product
	testToken         string
	adminToken        string
}

// SetupTest sets up the test environment before each test
func (suite *ProductHandlerTestSuite) SetupTest() {
	suite.helper = testhelpers.NewTestHelper()
	suite.mockProductSvc = new(mocks.MockProductService)
	suite.productHandler = handlers.NewProductHandler(suite.mockProductSvc)

	// Create test users
	suite.testUser = &models.User{
		ID:       "user-123",
		Username: "testuser",
		Role:     "user",
		IsActive: true,
	}

	suite.testAdmin = &models.User{
		ID:       "admin-123",
		Username: "admin",
		Role:     "admin",
		IsActive: true,
	}

	// Create test product
	suite.testProduct = &models.Product{
		ID:          "product-123",
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		Category:    "Electronics",
		SKU:         "TEST-001",
		Stock:       100,
		IsActive:    true,
		CreatedBy:   suite.testUser.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	suite.testToken = "valid-user-token"
	suite.adminToken = "valid-admin-token"

	// Setup router with mock middleware
	suite.router = gin.New()
	
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

	products := suite.router.Group("/products")
	{
		products.GET("", suite.productHandler.ListProducts)
		products.GET("/categories", suite.productHandler.GetCategories)
		products.GET("/:id", suite.productHandler.GetProduct)
		products.POST("", authMiddleware, suite.productHandler.CreateProduct)
		products.PUT("/:id", authMiddleware, suite.productHandler.UpdateProduct)
		products.DELETE("/:id", authMiddleware, suite.productHandler.DeleteProduct)
		products.PATCH("/:id/stock", authMiddleware, suite.productHandler.UpdateStock)
	}

	suite.helper.Router = suite.router
}

// TearDownTest cleans up after each test
func (suite *ProductHandlerTestSuite) TearDownTest() {
	suite.mockProductSvc.AssertExpectations(suite.T())
}

// TestGetProductSuccess tests successful product retrieval
func (suite *ProductHandlerTestSuite) TestGetProductSuccess() {
	// Arrange
	productID := "product-123"
	suite.mockProductSvc.On("GetByID", productID).Return(suite.testProduct, nil)

	// Act
	w := suite.helper.MakeRequest("GET", "/products/"+productID, nil, nil)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.Product
	err := suite.helper.ParseJSONResponse(w, &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.testProduct.Name, response.Name)
	assert.Equal(suite.T(), suite.testProduct.SKU, response.SKU)
}

// TestGetProductNotFound tests product retrieval when product doesn't exist
func (suite *ProductHandlerTestSuite) TestGetProductNotFound() {
	// Arrange
	productID := "nonexistent-product"
	suite.mockProductSvc.On("GetByID", productID).Return((*models.Product)(nil), errors.New("product not found"))

	// Act
	w := suite.helper.MakeRequest("GET", "/products/"+productID, nil, nil)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	suite.helper.AssertErrorResponse(suite.T(), w, http.StatusNotFound, "product_not_found")
}

// TestCreateProductSuccess tests successful product creation
func (suite *ProductHandlerTestSuite) TestCreateProductSuccess() {
	// Arrange
	createReq := models.ProductCreateRequest{
		Name:        "New Product",
		Description: "New Description",
		Price:       149.99,
		Category:    "Electronics",
		SKU:         "NEW-001",
		Stock:       50,
	}

	expectedProduct := &models.Product{
		ID:          "new-product-123",
		Name:        "New Product",
		Description: "New Description",
		Price:       149.99,
		Category:    "Electronics",
		SKU:         "NEW-001",
		Stock:       50,
		IsActive:    true,
		CreatedBy:   suite.testUser.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	suite.mockProductSvc.On("Create", mock.AnythingOfType("models.ProductCreateRequest"), suite.testUser.ID).Return(expectedProduct, nil)

	// Act
	w := suite.helper.MakeAuthenticatedRequest("POST", "/products", createReq, suite.testToken)

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response models.Product
	err := suite.helper.ParseJSONResponse(w, &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedProduct.Name, response.Name)
	assert.Equal(suite.T(), expectedProduct.SKU, response.SKU)
}

// TestCreateProductValidationError tests product creation with validation errors
func (suite *ProductHandlerTestSuite) TestCreateProductValidationError() {
	// Arrange
	invalidReq := models.ProductCreateRequest{
		Name:  "", // Required field missing
		Price: -10, // Invalid price
	}

	// Act
	w := suite.helper.MakeAuthenticatedRequest("POST", "/products", invalidReq, suite.testToken)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	suite.helper.AssertErrorResponse(suite.T(), w, http.StatusBadRequest, "validation_error")
}

// TestCreateProductSKUExists tests product creation with existing SKU
func (suite *ProductHandlerTestSuite) TestCreateProductSKUExists() {
	// Arrange
	createReq := models.ProductCreateRequest{
		Name:        "New Product",
		Description: "New Description",
		Price:       149.99,
		Category:    "Electronics",
		SKU:         "EXISTING-SKU",
		Stock:       50,
	}

	suite.mockProductSvc.On("Create", mock.AnythingOfType("models.ProductCreateRequest"), suite.testUser.ID).Return((*models.Product)(nil), errors.New("SKU already exists"))

	// Act
	w := suite.helper.MakeAuthenticatedRequest("POST", "/products", createReq, suite.testToken)

	// Assert
	assert.Equal(suite.T(), http.StatusConflict, w.Code)
	suite.helper.AssertErrorResponse(suite.T(), w, http.StatusConflict, "product_creation_failed")
}

// TestUpdateProductSuccess tests successful product update by owner
func (suite *ProductHandlerTestSuite) TestUpdateProductSuccess() {
	// Arrange
	productID := suite.testProduct.ID
	updateReq := models.ProductUpdateRequest{
		Name:  stringPtr("Updated Product"),
		Price: float64Ptr(199.99),
	}

	updatedProduct := *suite.testProduct
	updatedProduct.Name = "Updated Product"
	updatedProduct.Price = 199.99

	suite.mockProductSvc.On("GetByID", productID).Return(suite.testProduct, nil)
	suite.mockProductSvc.On("Update", productID, mock.AnythingOfType("models.ProductUpdateRequest")).Return(&updatedProduct, nil)

	// Act
	w := suite.helper.MakeAuthenticatedRequest("PUT", "/products/"+productID, updateReq, suite.testToken)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.Product
	err := suite.helper.ParseJSONResponse(w, &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated Product", response.Name)
	assert.Equal(suite.T(), 199.99, response.Price)
}

// TestUpdateProductForbidden tests updating product by non-owner
func (suite *ProductHandlerTestSuite) TestUpdateProductForbidden() {
	// Arrange
	productID := suite.testProduct.ID
	updateReq := models.ProductUpdateRequest{
		Name: stringPtr("Updated Product"),
	}

	// Product owned by different user
	otherProduct := *suite.testProduct
	otherProduct.CreatedBy = "other-user-123"

	suite.mockProductSvc.On("GetByID", productID).Return(&otherProduct, nil)

	// Act
	w := suite.helper.MakeAuthenticatedRequest("PUT", "/products/"+productID, updateReq, suite.testToken)

	// Assert
	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
	suite.helper.AssertErrorResponse(suite.T(), w, http.StatusForbidden, "forbidden")
}

// TestDeleteProductSuccess tests successful product deletion by owner
func (suite *ProductHandlerTestSuite) TestDeleteProductSuccess() {
	// Arrange
	productID := suite.testProduct.ID

	suite.mockProductSvc.On("GetByID", productID).Return(suite.testProduct, nil)
	suite.mockProductSvc.On("Delete", productID).Return(nil)

	// Act
	w := suite.helper.MakeAuthenticatedRequest("DELETE", "/products/"+productID, nil, suite.testToken)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.helper.AssertSuccessResponse(suite.T(), w, http.StatusOK, "Product deleted successfully")
}

// TestDeleteProductByAdmin tests successful product deletion by admin
func (suite *ProductHandlerTestSuite) TestDeleteProductByAdmin() {
	// Arrange
	productID := suite.testProduct.ID

	suite.mockProductSvc.On("GetByID", productID).Return(suite.testProduct, nil)
	suite.mockProductSvc.On("Delete", productID).Return(nil)

	// Act
	w := suite.helper.MakeAuthenticatedRequest("DELETE", "/products/"+productID, nil, suite.adminToken)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.helper.AssertSuccessResponse(suite.T(), w, http.StatusOK, "Product deleted successfully")
}

// TestListProductsSuccess tests successful product listing
func (suite *ProductHandlerTestSuite) TestListProductsSuccess() {
	// Arrange
	products := []models.Product{*suite.testProduct}
	pagination := models.PaginationMetadata{
		Page:       1,
		Limit:      10,
		Total:      1,
		TotalPages: 1,
		HasNext:    false,
		HasPrev:    false,
	}

	expectedQuery := models.ProductListQuery{
		Page:     1,
		Limit:    10,
		Category: "",
		Search:   "",
		SortBy:   "created_at",
		SortDir:  "desc",
	}

	suite.mockProductSvc.On("List", expectedQuery).Return(products, pagination, nil)

	// Act
	w := suite.helper.MakeRequest("GET", "/products?page=1&limit=10", nil, nil)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.ProductListResponse
	err := suite.helper.ParseJSONResponse(w, &response)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), response.Products, 1)
	assert.Equal(suite.T(), pagination.Total, response.Pagination.Total)
}

// TestGetCategoriesSuccess tests successful categories retrieval
func (suite *ProductHandlerTestSuite) TestGetCategoriesSuccess() {
	// Arrange
	categories := []string{"Electronics", "Clothing", "Books"}
	suite.mockProductSvc.On("GetCategories").Return(categories)

	// Act
	w := suite.helper.MakeRequest("GET", "/products/categories", nil, nil)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response struct {
		Categories []string `json:"categories"`
	}
	err := suite.helper.ParseJSONResponse(w, &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), categories, response.Categories)
}

// TestUpdateStockSuccess tests successful stock update
func (suite *ProductHandlerTestSuite) TestUpdateStockSuccess() {
	// Arrange
	productID := suite.testProduct.ID
	stockReq := struct {
		Stock int `json:"stock"`
	}{Stock: 150}

	suite.mockProductSvc.On("GetByID", productID).Return(suite.testProduct, nil)
	suite.mockProductSvc.On("UpdateStock", productID, 150).Return(nil)

	// Act
	w := suite.helper.MakeAuthenticatedRequest("PATCH", "/products/"+productID+"/stock", stockReq, suite.testToken)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.helper.AssertSuccessResponse(suite.T(), w, http.StatusOK, "Stock updated successfully")
}

// Helper functions
func float64Ptr(f float64) *float64 {
	return &f
}

// TestProductHandlerTestSuite runs the test suite
func TestProductHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ProductHandlerTestSuite))
}
