package models

import "time"

const (
	ErrCodeValidation              = "VALIDATION_ERROR"
	ErrCodeInvalidInput            = "INVALID_INPUT"
	ErrCodeInvalidFormat           = "INVALID_FORMAT"
	ErrCodeUnauthorized            = "UNAUTHORIZED"
	ErrCodeInvalidToken            = "INVALID_TOKEN"
	ErrCodeTokenExpired            = "TOKEN_EXPIRED"
	ErrCodeForbidden               = "FORBIDDEN"
	ErrCodeInsufficientPermissions = "INSUFFICIENT_PERMISSIONS"
	ErrCodeNotFound                = "NOT_FOUND"
	ErrCodeResourceNotFound        = "RESOURCE_NOT_FOUND"
	ErrCodeConflict                = "CONFLICT"
	ErrCodeDuplicateEntry          = "DUPLICATE_ENTRY"
	ErrCodeInternalServer          = "INTERNAL_SERVER_ERROR"
	ErrCodeDatabaseError           = "DATABASE_ERROR"
	ErrCodeExternalService         = "EXTERNAL_SERVICE_ERROR"
	ErrCodeBadRequest              = "BAD_REQUEST"
	ErrCodeServiceUnavailable      = "SERVICE_UNAVAILABLE"
)

func NewValidationErrorResponse(errors ValidationErrors) *APIResponse {
	return &APIResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    ErrCodeValidation,
			Message: "Validation failed",
			Details: errors,
		},
		Metadata: &Metadata{
			Timestamp: time.Now(),
		},
	}
}

func NewUnauthorizedErrorResponse() *APIResponse {
	return NewErrorResponse(
		ErrCodeUnauthorized,
		"Authentication required",
		nil,
	)
}

func NewForbiddenErrorResponse() *APIResponse {
	return NewErrorResponse(
		ErrCodeForbidden,
		"Access forbidden",
		nil,
	)
}

func NewNotFoundErrorResponse(resource string) *APIResponse {
	return NewErrorResponse(
		ErrCodeNotFound,
		resource+" not found",
		nil,
	)
}

func NewConflictErrorResponse(resource string) *APIResponse {
	return NewErrorResponse(
		ErrCodeConflict,
		resource+" already exists",
		nil,
	)
}

func NewInternalErrorResponse() *APIResponse {
	return NewErrorResponse(
		ErrCodeInternalServer,
		"Internal server error",
		nil,
	)
}

func NewBadRequestErrorResponse(message string) *APIResponse {
	return NewErrorResponse(
		ErrCodeBadRequest,
		message,
		nil,
	)
}

func NewServiceUnavailableErrorResponse() *APIResponse {
	return NewErrorResponse(
		ErrCodeServiceUnavailable,
		"Service temporarily unavailable",
		nil,
	)
}
