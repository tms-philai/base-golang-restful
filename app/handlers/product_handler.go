package handlers

import (
	"net/http"
	"strconv"

	"base-golang-restful-app/middleware"
	"base-golang-restful-app/models"
	"base-golang-restful-app/services"

	"github.com/gin-gonic/gin"
)

// ProductHandler handles product-related requests
type ProductHandler struct {
	productService *services.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct handles product creation
// @Summary Create a new product
// @Description Create a new product (authenticated users)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ProductCreateRequest true "Product creation request"
// @Success 201 {object} models.Product
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req models.ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid request payload",
			Details: err.Error(),
		})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists || userID == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	product, err := h.productService.Create(req, userID)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "SKU already exists" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, models.ErrorResponse{
			Error:   "product_creation_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProduct handles getting a product by ID
// @Summary Get product by ID
// @Description Get a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Product ID is required",
		})
		return
	}

	product, err := h.productService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "product_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct handles updating a product
// @Summary Update product
// @Description Update product information (owner or admin)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param request body models.ProductUpdateRequest true "Product update request"
// @Success 200 {object} models.Product
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Product ID is required",
		})
		return
	}

	// Get current user
	currentUserInterface, exists := middleware.GetCurrentUser(c)
	if !exists || currentUserInterface == nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	currentUser, ok := currentUserInterface.(*models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Check if product exists and user has permission to update
	product, err := h.productService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "product_not_found",
			Message: err.Error(),
		})
		return
	}

	// Allow admin to update any product, or owner to update their own product
	if !currentUser.HasRole("admin") && currentUser.ID.String() != product.CreatedBy {
		c.JSON(http.StatusForbidden, models.ErrorResponse{
			Error:   "forbidden",
			Message: "You can only update your own products",
		})
		return
	}

	var req models.ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid request payload",
			Details: err.Error(),
		})
		return
	}

	updatedProduct, err := h.productService.Update(id, req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "product not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "SKU already exists" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, models.ErrorResponse{
			Error:   "product_update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}

// DeleteProduct handles product deletion (soft delete)
// @Summary Delete product
// @Description Delete a product (soft delete - owner or admin)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Product ID is required",
		})
		return
	}

	// Get current user
	currentUserInterface, exists := middleware.GetCurrentUser(c)
	if !exists || currentUserInterface == nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	currentUser, ok := currentUserInterface.(*models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Check if product exists and user has permission to delete
	product, err := h.productService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "product_not_found",
			Message: err.Error(),
		})
		return
	}

	// Allow admin to delete any product, or owner to delete their own product
	if !currentUser.HasRole("admin") && currentUser.ID.String() != product.CreatedBy {
		c.JSON(http.StatusForbidden, models.ErrorResponse{
			Error:   "forbidden",
			Message: "You can only delete your own products",
		})
		return
	}

	err = h.productService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "product_deletion_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Product deleted successfully",
	})
}

// ListProducts handles listing products with filtering, searching, and pagination
// @Summary List products
// @Description Get a paginated list of products with filtering and searching
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param category query string false "Filter by category"
// @Param search query string false "Search in name, description, or SKU"
// @Param sort_by query string false "Sort by field" default(created_at)
// @Param sort_dir query string false "Sort direction (asc/desc)" default(desc)
// @Success 200 {object} models.ProductListResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/v1/products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	// Parse query parameters
	query := models.ProductListQuery{
		Page:     1,
		Limit:    10,
		SortBy:   "created_at",
		SortDir:  "desc",
		Category: c.Query("category"),
		Search:   c.Query("search"),
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			query.Page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			query.Limit = l
		}
	}

	if sortBy := c.Query("sort_by"); sortBy != "" {
		query.SortBy = sortBy
	}

	if sortDir := c.Query("sort_dir"); sortDir == "asc" || sortDir == "desc" {
		query.SortDir = sortDir
	}

	products, pagination, err := h.productService.List(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "list_products_failed",
			Message: err.Error(),
		})
		return
	}

	response := models.ProductListResponse{
		Products:   products,
		Pagination: pagination,
	}

	c.JSON(http.StatusOK, response)
}

// GetCategories handles getting all product categories
// @Summary Get product categories
// @Description Get all unique product categories
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} object{categories=[]string}
// @Router /api/v1/products/categories [get]
func (h *ProductHandler) GetCategories(c *gin.Context) {
	categories := h.productService.GetCategories()
	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

// UpdateStock handles updating product stock
// @Summary Update product stock
// @Description Update the stock quantity of a product (owner or admin)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param request body object{stock=int} true "Stock update request"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/products/{id}/stock [patch]
func (h *ProductHandler) UpdateStock(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Product ID is required",
		})
		return
	}

	// Get current user
	currentUserInterface, exists := middleware.GetCurrentUser(c)
	if !exists || currentUserInterface == nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	currentUser, ok := currentUserInterface.(*models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Check if product exists and user has permission to update
	product, err := h.productService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "product_not_found",
			Message: err.Error(),
		})
		return
	}

	// Allow admin to update any product, or owner to update their own product
	if !currentUser.HasRole("admin") && currentUser.ID.String() != product.CreatedBy {
		c.JSON(http.StatusForbidden, models.ErrorResponse{
			Error:   "forbidden",
			Message: "You can only update your own products",
		})
		return
	}

	var req struct {
		Stock int `json:"stock" binding:"min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid request payload",
			Details: err.Error(),
		})
		return
	}

	err = h.productService.UpdateStock(id, req.Stock)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "stock_update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Stock updated successfully",
	})
}
