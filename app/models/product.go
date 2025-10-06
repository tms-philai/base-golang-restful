package models

import (
	"time"

	"github.com/google/uuid"
)

// Product represents a product in the system
type Product struct {
	ID          string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string    `json:"name" binding:"required,min=1,max=200" example:"iPhone 15 Pro"`
	Description string    `json:"description" example:"Latest iPhone with advanced features"`
	Price       float64   `json:"price" binding:"required,min=0" example:"999.99"`
	Category    string    `json:"category" binding:"required" example:"Electronics"`
	SKU         string    `json:"sku" binding:"required" example:"IP15P-128-BLK"`
	Stock       int       `json:"stock" binding:"min=0" example:"100"`
	IsActive    bool      `json:"is_active" example:"true"`
	CreatedBy   string    `json:"created_by" example:"550e8400-e29b-41d4-a716-446655440000"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// ProductCreateRequest represents the request payload for creating a product
type ProductCreateRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=200" example:"iPhone 15 Pro"`
	Description string  `json:"description" example:"Latest iPhone with advanced features"`
	Price       float64 `json:"price" binding:"required,min=0" example:"999.99"`
	Category    string  `json:"category" binding:"required" example:"Electronics"`
	SKU         string  `json:"sku" binding:"required" example:"IP15P-128-BLK"`
	Stock       int     `json:"stock" binding:"min=0" example:"100"`
}

// ProductUpdateRequest represents the request payload for updating a product
type ProductUpdateRequest struct {
	Name        *string  `json:"name,omitempty" binding:"omitempty,min=1,max=200" example:"iPhone 15 Pro"`
	Description *string  `json:"description,omitempty" example:"Latest iPhone with advanced features"`
	Price       *float64 `json:"price,omitempty" binding:"omitempty,min=0" example:"999.99"`
	Category    *string  `json:"category,omitempty" example:"Electronics"`
	SKU         *string  `json:"sku,omitempty" example:"IP15P-128-BLK"`
	Stock       *int     `json:"stock,omitempty" binding:"omitempty,min=0" example:"100"`
	IsActive    *bool    `json:"is_active,omitempty" example:"true"`
}

// ProductListQuery represents query parameters for listing products
type ProductListQuery struct {
	Page     int    `form:"page,default=1" binding:"min=1" example:"1"`
	Limit    int    `form:"limit,default=10" binding:"min=1,max=100" example:"10"`
	Category string `form:"category" example:"Electronics"`
	Search   string `form:"search" example:"iPhone"`
	SortBy   string `form:"sort_by,default=created_at" example:"name"`
	SortDir  string `form:"sort_dir,default=desc" binding:"oneof=asc desc" example:"asc"`
}

// ProductListResponse represents the response payload for listing products
type ProductListResponse struct {
	Products   []Product          `json:"products"`
	Pagination PaginationMetadata `json:"pagination"`
}

// NewProduct creates a new product with generated ID and timestamps
func NewProduct(req ProductCreateRequest, createdBy string) *Product {
	now := time.Now()
	return &Product{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		SKU:         req.SKU,
		Stock:       req.Stock,
		IsActive:    true, // Default active status
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update applies updates from ProductUpdateRequest to the product
func (p *Product) Update(req ProductUpdateRequest) {
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Description != nil {
		p.Description = *req.Description
	}
	if req.Price != nil {
		p.Price = *req.Price
	}
	if req.Category != nil {
		p.Category = *req.Category
	}
	if req.SKU != nil {
		p.SKU = *req.SKU
	}
	if req.Stock != nil {
		p.Stock = *req.Stock
	}
	if req.IsActive != nil {
		p.IsActive = *req.IsActive
	}
	p.UpdatedAt = time.Now()
}
