package models

import (
	"testing"
)

func TestNewPaginationRequest(t *testing.T) {
	req := NewPaginationRequest()

	if req.Page != 1 {
		t.Errorf("Expected default page to be 1, got %d", req.Page)
	}
	if req.PageSize != 10 {
		t.Errorf("Expected default pageSize to be 10, got %d", req.PageSize)
	}
	if req.Order != "asc" {
		t.Errorf("Expected default order to be 'asc', got '%s'", req.Order)
	}
}

func TestPaginationRequest_Validate(t *testing.T) {
	tests := []struct {
		name     string
		input    *PaginationRequest
		expected *PaginationRequest
	}{
		{
			name: "Valid request",
			input: &PaginationRequest{
				Page:     2,
				PageSize: 20,
				Order:    "desc",
			},
			expected: &PaginationRequest{
				Page:     2,
				PageSize: 20,
				Order:    "desc",
			},
		},
		{
			name: "Invalid page (< 1)",
			input: &PaginationRequest{
				Page:     0,
				PageSize: 10,
			},
			expected: &PaginationRequest{
				Page:     1,
				PageSize: 10,
				Order:    "asc",
			},
		},
		{
			name: "Invalid pageSize (< 1)",
			input: &PaginationRequest{
				Page:     1,
				PageSize: 0,
			},
			expected: &PaginationRequest{
				Page:     1,
				PageSize: 10,
				Order:    "asc",
			},
		},
		{
			name: "PageSize exceeds max (> 100)",
			input: &PaginationRequest{
				Page:     1,
				PageSize: 150,
			},
			expected: &PaginationRequest{
				Page:     1,
				PageSize: 100,
				Order:    "asc",
			},
		},
		{
			name: "Empty order",
			input: &PaginationRequest{
				Page:     1,
				PageSize: 10,
				Order:    "",
			},
			expected: &PaginationRequest{
				Page:     1,
				PageSize: 10,
				Order:    "asc",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.Validate()

			if tt.input.Page != tt.expected.Page {
				t.Errorf("Expected page %d, got %d", tt.expected.Page, tt.input.Page)
			}
			if tt.input.PageSize != tt.expected.PageSize {
				t.Errorf("Expected pageSize %d, got %d", tt.expected.PageSize, tt.input.PageSize)
			}
			if tt.input.Order != tt.expected.Order {
				t.Errorf("Expected order '%s', got '%s'", tt.expected.Order, tt.input.Order)
			}
		})
	}
}

func TestPaginationRequest_GetOffset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		expected int
	}{
		{"Page 1", 1, 10, 0},
		{"Page 2", 2, 10, 10},
		{"Page 3", 3, 20, 40},
		{"Page 5", 5, 15, 60},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &PaginationRequest{
				Page:     tt.page,
				PageSize: tt.pageSize,
			}

			offset := req.GetOffset()
			if offset != tt.expected {
				t.Errorf("Expected offset %d, got %d", tt.expected, offset)
			}
		})
	}
}

func TestPaginationRequest_GetLimit(t *testing.T) {
	req := &PaginationRequest{
		Page:     1,
		PageSize: 25,
	}

	limit := req.GetLimit()
	if limit != 25 {
		t.Errorf("Expected limit 25, got %d", limit)
	}
}

func TestNewCursorPaginationRequest(t *testing.T) {
	req := NewCursorPaginationRequest()

	if req.PageSize != 10 {
		t.Errorf("Expected default pageSize to be 10, got %d", req.PageSize)
	}
	if req.Order != "asc" {
		t.Errorf("Expected default order to be 'asc', got '%s'", req.Order)
	}
}

func TestCursorPaginationRequest_Validate(t *testing.T) {
	tests := []struct {
		name     string
		input    *CursorPaginationRequest
		expected *CursorPaginationRequest
	}{
		{
			name: "Valid request",
			input: &CursorPaginationRequest{
				Cursor:   "abc123",
				PageSize: 20,
				Order:    "desc",
			},
			expected: &CursorPaginationRequest{
				Cursor:   "abc123",
				PageSize: 20,
				Order:    "desc",
			},
		},
		{
			name: "Invalid pageSize (< 1)",
			input: &CursorPaginationRequest{
				PageSize: 0,
			},
			expected: &CursorPaginationRequest{
				PageSize: 10,
				Order:    "asc",
			},
		},
		{
			name: "PageSize exceeds max (> 100)",
			input: &CursorPaginationRequest{
				PageSize: 200,
			},
			expected: &CursorPaginationRequest{
				PageSize: 100,
				Order:    "asc",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.Validate()

			if tt.input.PageSize != tt.expected.PageSize {
				t.Errorf("Expected pageSize %d, got %d", tt.expected.PageSize, tt.input.PageSize)
			}
			if tt.input.Order != tt.expected.Order {
				t.Errorf("Expected order '%s', got '%s'", tt.expected.Order, tt.input.Order)
			}
		})
	}
}
