package utils

import (
	"testing"
)

func TestNewPaginationHelper(t *testing.T) {
	baseURL := "http://example.com/api/users"
	helper := NewPaginationHelper(baseURL)

	if helper == nil {
		t.Error("Expected NewPaginationHelper to return non-nil helper")
	}
	if helper.baseURL != baseURL {
		t.Errorf("Expected baseURL to be '%s', got '%s'", baseURL, helper.baseURL)
	}
}

func TestPaginationHelper_CreatePaginationResponse(t *testing.T) {
	helper := NewPaginationHelper("http://example.com/api/users")

	data := []string{"item1", "item2", "item3"}
	response := helper.CreatePaginationResponse(data, 1, 10, 25, false)

	if response.Pagination.Page != 1 {
		t.Errorf("Expected page to be 1, got %d", response.Pagination.Page)
	}
	if response.Pagination.PageSize != 10 {
		t.Errorf("Expected pageSize to be 10, got %d", response.Pagination.PageSize)
	}
	if response.Pagination.TotalCount != 25 {
		t.Errorf("Expected totalCount to be 25, got %d", response.Pagination.TotalCount)
	}
	if response.Pagination.TotalPages != 3 {
		t.Errorf("Expected totalPages to be 3, got %d", response.Pagination.TotalPages)
	}
	if !response.Pagination.HasNext {
		t.Error("Expected hasNext to be true")
	}
	if response.Pagination.HasPrev {
		t.Error("Expected hasPrev to be false")
	}
}

func TestPaginationHelper_CreatePaginationResponse_WithLinks(t *testing.T) {
	helper := NewPaginationHelper("http://example.com/api/users")

	data := []string{"item1", "item2"}
	response := helper.CreatePaginationResponse(data, 2, 10, 25, true)

	if response.Links == nil {
		t.Fatal("Expected links to be generated")
	}

	if response.Links.First != "http://example.com/api/users?page=1&pageSize=10" {
		t.Errorf("Unexpected first link: %s", response.Links.First)
	}
	if response.Links.Last != "http://example.com/api/users?page=3&pageSize=10" {
		t.Errorf("Unexpected last link: %s", response.Links.Last)
	}
	if response.Links.Next == nil || *response.Links.Next != "http://example.com/api/users?page=3&pageSize=10" {
		t.Error("Expected next link to be set")
	}
	if response.Links.Prev == nil || *response.Links.Prev != "http://example.com/api/users?page=1&pageSize=10" {
		t.Error("Expected prev link to be set")
	}
}

func TestCalculateTotalPages(t *testing.T) {
	tests := []struct {
		name       string
		totalCount int64
		pageSize   int
		expected   int
	}{
		{"Exact division", 100, 10, 10},
		{"With remainder", 105, 10, 11},
		{"Less than pageSize", 5, 10, 1},
		{"Zero count", 0, 10, 0},
		{"Invalid pageSize", 100, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateTotalPages(tt.totalCount, tt.pageSize)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestValidatePaginationParams(t *testing.T) {
	tests := []struct {
		name             string
		inputPage        int
		inputPageSize    int
		expectedPage     int
		expectedPageSize int
	}{
		{"Valid params", 2, 20, 2, 20},
		{"Invalid page (< 1)", 0, 10, 1, 10},
		{"Invalid pageSize (< 1)", 1, 0, 1, 10},
		{"PageSize exceeds max", 1, 150, 1, 100},
		{"Negative values", -1, -5, 1, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page, pageSize := ValidatePaginationParams(tt.inputPage, tt.inputPageSize)

			if page != tt.expectedPage {
				t.Errorf("Expected page %d, got %d", tt.expectedPage, page)
			}
			if pageSize != tt.expectedPageSize {
				t.Errorf("Expected pageSize %d, got %d", tt.expectedPageSize, pageSize)
			}
		})
	}
}
