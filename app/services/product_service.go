package services

import (
	"context"
	"errors"
	"fmt"

	"base-golang-restful-app/models"
	"base-golang-restful-app/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) Create(req models.ProductCreateRequest, createdBy string) (*models.Product, error) {
	ctx := context.Background()

	// Check if SKU already exists
	_, err := s.productRepo.FindBySKU(ctx, req.SKU)
	if err == nil {
		return nil, errors.New("SKU already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check SKU: %w", err)
	}

	// Create new product
	product := models.NewProduct(req, createdBy)

	// Save to database
	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (s *ProductService) GetByID(id string) (*models.Product, error) {
	ctx := context.Background()

	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	var product models.Product
	if err := s.productRepo.FindByID(ctx, productID, &product); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

func (s *ProductService) GetBySKU(sku string) (*models.Product, error) {
	ctx := context.Background()

	product, err := s.productRepo.FindBySKU(ctx, sku)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

func (s *ProductService) Update(id string, req models.ProductUpdateRequest) (*models.Product, error) {
	ctx := context.Background()

	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	// Get existing product
	var product models.Product
	if err := s.productRepo.FindByID(ctx, productID, &product); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Check SKU uniqueness if SKU is being updated
	if req.SKU != nil && *req.SKU != product.SKU {
		_, err := s.productRepo.FindBySKU(ctx, *req.SKU)
		if err == nil {
			return nil, errors.New("SKU already exists")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to check SKU: %w", err)
		}
	}

	// Apply updates
	product.Update(req)

	// Save updates
	if err := s.productRepo.Update(ctx, &product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return &product, nil
}

func (s *ProductService) Delete(id string) error {
	ctx := context.Background()

	productID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid product ID")
	}

	// Get product to verify it exists
	var product models.Product
	if err := s.productRepo.FindByID(ctx, productID, &product); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Soft delete by setting IsActive to false
	product.IsActive = false
	if err := s.productRepo.Update(ctx, &product); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func (s *ProductService) List(query models.ProductListQuery) ([]models.Product, models.PaginationMetadata, error) {
	ctx := context.Background()

	// If search or category filter is specified, use search
	if query.Search != "" || query.Category != "" {
		return s.listWithFilters(ctx, query)
	}

	// Otherwise use pagination
	products, total, err := s.productRepo.FindAllWithPagination(ctx, query.Page, query.Limit)
	if err != nil {
		return nil, models.PaginationMetadata{}, fmt.Errorf("failed to list products: %w", err)
	}

	pagination := models.NewPaginationMetadata(query.Page, query.Limit, total)
	return products, pagination, nil
}

func (s *ProductService) listWithFilters(ctx context.Context, query models.ProductListQuery) ([]models.Product, models.PaginationMetadata, error) {
	var products []models.Product
	var err error

	// Apply filters
	if query.Category != "" {
		products, err = s.productRepo.FindByCategory(ctx, query.Category)
	} else if query.Search != "" {
		products, err = s.productRepo.Search(ctx, query.Search)
	}

	if err != nil {
		return nil, models.PaginationMetadata{}, fmt.Errorf("failed to list products: %w", err)
	}

	total := int64(len(products))

	// Apply pagination manually for filtered results
	offset := (query.Page - 1) * query.Limit
	end := offset + query.Limit

	if offset > len(products) {
		return []models.Product{}, models.NewPaginationMetadata(query.Page, query.Limit, total), nil
	}

	if end > len(products) {
		end = len(products)
	}

	paginatedProducts := products[offset:end]
	pagination := models.NewPaginationMetadata(query.Page, query.Limit, total)

	return paginatedProducts, pagination, nil
}

func (s *ProductService) GetCategories() []string {
	ctx := context.Background()

	products, err := s.productRepo.FindByCategory(ctx, "")
	if err != nil {
		return []string{}
	}

	categoryMap := make(map[string]bool)
	for _, product := range products {
		if product.IsActive {
			categoryMap[product.Category] = true
		}
	}

	var categories []string
	for category := range categoryMap {
		if category != "" {
			categories = append(categories, category)
		}
	}

	return categories
}

func (s *ProductService) UpdateStock(id string, stock int) error {
	ctx := context.Background()

	if stock < 0 {
		return errors.New("stock cannot be negative")
	}

	productID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid product ID")
	}

	// Check if product exists
	var product models.Product
	if err := s.productRepo.FindByID(ctx, productID, &product); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Update stock
	if err := s.productRepo.UpdateStock(ctx, productID, stock); err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	return nil
}

func (s *ProductService) CheckStock(id string, quantity int) (bool, error) {
	product, err := s.GetByID(id)
	if err != nil {
		return false, err
	}

	return product.Stock >= quantity, nil
}
