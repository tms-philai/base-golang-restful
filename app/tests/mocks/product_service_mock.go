package mocks

import (
	"base-golang-restful-app/models"

	"github.com/stretchr/testify/mock"
)

// MockProductService is a mock implementation of ProductService
type MockProductService struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockProductService) Create(req models.ProductCreateRequest, createdBy string) (*models.Product, error) {
	args := m.Called(req, createdBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

// GetByID mocks the GetByID method
func (m *MockProductService) GetByID(id string) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

// GetBySKU mocks the GetBySKU method
func (m *MockProductService) GetBySKU(sku string) (*models.Product, error) {
	args := m.Called(sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

// Update mocks the Update method
func (m *MockProductService) Update(id string, req models.ProductUpdateRequest) (*models.Product, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

// Delete mocks the Delete method
func (m *MockProductService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// List mocks the List method
func (m *MockProductService) List(query models.ProductListQuery) ([]models.Product, models.PaginationMetadata, error) {
	args := m.Called(query)
	return args.Get(0).([]models.Product), args.Get(1).(models.PaginationMetadata), args.Error(2)
}

// GetCategories mocks the GetCategories method
func (m *MockProductService) GetCategories() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

// UpdateStock mocks the UpdateStock method
func (m *MockProductService) UpdateStock(id string, stock int) error {
	args := m.Called(id, stock)
	return args.Error(0)
}

// CheckStock mocks the CheckStock method
func (m *MockProductService) CheckStock(id string, quantity int) (bool, error) {
	args := m.Called(id, quantity)
	return args.Bool(0), args.Error(1)
}
