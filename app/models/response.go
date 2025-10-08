package models

import "time"

type APIResponse struct {
	Success  bool        `json:"success"`
	Data     interface{} `json:"data,omitempty"`
	Metadata *Metadata   `json:"metadata,omitempty"`
	Error    *ErrorInfo  `json:"error,omitempty"`
	Message  string      `json:"message,omitempty"`
}

type Metadata struct {
	Timestamp     time.Time `json:"timestamp"`
	Version       string    `json:"version,omitempty"`
	CorrelationID string    `json:"correlation_id,omitempty"`
	RequestID     string    `json:"request_id,omitempty"`
}

type ErrorInfo struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type ListResponse struct {
	Success    bool               `json:"success"`
	Data       interface{}        `json:"data"`
	Pagination *PaginationMeta    `json:"pagination,omitempty"`
	Links      *PaginationLinks   `json:"links,omitempty"`
	Metadata   *Metadata          `json:"metadata,omitempty"`
}

func NewSuccessResponse(data interface{}, message string) *APIResponse {
	return &APIResponse{
		Success: true,
		Data:    data,
		Message: message,
		Metadata: &Metadata{
			Timestamp: time.Now(),
		},
	}
}

func NewErrorResponse(code, message string, details interface{}) *APIResponse {
	return &APIResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
		Metadata: &Metadata{
			Timestamp: time.Now(),
		},
	}
}

func NewListResponse(data interface{}, pagination *PaginationMeta) *ListResponse {
	return &ListResponse{
		Success:    true,
		Data:       data,
		Pagination: pagination,
		Metadata: &Metadata{
			Timestamp: time.Now(),
		},
	}
}

func NewCreatedResponse(data interface{}) *APIResponse {
	return &APIResponse{
		Success: true,
		Data:    data,
		Message: "Resource created successfully",
		Metadata: &Metadata{
			Timestamp: time.Now(),
		},
	}
}

func NewUpdatedResponse(data interface{}) *APIResponse {
	return &APIResponse{
		Success: true,
		Data:    data,
		Message: "Resource updated successfully",
		Metadata: &Metadata{
			Timestamp: time.Now(),
		},
	}
}

func NewDeletedResponse() *APIResponse {
	return &APIResponse{
		Success: true,
		Message: "Resource deleted successfully",
		Metadata: &Metadata{
			Timestamp: time.Now(),
		},
	}
}

func NewNoContentResponse() *APIResponse {
	return &APIResponse{
		Success: true,
		Message: "No content",
		Metadata: &Metadata{
			Timestamp: time.Now(),
		},
	}
}
