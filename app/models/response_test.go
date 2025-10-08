package models

import (
	"testing"
)

func TestNewSuccessResponse(t *testing.T) {
	data := map[string]string{"key": "value"}
	message := "Success"

	response := NewSuccessResponse(data, message)

	if !response.Success {
		t.Error("Expected success to be true")
	}
	if response.Data == nil {
		t.Error("Expected data to be set")
	}
	if response.Message != message {
		t.Errorf("Expected message '%s', got '%s'", message, response.Message)
	}
	if response.Metadata == nil {
		t.Error("Expected metadata to be set")
	}
	if response.Error != nil {
		t.Error("Expected error to be nil")
	}
}

func TestNewErrorResponse(t *testing.T) {
	code := "TEST_ERROR"
	message := "Test error message"
	details := map[string]string{"field": "test"}

	response := NewErrorResponse(code, message, details)

	if response.Success {
		t.Error("Expected success to be false")
	}
	if response.Error == nil {
		t.Fatal("Expected error to be set")
	}
	if response.Error.Code != code {
		t.Errorf("Expected error code '%s', got '%s'", code, response.Error.Code)
	}
	if response.Error.Message != message {
		t.Errorf("Expected error message '%s', got '%s'", message, response.Error.Message)
	}
	if response.Metadata == nil {
		t.Error("Expected metadata to be set")
	}
}

func TestNewListResponse(t *testing.T) {
	data := []string{"item1", "item2"}
	pagination := &PaginationMeta{
		Page:       1,
		PageSize:   10,
		TotalCount: 100,
		TotalPages: 10,
		HasNext:    true,
		HasPrev:    false,
	}

	response := NewListResponse(data, pagination)

	if !response.Success {
		t.Error("Expected success to be true")
	}
	if response.Data == nil {
		t.Error("Expected data to be set")
	}
	if response.Pagination == nil {
		t.Error("Expected pagination to be set")
	}
	if response.Pagination.Page != 1 {
		t.Errorf("Expected page 1, got %d", response.Pagination.Page)
	}
}

func TestNewCreatedResponse(t *testing.T) {
	data := map[string]string{"id": "123"}

	response := NewCreatedResponse(data)

	if !response.Success {
		t.Error("Expected success to be true")
	}
	if response.Data == nil {
		t.Error("Expected data to be set")
	}
	if response.Message != "Resource created successfully" {
		t.Errorf("Expected created message, got '%s'", response.Message)
	}
}

func TestNewUpdatedResponse(t *testing.T) {
	data := map[string]string{"id": "123"}

	response := NewUpdatedResponse(data)

	if !response.Success {
		t.Error("Expected success to be true")
	}
	if response.Data == nil {
		t.Error("Expected data to be set")
	}
	if response.Message != "Resource updated successfully" {
		t.Errorf("Expected updated message, got '%s'", response.Message)
	}
}

func TestNewDeletedResponse(t *testing.T) {
	response := NewDeletedResponse()

	if !response.Success {
		t.Error("Expected success to be true")
	}
	if response.Data != nil {
		t.Error("Expected data to be nil")
	}
	if response.Message != "Resource deleted successfully" {
		t.Errorf("Expected deleted message, got '%s'", response.Message)
	}
}

func TestNewNoContentResponse(t *testing.T) {
	response := NewNoContentResponse()

	if !response.Success {
		t.Error("Expected success to be true")
	}
	if response.Data != nil {
		t.Error("Expected data to be nil")
	}
	if response.Message != "No content" {
		t.Errorf("Expected no content message, got '%s'", response.Message)
	}
}
