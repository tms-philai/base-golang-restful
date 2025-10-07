package services

import (
	"testing"

	"base-golang-restful-app/models"
	"base-golang-restful-app/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ProductServiceTestSuite defines the test suite for ProductService
type ProductServiceTestSuite struct {
	suite.Suite
	productService *services.ProductService
}

// SetupTest sets up the test environment before each test
func (suite *ProductServiceTestSuite) SetupTest() {
	suite.productService = services.NewProductService()
}

// TestCreateProductSuccess tests successful product creation
func (suite *ProductServiceTestSuite) TestCreateProductSuccess() {
	// Arrange
	req := models.ProductCreateRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		Category:    "Electronics",
		SKU:         "TEST-001",
		Stock:       100,
	}
	createdBy := "user-123"

	// Act
	product, err := suite.productService.Create(req, createdBy)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), product)
	assert.Equal(suite.T(), req.Name, product.Name)
	assert.Equal(suite.T(), req.Description, product.Description)
	assert.Equal(suite.T(), req.Price, product.Price)
	assert.Equal(suite.T(), req.Category, product.Category)
	assert.Equal(suite.T(), req.SKU, product.SKU)
	assert.Equal(suite.T(), req.Stock, product.Stock)
	assert.Equal(suite.T(), createdBy, product.CreatedBy)
	assert.True(suite.T(), product.IsActive)
	assert.NotEmpty(suite.T(), product.ID)
}

// TestCreateProductDuplicateSKU tests product creation with duplicate SKU
func (suite *ProductServiceTestSuite) TestCreateProductDuplicateSKU() {
	// Arrange
	req1 := models.ProductCreateRequest{
		Name:        "Product 1",
		Description: "Description 1",
		Price:       99.99,
		Category:    "Electronics",
		SKU:         "DUPLICATE-SKU",
		Stock:       100,
	}

	req2 := models.ProductCreateRequest{
		Name:        "Product 2",
		Description: "Description 2",
		Price:       149.99,
		Category:    "Electronics",
		SKU:         "DUPLICATE-SKU", // Same SKU
		Stock:       50,
	}

	// Act
	product1, err1 := suite.productService.Create(req1, "user-123")
	product2, err2 := suite.productService.Create(req2, "user-456")

	// Assert
	assert.NoError(suite.T(), err1)
	assert.NotNil(suite.T(), product1)
	assert.Error(suite.T(), err2)
	assert.Nil(suite.T(), product2)
	assert.Contains(suite.T(), err2.Error(), "SKU already exists")
}


// TestGetProductByID tests product retrieval by ID
func (suite *ProductServiceTestSuite) TestGetProductByID() {
	// Arrange
	req := models.ProductCreateRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		Category:    "Electronics",
		SKU:         "TEST-001",
		Stock:       100,
	}
	createdProduct, _ := suite.productService.Create(req, "user-123")

	// Act
	foundProduct, err := suite.productService.GetByID(createdProduct.ID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), foundProduct)
	assert.Equal(suite.T(), createdProduct.ID, foundProduct.ID)
	assert.Equal(suite.T(), createdProduct.Name, foundProduct.Name)
}

// TestGetProductByIDNotFound tests product retrieval with non-existent ID
func (suite *ProductServiceTestSuite) TestGetProductByIDNotFound() {
	// Act
	product, err := suite.productService.GetByID("nonexistent-id")

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), product)
	assert.Contains(suite.T(), err.Error(), "product not found")
}

// TestGetProductBySKU tests product retrieval by SKU
func (suite *ProductServiceTestSuite) TestGetProductBySKU() {
	// Arrange
	req := models.ProductCreateRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		Category:    "Electronics",
		SKU:         "TEST-001",
		Stock:       100,
	}
	createdProduct, _ := suite.productService.Create(req, "user-123")

	// Act
	foundProduct, err := suite.productService.GetBySKU(createdProduct.SKU)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), foundProduct)
	assert.Equal(suite.T(), createdProduct.SKU, foundProduct.SKU)
}

// TestUpdateProductSuccess tests successful product update
func (suite *ProductServiceTestSuite) TestUpdateProductSuccess() {
	// Arrange
	req := models.ProductCreateRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		Category:    "Electronics",
		SKU:         "TEST-001",
		Stock:       100,
	}
	createdProduct, _ := suite.productService.Create(req, "user-123")

	updateReq := models.ProductUpdateRequest{
		Name:        productStringPtr("Updated Product"),
		Description: productStringPtr("Updated Description"),
		Price:       productFloat64Ptr(149.99),
	}

	// Act
	updatedProduct, err := suite.productService.Update(createdProduct.ID, updateReq)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), updatedProduct)
	assert.Equal(suite.T(), "Updated Product", updatedProduct.Name)
	assert.Equal(suite.T(), "Updated Description", updatedProduct.Description)
	assert.Equal(suite.T(), 149.99, updatedProduct.Price)
	assert.Equal(suite.T(), createdProduct.SKU, updatedProduct.SKU) // Unchanged
}

// TestDeleteProduct tests product deletion (soft delete)
func (suite *ProductServiceTestSuite) TestDeleteProduct() {
	// Arrange
	req := models.ProductCreateRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		Category:    "Electronics",
		SKU:         "TEST-001",
		Stock:       100,
	}
	createdProduct, _ := suite.productService.Create(req, "user-123")

	// Act
	err := suite.productService.Delete(createdProduct.ID)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify product is soft deleted
	deletedProduct, getErr := suite.productService.GetByID(createdProduct.ID)
	assert.NoError(suite.T(), getErr)
	assert.False(suite.T(), deletedProduct.IsActive)
}

// TestListProducts tests product listing with pagination
func (suite *ProductServiceTestSuite) TestListProducts() {
	// Arrange - Create multiple products
	for i := 0; i < 5; i++ {
		req := models.ProductCreateRequest{
			Name:        "Product " + string(rune(i+'1')),
			Description: "Description " + string(rune(i+'1')),
			Price:       99.99 + float64(i*10),
			Category:    "Electronics",
			SKU:         "TEST-00" + string(rune(i+'1')),
			Stock:       100,
		}
		suite.productService.Create(req, "user-123")
	}

	// Act
	query := models.ProductListQuery{
		Page:    1,
		Limit:   3,
		SortBy:  "created_at",
		SortDir: "desc",
	}
	products, pagination, err := suite.productService.List(query)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), products, 3)
	assert.Equal(suite.T(), 1, pagination.Page)
	assert.Equal(suite.T(), 3, pagination.Limit)
	assert.Equal(suite.T(), int64(5), pagination.Total)
	assert.Equal(suite.T(), 2, pagination.TotalPages)
	assert.True(suite.T(), pagination.HasNext)
	assert.False(suite.T(), pagination.HasPrev)
}

// TestListProductsWithSearch tests product listing with search
func (suite *ProductServiceTestSuite) TestListProductsWithSearch() {
	// Arrange
	req1 := models.ProductCreateRequest{
		Name:        "iPhone 15",
		Description: "Latest iPhone",
		Price:       999.99,
		Category:    "Electronics",
		SKU:         "IPHONE-15",
		Stock:       50,
	}
	req2 := models.ProductCreateRequest{
		Name:        "Samsung Galaxy",
		Description: "Android phone",
		Price:       799.99,
		Category:    "Electronics",
		SKU:         "SAMSUNG-GALAXY",
		Stock:       30,
	}
	suite.productService.Create(req1, "user-123")
	suite.productService.Create(req2, "user-123")

	// Act
	query := models.ProductListQuery{
		Page:    1,
		Limit:   10,
		Search:  "iPhone",
		SortBy:  "name",
		SortDir: "asc",
	}
	products, pagination, err := suite.productService.List(query)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), products, 1)
	assert.Equal(suite.T(), "iPhone 15", products[0].Name)
	assert.Equal(suite.T(), int64(1), pagination.Total)
}

// TestGetCategories tests getting all product categories
func (suite *ProductServiceTestSuite) TestGetCategories() {
	// Arrange - Create products with different categories
	categories := []string{"Electronics", "Clothing", "Books", "Electronics"} // Duplicate Electronics
	for i, category := range categories {
		req := models.ProductCreateRequest{
			Name:        "Product " + string(rune(i+'1')),
			Description: "Description",
			Price:       99.99,
			Category:    category,
			SKU:         "TEST-00" + string(rune(i+'1')),
			Stock:       100,
		}
		suite.productService.Create(req, "user-123")
	}

	// Act
	categories = suite.productService.GetCategories()

	// Assert
	assert.Contains(suite.T(), categories, "Electronics")
	assert.Contains(suite.T(), categories, "Clothing")
	assert.Contains(suite.T(), categories, "Books")
	// Should not have duplicates
	assert.Equal(suite.T(), 3, len(categories))
}

// TestUpdateStock tests stock update functionality
func (suite *ProductServiceTestSuite) TestUpdateStock() {
	// Arrange
	req := models.ProductCreateRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		Category:    "Electronics",
		SKU:         "TEST-001",
		Stock:       100,
	}
	createdProduct, _ := suite.productService.Create(req, "user-123")

	// Act
	err := suite.productService.UpdateStock(createdProduct.ID, 150)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify stock was updated
	updatedProduct, getErr := suite.productService.GetByID(createdProduct.ID)
	assert.NoError(suite.T(), getErr)
	assert.Equal(suite.T(), 150, updatedProduct.Stock)
}

// TestCheckStock tests stock checking functionality
func (suite *ProductServiceTestSuite) TestCheckStock() {
	// Arrange
	req := models.ProductCreateRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		Category:    "Electronics",
		SKU:         "TEST-001",
		Stock:       100,
	}
	createdProduct, _ := suite.productService.Create(req, "user-123")

	// Act & Assert - Available stock
	available, err := suite.productService.CheckStock(createdProduct.ID, 50)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), available)

	// Act & Assert - Insufficient stock
	available, err = suite.productService.CheckStock(createdProduct.ID, 150)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), available)
}

// Helper functions
func productStringPtr(s string) *string {
	return &s
}

func productFloat64Ptr(f float64) *float64 {
	return &f
}

// TestProductServiceTestSuite runs the test suite
func TestProductServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceTestSuite))
}
