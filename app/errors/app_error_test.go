package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		appErr   *AppError
		expected string
	}{
		{
			name: "error without wrapped error",
			appErr: &AppError{
				Code:    "TEST_ERROR",
				Message: "Test error message",
			},
			expected: "TEST_ERROR: Test error message",
		},
		{
			name: "error with wrapped error",
			appErr: &AppError{
				Code:    "TEST_ERROR",
				Message: "Test error message",
				Err:     errors.New("wrapped error"),
			},
			expected: "TEST_ERROR: Test error message (wrapped error)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.appErr.Error()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	wrappedErr := errors.New("wrapped error")
	appErr := &AppError{
		Code:    "TEST_ERROR",
		Message: "Test message",
		Err:     wrappedErr,
	}

	unwrapped := appErr.Unwrap()
	assert.Equal(t, wrappedErr, unwrapped)
}

func TestAppError_WithDetails(t *testing.T) {
	appErr := &AppError{
		Code:    "TEST_ERROR",
		Message: "Test message",
	}

	details := map[string]interface{}{
		"field": "value",
	}

	result := appErr.WithDetails(details)

	assert.Equal(t, details, result.Details)
	assert.Equal(t, appErr, result)
}

func TestAppError_WithError(t *testing.T) {
	appErr := &AppError{
		Code:    "TEST_ERROR",
		Message: "Test message",
	}

	wrappedErr := errors.New("wrapped error")
	result := appErr.WithError(wrappedErr)

	assert.Equal(t, wrappedErr, result.Err)
	assert.Equal(t, appErr, result)
}

func TestNewAppError(t *testing.T) {
	code := "TEST_ERROR"
	message := "Test message"
	statusCode := http.StatusBadRequest

	appErr := NewAppError(code, message, statusCode)

	assert.Equal(t, code, appErr.Code)
	assert.Equal(t, message, appErr.Message)
	assert.Equal(t, statusCode, appErr.StatusCode)
}

func TestErrorConstructors(t *testing.T) {
	tests := []struct {
		name           string
		constructor    func(string) *AppError
		message        string
		expectedCode   string
		expectedStatus int
	}{
		{
			name:           "BadRequest",
			constructor:    BadRequest,
			message:        "Bad request message",
			expectedCode:   "BAD_REQUEST",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Unauthorized",
			constructor:    Unauthorized,
			message:        "Unauthorized message",
			expectedCode:   "UNAUTHORIZED",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Forbidden",
			constructor:    Forbidden,
			message:        "Forbidden message",
			expectedCode:   "FORBIDDEN",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "NotFound",
			constructor:    NotFound,
			message:        "Not found message",
			expectedCode:   "NOT_FOUND",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Conflict",
			constructor:    Conflict,
			message:        "Conflict message",
			expectedCode:   "CONFLICT",
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "Validation",
			constructor:    Validation,
			message:        "Validation message",
			expectedCode:   "VALIDATION_ERROR",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "InternalServer",
			constructor:    InternalServer,
			message:        "Internal server message",
			expectedCode:   "INTERNAL_SERVER_ERROR",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "DatabaseError",
			constructor:    DatabaseError,
			message:        "Database error message",
			expectedCode:   "DATABASE_ERROR",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "ServiceUnavailable",
			constructor:    ServiceUnavailable,
			message:        "Service unavailable message",
			expectedCode:   "SERVICE_UNAVAILABLE",
			expectedStatus: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.constructor(tt.message)

			assert.Equal(t, tt.expectedCode, err.Code)
			assert.Equal(t, tt.message, err.Message)
			assert.Equal(t, tt.expectedStatus, err.StatusCode)
		})
	}
}
