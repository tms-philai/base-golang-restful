package unit

import (
	"base-gin/internal/app/handlers"
	"base-gin/internal/domain/models"
	"base-gin/internal/domain/services"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFileService is a mock implementation of FileService
type MockFileService struct {
	mock.Mock
}

func (m *MockFileService) UploadFile(ctx context.Context, input services.UploadFileInput) (*services.UploadFileResult, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*services.UploadFileResult), args.Error(1)
}

func (m *MockFileService) UploadMultipleFiles(ctx context.Context, fileHeaders []*multipart.FileHeader, uploadedBy *uuid.UUID, isPublic bool, metadata string) ([]*services.UploadFileResult, error) {
	args := m.Called(ctx, fileHeaders, uploadedBy, isPublic, metadata)
	return args.Get(0).([]*services.UploadFileResult), args.Error(1)
}

func (m *MockFileService) GetFile(ctx context.Context, fileID uuid.UUID) (*models.File, error) {
	args := m.Called(ctx, fileID)
	return args.Get(0).(*models.File), args.Error(1)
}

func (m *MockFileService) DeleteFile(ctx context.Context, fileID uuid.UUID) error {
	args := m.Called(ctx, fileID)
	return args.Error(0)
}

func (m *MockFileService) GetFilesByUser(ctx context.Context, userID uuid.UUID) ([]models.File, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.File), args.Error(1)
}

func (m *MockFileService) GetFileContent(ctx context.Context, fileID uuid.UUID) ([]byte, error) {
	args := m.Called(ctx, fileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockFileService) GetFilePath(ctx context.Context, fileID uuid.UUID) (string, error) {
	args := m.Called(ctx, fileID)
	return args.String(0), args.Error(1)
}

func TestFileHandler_UploadFile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		fileName       string
		fileContent    string
		isPublic       bool
		metadata       string
		mockSetup      func(*MockFileService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:        "successful file upload",
			userID:      "user-123",
			fileName:    "test.txt",
			fileContent: "Hello, World!",
			isPublic:    false,
			metadata:    "test metadata",
			mockSetup: func(fileService *MockFileService) {
				file := &models.File{
					ID:       uuid.New(),
					MimeType: "text/plain",
					IsPublic: false,
					Metadata: "test metadata",
				}
				result := &services.UploadFileResult{
					File: file,
					URL:  "http://example.com/files/" + file.ID.String(),
				}
				fileService.On("UploadFile", mock.Anything, mock.AnythingOfType("services.UploadFileInput")).Return(result, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:        "successful public file upload",
			userID:      "user-123",
			fileName:    "public.txt",
			fileContent: "Public content",
			isPublic:    true,
			metadata:    "",
			mockSetup: func(fileService *MockFileService) {
				file := &models.File{
					ID:       uuid.New(),
					FileName: "public.txt",
					FileSize: 14,
					MimeType: "text/plain",
					IsPublic: true,
					Metadata: "",
				}
				result := &services.UploadFileResult{
					File: file,
					URL:  "http://example.com/files/" + file.ID.String(),
				}
				fileService.On("UploadFile", mock.Anything, mock.AnythingOfType("services.UploadFileInput")).Return(result, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "unauthorized - missing user ID",
			userID:         "",
			fileName:       "test.txt",
			fileContent:    "Hello, World!",
			isPublic:       false,
			metadata:       "",
			mockSetup:      func(*MockFileService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "unauthorized",
		},
		{
			name:        "file upload service error",
			userID:      "user-123",
			fileName:    "test.txt",
			fileContent: "Hello, World!",
			isPublic:    false,
			metadata:    "",
			mockSetup: func(fileService *MockFileService) {
				fileService.On("UploadFile", mock.Anything, mock.AnythingOfType("services.UploadFileInput")).Return((*services.UploadFileResult)(nil), errors.New("upload failed"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "upload_failed",
		},
		{
			name:           "missing file",
			userID:         "user-123",
			fileName:       "",
			fileContent:    "",
			isPublic:       false,
			metadata:       "",
			mockSetup:      func(*MockFileService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "file_required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockFileService := new(MockFileService)
			tt.mockSetup(mockFileService)

			// Create handler
			handler := handlers.NewFileHandler(mockFileService)

			// Setup router with middleware mock
			router := gin.New()
			router.POST("/files/upload", func(c *gin.Context) {
				// Mock middleware setting user ID
				if tt.userID != "" {
					c.Set("user_id", tt.userID)
				}
				handler.UploadFile(c)
			})

			// Create multipart form data
			var body bytes.Buffer
			writer := multipart.NewWriter(&body)

			// Add file if provided
			if tt.fileName != "" {
				fileWriter, err := writer.CreateFormFile("file", tt.fileName)
				assert.NoError(t, err)
				fileWriter.Write([]byte(tt.fileContent))
			}

			// Add form fields
			writer.WriteField("is_public", "false")
			if tt.isPublic {
				writer.WriteField("is_public", "true")
			}
			writer.WriteField("metadata", tt.metadata)

			writer.Close()

			// Create request
			req, _ := http.NewRequest("POST", "/files/upload", &body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response.Error)
			}

			// Verify mock expectations
			mockFileService.AssertExpectations(t)
		})
	}
}

func TestFileHandler_GetFile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		fileID         string
		mockSetup      func(*MockFileService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successful get file",
			fileID: "file-123",
			mockSetup: func(fileService *MockFileService) {
				file := &models.File{
					ID:       uuid.New(),
					FileName: "test.txt",
					FileSize: 13,
					MimeType: "text/plain",
					IsPublic: true,
				}
				fileService.On("GetFile", "file-123").Return(file, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "file not found",
			fileID: "nonexistent",
			mockSetup: func(fileService *MockFileService) {
				fileService.On("GetFile", "nonexistent").Return((*models.File)(nil), errors.New("file not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "file_not_found",
		},
		{
			name:           "missing file ID",
			fileID:         "",
			mockSetup:      func(*MockFileService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockFileService := new(MockFileService)
			tt.mockSetup(mockFileService)

			// Create handler
			handler := handlers.NewFileHandler(mockFileService)

			// Setup router
			router := gin.New()
			router.GET("/files/:id", handler.GetFile)

			// Create request
			url := "/files/" + tt.fileID
			req, _ := http.NewRequest("GET", url, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response.Error)
			}

			// Verify mock expectations
			mockFileService.AssertExpectations(t)
		})
	}
}

func TestFileHandler_GetFileURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		fileID         string
		mockSetup      func(*MockFileService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successful get file URL",
			fileID: "file-123",
			mockSetup: func(fileService *MockFileService) {
				fileService.On("GetFileURL", "file-123").Return("http://example.com/files/file-123", nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "file not found",
			fileID: "nonexistent",
			mockSetup: func(fileService *MockFileService) {
				fileService.On("GetFileURL", "nonexistent").Return("", errors.New("file not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "file_not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockFileService := new(MockFileService)
			tt.mockSetup(mockFileService)

			// Create handler
			handler := handlers.NewFileHandler(mockFileService)

			// Setup router
			router := gin.New()
			router.GET("/files/:id/url", handler.GetFileURL)

			// Create request
			url := "/files/" + tt.fileID + "/url"
			req, _ := http.NewRequest("GET", url, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response.Error)
			}

			// Verify mock expectations
			mockFileService.AssertExpectations(t)
		})
	}
}

func TestFileHandler_ListFiles(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		queryParams    string
		mockSetup      func(*MockFileService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:        "successful list files",
			userID:      "user-123",
			queryParams: "?page=1&limit=10",
			mockSetup: func(fileService *MockFileService) {
				files := []models.File{
					{
						ID:       uuid.New(),
						FileName: "file1.txt",
						FileSize: 100,
						MimeType: "text/plain",
						IsPublic: false,
					},
					{
						ID:       uuid.New(),
						FileName: "file2.jpg",
						FileSize: 200,
						MimeType: "image/jpeg",
						IsPublic: true,
					},
				}
				pagination := models.PaginationMetadata{
					Page:       1,
					Limit:      10,
					Total:      2,
					TotalPages: 1,
				}
				fileService.On("ListFiles", "user-123", mock.AnythingOfType("models.FileListQuery")).Return(files, pagination, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "list files service error",
			userID:      "user-123",
			queryParams: "?page=1&limit=10",
			mockSetup: func(fileService *MockFileService) {
				fileService.On("ListFiles", "user-123", mock.AnythingOfType("models.FileListQuery")).Return(([]models.File)(nil), models.PaginationMetadata{}, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "list_files_failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockFileService := new(MockFileService)
			tt.mockSetup(mockFileService)

			// Create handler
			handler := handlers.NewFileHandler(mockFileService)

			// Setup router with middleware mock
			router := gin.New()
			router.GET("/files", func(c *gin.Context) {
				// Mock middleware setting user ID
				c.Set("user_id", tt.userID)
				handler.ListFiles(c)
			})

			// Create request
			url := "/files" + tt.queryParams
			req, _ := http.NewRequest("GET", url, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response.Error)
			}

			// Verify mock expectations
			mockFileService.AssertExpectations(t)
		})
	}
}

func TestFileHandler_DeleteFile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		fileID         string
		userID         string
		mockSetup      func(*MockFileService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successful delete file",
			fileID: "file-123",
			userID: "user-123",
			mockSetup: func(fileService *MockFileService) {
				fileService.On("DeleteFile", "file-123", "user-123").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "file not found",
			fileID: "nonexistent",
			userID: "user-123",
			mockSetup: func(fileService *MockFileService) {
				fileService.On("DeleteFile", "nonexistent", "user-123").Return(errors.New("file not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "file_not_found",
		},
		{
			name:           "missing file ID",
			fileID:         "",
			userID:         "user-123",
			mockSetup:      func(*MockFileService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockFileService := new(MockFileService)
			tt.mockSetup(mockFileService)

			// Create handler
			handler := handlers.NewFileHandler(mockFileService)

			// Setup router with middleware mock
			router := gin.New()
			router.DELETE("/files/:id", func(c *gin.Context) {
				// Mock middleware setting user ID
				c.Set("user_id", tt.userID)
				handler.DeleteFile(c)
			})

			// Create request
			url := "/files/" + tt.fileID
			req, _ := http.NewRequest("DELETE", url, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response.Error)
			}

			// Verify mock expectations
			mockFileService.AssertExpectations(t)
		})
	}
}
