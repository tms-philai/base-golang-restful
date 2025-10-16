package handlers

import (
	"base-gin/internal/app/middleware"
	"base-gin/internal/domain/models"
	"base-gin/internal/domain/services"
	"base-gin/internal/pkg/validate"
	"context"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// FileServiceInterface defines the interface for FileService
type FileServiceInterface interface {
	UploadFile(ctx context.Context, input services.UploadFileInput) (*services.UploadFileResult, error)
	UploadMultipleFiles(ctx context.Context, fileHeaders []*multipart.FileHeader, uploadedBy *uuid.UUID, isPublic bool, metadata string) ([]*services.UploadFileResult, error)
	GetFile(ctx context.Context, fileID uuid.UUID) (*models.File, error)
	DeleteFile(ctx context.Context, fileID uuid.UUID) error
	GetFilesByUser(ctx context.Context, userID uuid.UUID) ([]models.File, error)
	GetFileContent(ctx context.Context, fileID uuid.UUID) ([]byte, error)
	GetFilePath(ctx context.Context, fileID uuid.UUID) (string, error)
}

type FileHandler struct {
	fileService FileServiceInterface
}

func NewFileHandler(fileService FileServiceInterface) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

// UploadFile godoc
// @Summary Upload a single file
// @Description Upload a single file with optional metadata
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File to upload"
// @Param is_public formData bool false "Whether the file is public (default: false)"
// @Param metadata formData string false "Additional metadata for the file"
// @Success 201 {object} models.FileUploadResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 413 {object} models.ErrorResponse
// @Router /api/v1/files/upload [post]
func (h *FileHandler) UploadFile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists || userID == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_user_id",
			Message: "Invalid user ID format",
		})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "file_required",
			Message: "File is required",
			Details: err.Error(),
		})
		return
	}

	isPublic := c.PostForm("is_public") == "true"
	metadata := c.PostForm("metadata")

	input := services.UploadFileInput{
		FileHeader: fileHeader,
		UploadedBy: &userUUID,
		IsPublic:   isPublic,
		Metadata:   metadata,
	}

	result, err := h.fileService.UploadFile(c.Request.Context(), input)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == validate.ErrFileTooLarge {
			statusCode = http.StatusRequestEntityTooLarge
		}
		c.JSON(statusCode, models.ErrorResponse{
			Error:   "upload_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.FileUploadResponse{
		File:    result.File,
		URL:     result.URL,
		Message: "File uploaded successfully",
	})
}

// UploadMultipleFiles godoc
// @Summary Upload multiple files
// @Description Upload multiple files at once
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param files formData []file true "Files to upload"
// @Param is_public formData bool false "Whether the files are public (default: false)"
// @Param metadata formData string false "Additional metadata for the files"
// @Success 201 {object} models.MultipleFileUploadResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 413 {object} models.ErrorResponse
// @Router /api/v1/files/upload/multiple [post]
func (h *FileHandler) UploadMultipleFiles(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists || userID == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_user_id",
			Message: "Invalid user ID format",
		})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "multipart_error",
			Message: "Failed to parse multipart form",
			Details: err.Error(),
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "files_required",
			Message: "At least one file is required",
		})
		return
	}

	isPublic := c.PostForm("is_public") == "true"
	metadata := c.PostForm("metadata")

	results, err := h.fileService.UploadMultipleFiles(c.Request.Context(), files, &userUUID, isPublic, metadata)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == validate.ErrFileTooLarge {
			statusCode = http.StatusRequestEntityTooLarge
		}
		c.JSON(statusCode, models.ErrorResponse{
			Error:   "upload_failed",
			Message: err.Error(),
		})
		return
	}

	var responseFiles []models.FileUploadResponse
	for _, result := range results {
		responseFiles = append(responseFiles, models.FileUploadResponse{
			File:    result.File,
			URL:     result.URL,
			Message: "File uploaded successfully",
		})
	}

	c.JSON(http.StatusCreated, models.MultipleFileUploadResponse{
		Files:   responseFiles,
		Count:   len(responseFiles),
		Message: "Files uploaded successfully",
	})
}

// GetFile godoc
// @Summary Get file details
// @Description Get detailed information about a specific file
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "File ID"
// @Success 200 {object} models.File
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/files/{id} [get]
func (h *FileHandler) GetFile(c *gin.Context) {
	fileIDStr := c.Param("id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_file_id",
			Message: "Invalid file ID format",
		})
		return
	}

	file, err := h.fileService.GetFile(c.Request.Context(), fileID)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "file_not_found",
				Message: "File not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "get_file_failed",
				Message: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, file)
}

// DownloadFile godoc
// @Summary Download a file
// @Description Download a file by its ID
// @Tags files
// @Accept json
// @Produce application/octet-stream
// @Security BearerAuth
// @Param id path string true "File ID"
// @Success 200 {file} file
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/files/{id}/download [get]
func (h *FileHandler) DownloadFile(c *gin.Context) {
	fileIDStr := c.Param("id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_file_id",
			Message: "Invalid file ID format",
		})
		return
	}

	file, err := h.fileService.GetFile(c.Request.Context(), fileID)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "file_not_found",
				Message: "File not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "get_file_failed",
				Message: err.Error(),
			})
		}
		return
	}

	filePath, err := h.fileService.GetFilePath(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "file_path_error",
			Message: err.Error(),
		})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=\""+file.OriginalName+"\"")
	c.Header("Content-Type", file.MimeType)
	c.File(filePath)
}

// ListFiles godoc
// @Summary List user's files
// @Description Get a paginated list of files uploaded by the authenticated user
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Param type query string false "Filter by file type (image, video, document)"
// @Success 200 {object} models.FileListResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/files [get]
func (h *FileHandler) ListFiles(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists || userID == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_user_id",
			Message: "Invalid user ID format",
		})
		return
	}

	files, err := h.fileService.GetFilesByUser(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "list_files_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.FileListResponse{
		Files: files,
		Count: len(files),
	})
}

// DeleteFile godoc
// @Summary Delete a file
// @Description Delete a file by its ID (only owner can delete)
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "File ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/files/{id} [delete]
func (h *FileHandler) DeleteFile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists || userID == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	fileIDStr := c.Param("id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_file_id",
			Message: "Invalid file ID format",
		})
		return
	}

	file, err := h.fileService.GetFile(c.Request.Context(), fileID)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "file_not_found",
				Message: "File not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "get_file_failed",
				Message: err.Error(),
			})
		}
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_user_id",
			Message: "Invalid user ID format",
		})
		return
	}

	if file.UploadedBy != nil && *file.UploadedBy != userUUID {
		c.JSON(http.StatusForbidden, models.ErrorResponse{
			Error:   "access_denied",
			Message: "You can only delete your own files",
		})
		return
	}

	err = h.fileService.DeleteFile(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "delete_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "File deleted successfully",
	})
}

// GetFileURL godoc
// @Summary Get file URL
// @Description Get the public URL for a file
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "File ID"
// @Success 200 {object} models.FileURLResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/files/{id}/url [get]
func (h *FileHandler) GetFileURL(c *gin.Context) {
	fileIDStr := c.Param("id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_file_id",
			Message: "Invalid file ID format",
		})
		return
	}

	file, err := h.fileService.GetFile(c.Request.Context(), fileID)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "file_not_found",
				Message: "File not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "get_file_failed",
				Message: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.FileURLResponse{
		ID:           file.ID,
		OriginalName: file.OriginalName,
		URL:          file.URL,
		ThumbnailURL: file.ThumbnailURL,
		IsPublic:     file.IsPublic,
	})
}
