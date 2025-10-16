package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	StatusCode int                    `json:"-"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Err        error                  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

func NewAppError(code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

var (
	ErrBadRequest = &AppError{
		Code:       "BAD_REQUEST",
		Message:    "Bad request",
		StatusCode: http.StatusBadRequest,
	}

	ErrUnauthorized = &AppError{
		Code:       "UNAUTHORIZED",
		Message:    "Unauthorized",
		StatusCode: http.StatusUnauthorized,
	}

	ErrForbidden = &AppError{
		Code:       "FORBIDDEN",
		Message:    "Forbidden",
		StatusCode: http.StatusForbidden,
	}

	ErrNotFound = &AppError{
		Code:       "NOT_FOUND",
		Message:    "Resource not found",
		StatusCode: http.StatusNotFound,
	}

	ErrConflict = &AppError{
		Code:       "CONFLICT",
		Message:    "Resource conflict",
		StatusCode: http.StatusConflict,
	}

	ErrValidation = &AppError{
		Code:       "VALIDATION_ERROR",
		Message:    "Validation error",
		StatusCode: http.StatusUnprocessableEntity,
	}

	ErrInternalServer = &AppError{
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    "Internal server error",
		StatusCode: http.StatusInternalServerError,
	}

	ErrDatabaseError = &AppError{
		Code:       "DATABASE_ERROR",
		Message:    "Database error",
		StatusCode: http.StatusInternalServerError,
	}

	ErrServiceUnavailable = &AppError{
		Code:       "SERVICE_UNAVAILABLE",
		Message:    "Service unavailable",
		StatusCode: http.StatusServiceUnavailable,
	}
)

func BadRequest(message string) *AppError {
	return &AppError{
		Code:       ErrBadRequest.Code,
		Message:    message,
		StatusCode: ErrBadRequest.StatusCode,
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		Code:       ErrUnauthorized.Code,
		Message:    message,
		StatusCode: ErrUnauthorized.StatusCode,
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		Code:       ErrForbidden.Code,
		Message:    message,
		StatusCode: ErrForbidden.StatusCode,
	}
}

func NotFound(message string) *AppError {
	return &AppError{
		Code:       ErrNotFound.Code,
		Message:    message,
		StatusCode: ErrNotFound.StatusCode,
	}
}

func Conflict(message string) *AppError {
	return &AppError{
		Code:       ErrConflict.Code,
		Message:    message,
		StatusCode: ErrConflict.StatusCode,
	}
}

func Validation(message string) *AppError {
	return &AppError{
		Code:       ErrValidation.Code,
		Message:    message,
		StatusCode: ErrValidation.StatusCode,
	}
}

func InternalServer(message string) *AppError {
	return &AppError{
		Code:       ErrInternalServer.Code,
		Message:    message,
		StatusCode: ErrInternalServer.StatusCode,
	}
}

func DatabaseError(message string) *AppError {
	return &AppError{
		Code:       ErrDatabaseError.Code,
		Message:    message,
		StatusCode: ErrDatabaseError.StatusCode,
	}
}

func ServiceUnavailable(message string) *AppError {
	return &AppError{
		Code:       ErrServiceUnavailable.Code,
		Message:    message,
		StatusCode: ErrServiceUnavailable.StatusCode,
	}
}
