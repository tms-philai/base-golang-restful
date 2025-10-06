package services

import (
	"errors"
	"sort"
	"strings"
	"sync"

	"base-golang-restful-app/models"
)

// ProductService handles product-related business logic
type ProductService struct {
	products map[string]*models.Product // In-memory storage for demo
	mutex    sync.RWMutex
}

// NewProductService creates a new product service
func NewProductService() *ProductService {
	return &ProductService{
		products: make(map[string]*models.Product),
		mutex:    sync.RWMutex{},
	}
}

// Create creates a new product
func (s *ProductService) Create(req models.ProductCreateRequest, createdBy string) (*models.Product, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if SKU already exists
	for _, product := range s.products {
		if product.SKU == req.SKU {
			return nil, errors.New("SKU already exists")
		}
	}

	// Create new product
	product := models.NewProduct(req, createdBy)

	// Store product
	s.products[product.ID] = product

	return product, nil
}

// GetByID retrieves a product by ID
func (s *ProductService) GetByID(id string) (*models.Product, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	product, exists := s.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}

	return product, nil
}

// GetBySKU retrieves a product by SKU
func (s *ProductService) GetBySKU(sku string) (*models.Product, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, product := range s.products {
		if product.SKU == sku {
			return product, nil
		}
	}

	return nil, errors.New("product not found")
}

// Update updates a product
func (s *ProductService) Update(id string, req models.ProductUpdateRequest) (*models.Product, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	product, exists := s.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}

	// Check for SKU conflicts if SKU is being updated
	if req.SKU != nil {
		for _, existingProduct := range s.products {
			if existingProduct.ID != id && existingProduct.SKU == *req.SKU {
				return nil, errors.New("SKU already exists")
			}
		}
	}

	// Apply updates
	product.Update(req)

	return product, nil
}

// Delete deletes a product (soft delete by setting IsActive to false)
func (s *ProductService) Delete(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	product, exists := s.products[id]
	if !exists {
		return errors.New("product not found")
	}

	product.IsActive = false
	return nil
}

// List retrieves products with filtering, searching, sorting, and pagination
func (s *ProductService) List(query models.ProductListQuery) ([]models.Product, models.PaginationMetadata, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Convert map to slice and apply filters
	var filteredProducts []models.Product
	for _, product := range s.products {
		if !product.IsActive {
			continue // Skip inactive products
		}

		// Apply category filter
		if query.Category != "" && product.Category != query.Category {
			continue
		}

		// Apply search filter
		if query.Search != "" {
			searchLower := strings.ToLower(query.Search)
			if !strings.Contains(strings.ToLower(product.Name), searchLower) &&
				!strings.Contains(strings.ToLower(product.Description), searchLower) &&
				!strings.Contains(strings.ToLower(product.SKU), searchLower) {
				continue
			}
		}

		filteredProducts = append(filteredProducts, *product)
	}

	// Apply sorting
	s.sortProducts(filteredProducts, query.SortBy, query.SortDir)

	total := int64(len(filteredProducts))

	// Apply pagination
	offset := (query.Page - 1) * query.Limit
	end := offset + query.Limit

	if offset > len(filteredProducts) {
		return []models.Product{}, models.NewPaginationMetadata(query.Page, query.Limit, total), nil
	}

	if end > len(filteredProducts) {
		end = len(filteredProducts)
	}

	products := filteredProducts[offset:end]
	pagination := models.NewPaginationMetadata(query.Page, query.Limit, total)

	return products, pagination, nil
}

// sortProducts sorts products based on the specified field and direction
func (s *ProductService) sortProducts(products []models.Product, sortBy, sortDir string) {
	sort.Slice(products, func(i, j int) bool {
		var less bool

		switch sortBy {
		case "name":
			less = products[i].Name < products[j].Name
		case "price":
			less = products[i].Price < products[j].Price
		case "category":
			less = products[i].Category < products[j].Category
		case "sku":
			less = products[i].SKU < products[j].SKU
		case "stock":
			less = products[i].Stock < products[j].Stock
		case "created_at":
			less = products[i].CreatedAt.Before(products[j].CreatedAt)
		case "updated_at":
			less = products[i].UpdatedAt.Before(products[j].UpdatedAt)
		default:
			// Default sort by created_at
			less = products[i].CreatedAt.Before(products[j].CreatedAt)
		}

		if sortDir == "desc" {
			return !less
		}
		return less
	})
}

// GetCategories retrieves all unique product categories
func (s *ProductService) GetCategories() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	categoryMap := make(map[string]bool)
	for _, product := range s.products {
		if product.IsActive {
			categoryMap[product.Category] = true
		}
	}

	var categories []string
	for category := range categoryMap {
		categories = append(categories, category)
	}

	sort.Strings(categories)
	return categories
}

// UpdateStock updates product stock
func (s *ProductService) UpdateStock(id string, stock int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	product, exists := s.products[id]
	if !exists {
		return errors.New("product not found")
	}

	if stock < 0 {
		return errors.New("stock cannot be negative")
	}

	product.Stock = stock
	return nil
}

// CheckStock checks if product has sufficient stock
func (s *ProductService) CheckStock(id string, quantity int) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	product, exists := s.products[id]
	if !exists {
		return false, errors.New("product not found")
	}

	return product.Stock >= quantity, nil
}
