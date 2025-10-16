package models

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string      `json:"error" example:"Invalid input"`
	Message string      `json:"message" example:"Validation failed"`
	Details interface{} `json:"details,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginationMetadata represents pagination information
type PaginationMetadata struct {
	Page       int   `json:"page" example:"1"`
	Limit      int   `json:"limit" example:"10"`
	Total      int64 `json:"total" example:"100"`
	TotalPages int   `json:"total_pages" example:"10"`
	HasNext    bool  `json:"has_next" example:"true"`
	HasPrev    bool  `json:"has_prev" example:"false"`
}

// NewPaginationMetadata creates pagination metadata
func NewPaginationMetadata(page, limit int, total int64) PaginationMetadata {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return PaginationMetadata{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"must be a valid email address"`
	Value   string `json:"value,omitempty" example:"invalid-email"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError
