package models

import (
	"testing"
)

func TestErrorCodes(t *testing.T) {
	codes := map[string]string{
		ErrCodeValidation:              "VALIDATION_ERROR",
		ErrCodeInvalidInput:            "INVALID_INPUT",
		ErrCodeInvalidFormat:           "INVALID_FORMAT",
		ErrCodeUnauthorized:            "UNAUTHORIZED",
		ErrCodeInvalidToken:            "INVALID_TOKEN",
		ErrCodeTokenExpired:            "TOKEN_EXPIRED",
		ErrCodeForbidden:               "FORBIDDEN",
		ErrCodeInsufficientPermissions: "INSUFFICIENT_PERMISSIONS",
		ErrCodeNotFound:                "NOT_FOUND",
		ErrCodeResourceNotFound:        "RESOURCE_NOT_FOUND",
		ErrCodeConflict:                "CONFLICT",
		ErrCodeDuplicateEntry:          "DUPLICATE_ENTRY",
		ErrCodeInternalServer:          "INTERNAL_SERVER_ERROR",
		ErrCodeDatabaseError:           "DATABASE_ERROR",
		ErrCodeExternalService:         "EXTERNAL_SERVICE_ERROR",
		ErrCodeBadRequest:              "BAD_REQUEST",
		ErrCodeServiceUnavailable:      "SERVICE_UNAVAILABLE",
	}

	for code, expected := range codes {
		if code != expected {
			t.Errorf("Expected error code '%s', got '%s'", expected, code)
		}
	}
}

func TestNewValidationErrorResponse(t *testing.T) {
	errors := ValidationErrors{
		{Field: "email", Message: "Invalid email"},
		{Field: "password", Message: "Password too short"},
	}

	response := NewValidationErrorResponse(errors)

	if response.Success {
		t.Error("Expected success to be false")
	}
	if response.Error == nil {
		t.Fatal("Expected error to be set")
	}
	if response.Error.Code != ErrCodeValidation {
		t.Errorf("Expected error code '%s', got '%s'", ErrCodeValidation, response.Error.Code)
	}
	if response.Error.Message != "Validation failed" {
		t.Errorf("Expected validation message, got '%s'", response.Error.Message)
	}
}

func TestNewUnauthorizedErrorResponse(t *testing.T) {
	response := NewUnauthorizedErrorResponse()

	if response.Success {
		t.Error("Expected success to be false")
	}
	if response.Error == nil {
		t.Fatal("Expected error to be set")
	}
	if response.Error.Code != ErrCodeUnauthorized {
		t.Errorf("Expected error code '%s', got '%s'", ErrCodeUnauthorized, response.Error.Code)
	}
}

func TestNewForbiddenErrorResponse(t *testing.T) {
	response := NewForbiddenErrorResponse()

	if response.Success {
		t.Error("Expected success to be false")
	}
	if response.Error == nil {
		t.Fatal("Expected error to be set")
	}
	if response.Error.Code != ErrCodeForbidden {
		t.Errorf("Expected error code '%s', got '%s'", ErrCodeForbidden, response.Error.Code)
	}
}

func TestNewNotFoundErrorResponse(t *testing.T) {
	resource := "User"
	response := NewNotFoundErrorResponse(resource)

	if response.Success {
		t.Error("Expected success to be false")
	}
	if response.Error == nil {
		t.Fatal("Expected error to be set")
	}
	if response.Error.Code != ErrCodeNotFound {
		t.Errorf("Expected error code '%s', got '%s'", ErrCodeNotFound, response.Error.Code)
	}
	if response.Error.Message != "User not found" {
		t.Errorf("Expected message 'User not found', got '%s'", response.Error.Message)
	}
}

func TestNewConflictErrorResponse(t *testing.T) {
	resource := "Email"
	response := NewConflictErrorResponse(resource)

	if response.Success {
		t.Error("Expected success to be false")
	}
	if response.Error == nil {
		t.Fatal("Expected error to be set")
	}
	if response.Error.Code != ErrCodeConflict {
		t.Errorf("Expected error code '%s', got '%s'", ErrCodeConflict, response.Error.Code)
	}
	if response.Error.Message != "Email already exists" {
		t.Errorf("Expected message 'Email already exists', got '%s'", response.Error.Message)
	}
}

func TestNewInternalErrorResponse(t *testing.T) {
	response := NewInternalErrorResponse()

	if response.Success {
		t.Error("Expected success to be false")
	}
	if response.Error == nil {
		t.Fatal("Expected error to be set")
	}
	if response.Error.Code != ErrCodeInternalServer {
		t.Errorf("Expected error code '%s', got '%s'", ErrCodeInternalServer, response.Error.Code)
	}
}

func TestNewBadRequestErrorResponse(t *testing.T) {
	message := "Invalid request"
	response := NewBadRequestErrorResponse(message)

	if response.Success {
		t.Error("Expected success to be false")
	}
	if response.Error == nil {
		t.Fatal("Expected error to be set")
	}
	if response.Error.Code != ErrCodeBadRequest {
		t.Errorf("Expected error code '%s', got '%s'", ErrCodeBadRequest, response.Error.Code)
	}
	if response.Error.Message != message {
		t.Errorf("Expected message '%s', got '%s'", message, response.Error.Message)
	}
}

func TestNewServiceUnavailableErrorResponse(t *testing.T) {
	response := NewServiceUnavailableErrorResponse()

	if response.Success {
		t.Error("Expected success to be false")
	}
	if response.Error == nil {
		t.Fatal("Expected error to be set")
	}
	if response.Error.Code != ErrCodeServiceUnavailable {
		t.Errorf("Expected error code '%s', got '%s'", ErrCodeServiceUnavailable, response.Error.Code)
	}
}
