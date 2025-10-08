package utils

import (
	"base-golang-restful-app/models"
	"testing"

	"gorm.io/gorm"
)

func TestNewQueryBuilder(t *testing.T) {
	var db *gorm.DB
	qb := NewQueryBuilder(db)

	if qb == nil {
		t.Error("Expected NewQueryBuilder to return non-nil query builder")
	}
}

func TestParseSortFields(t *testing.T) {
	tests := []struct {
		name     string
		sortBy   string
		order    string
		expected []SortField
	}{
		{
			name:   "Single field",
			sortBy: "name",
			order:  "asc",
			expected: []SortField{
				{Field: "name", Order: "asc"},
			},
		},
		{
			name:   "Multiple fields with orders",
			sortBy: "name,created_at",
			order:  "asc,desc",
			expected: []SortField{
				{Field: "name", Order: "asc"},
				{Field: "created_at", Order: "desc"},
			},
		},
		{
			name:   "Multiple fields, single order",
			sortBy: "name,email",
			order:  "desc",
			expected: []SortField{
				{Field: "name", Order: "desc"},
				{Field: "email", Order: "asc"},
			},
		},
		{
			name:     "Empty sortBy",
			sortBy:   "",
			order:    "asc",
			expected: nil,
		},
		{
			name:   "Fields with spaces",
			sortBy: " name , email ",
			order:  "asc,desc",
			expected: []SortField{
				{Field: "name", Order: "asc"},
				{Field: "email", Order: "desc"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseSortFields(tt.sortBy, tt.order)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d sort fields, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i].Field != expected.Field {
					t.Errorf("Expected field '%s', got '%s'", expected.Field, result[i].Field)
				}
				if result[i].Order != expected.Order {
					t.Errorf("Expected order '%s', got '%s'", expected.Order, result[i].Order)
				}
			}
		})
	}
}

func TestPaginationRequest_Methods(t *testing.T) {
	req := &models.PaginationRequest{
		Page:     3,
		PageSize: 20,
		SortBy:   "name",
		Order:    "desc",
	}

	offset := req.GetOffset()
	if offset != 40 {
		t.Errorf("Expected offset 40, got %d", offset)
	}

	limit := req.GetLimit()
	if limit != 20 {
		t.Errorf("Expected limit 20, got %d", limit)
	}
}
