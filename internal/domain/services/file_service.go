package services

import (
	"base-gin/internal/domain/models"
	"base-gin/internal/domain/repository"
	"base-gin/internal/pkg/utils"
	"base-gin/internal/pkg/validate"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileService struct {
	*BaseService
	fileRepo  repository.BaseRepository
	validator *validate.FileValidator
	uploadDir string
	baseURL   string
}

type FileServiceConfig struct {
	DB        *gorm.DB
	FileRepo  repository.BaseRepository
	Validator *validate.FileValidator
	UploadDir string
	BaseURL   string
}

func NewFileService(config FileServiceConfig) *FileService {
	if config.UploadDir == "" {
		config.UploadDir = "./uploads"
	}
	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:8080"
	}
	if config.Validator == nil {
		config.Validator = validate.NewFileValidator(validate.FileValidationConfig{
			MaxFileSize: 10 * 1024 * 1024, // 10MB
		})
	}

	return &FileService{
		BaseService: NewBaseService(config.DB),
		fileRepo:    config.FileRepo,
		validator:   config.Validator,
		uploadDir:   config.UploadDir,
		baseURL:     config.BaseURL,
	}
}

type UploadFileInput struct {
	FileHeader *multipart.FileHeader
	UploadedBy *uuid.UUID
	IsPublic   bool
	Metadata   string
}

type UploadFileResult struct {
	File *models.File
	URL  string
}

func (s *FileService) UploadFile(ctx context.Context, input UploadFileInput) (*UploadFileResult, error) {
	if err := s.validator.ValidateFile(input.FileHeader); err != nil {
		return nil, err
	}

	file, err := input.FileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	mimeType, err := utils.DetectMimeType(file)
	if err != nil {
		return nil, fmt.Errorf("failed to detect mime type: %w", err)
	}

	if err := s.validator.ValidateFileContent(file, mimeType); err != nil {
		return nil, err
	}

	fileName, err := s.generateUniqueFileName(input.FileHeader.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to generate file name: %w", err)
	}

	subDir := s.getSubDirectory()
	fullUploadDir := filepath.Join(s.uploadDir, subDir)

	if err := os.MkdirAll(fullUploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	filePath := filepath.Join(fullUploadDir, fileName)

	if err := s.saveFile(file, filePath); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(input.FileHeader.Filename))
	relativePath := filepath.Join(subDir, fileName)
	fileURL := s.generateFileURL(relativePath)

	// Process metadata - convert to valid JSON or null
	var metadataValue string
	if strings.TrimSpace(input.Metadata) == "" {
		metadataValue = "null"
	} else {
		// If it's already valid JSON, use it as is
		var temp interface{}
		if json.Unmarshal([]byte(input.Metadata), &temp) == nil {
			metadataValue = input.Metadata
		} else {
			// If it's plain text, wrap it in quotes to make it a valid JSON string
			metadataValue = fmt.Sprintf(`"%s"`, strings.ReplaceAll(input.Metadata, `"`, `\"`))
		}
	}

	fileModel := &models.File{
		OriginalName: input.FileHeader.Filename,
		FileName:     fileName,
		FilePath:     relativePath,
		FileSize:     input.FileHeader.Size,
		MimeType:     mimeType,
		Extension:    ext,
		StorageType:  models.StorageTypeLocal,
		URL:          fileURL,
		UploadedBy:   input.UploadedBy,
		IsPublic:     input.IsPublic,
		Metadata:     metadataValue,
	}

	if err := s.fileRepo.Create(ctx, fileModel); err != nil {
		_ = os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}

	return &UploadFileResult{
		File: fileModel,
		URL:  fileURL,
	}, nil
}

func (s *FileService) UploadMultipleFiles(ctx context.Context, fileHeaders []*multipart.FileHeader, uploadedBy *uuid.UUID, isPublic bool, metadata string) ([]*UploadFileResult, error) {
	results := make([]*UploadFileResult, 0, len(fileHeaders))

	for _, fileHeader := range fileHeaders {
		result, err := s.UploadFile(ctx, UploadFileInput{
			FileHeader: fileHeader,
			UploadedBy: uploadedBy,
			IsPublic:   isPublic,
			Metadata:   metadata,
		})

		if err != nil {
			return results, err
		}

		results = append(results, result)
	}

	return results, nil
}

func (s *FileService) GetFile(ctx context.Context, fileID uuid.UUID) (*models.File, error) {
	var file models.File
	if err := s.fileRepo.FindByID(ctx, fileID, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

func (s *FileService) DeleteFile(ctx context.Context, fileID uuid.UUID) error {
	var file models.File
	if err := s.fileRepo.FindByID(ctx, fileID, &file); err != nil {
		return err
	}

	if file.StorageType == models.StorageTypeLocal {
		fullPath := filepath.Join(s.uploadDir, file.FilePath)
		if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete physical file: %w", err)
		}
	}

	if err := s.fileRepo.Delete(ctx, fileID, &models.File{}); err != nil {
		return fmt.Errorf("failed to delete file metadata: %w", err)
	}

	return nil
}

func (s *FileService) GetFilesByUser(ctx context.Context, userID uuid.UUID) ([]models.File, error) {
	var files []models.File
	err := s.db.WithContext(ctx).
		Where("uploaded_by = ?", userID).
		Order("created_at DESC").
		Find(&files).Error

	if err != nil {
		return nil, err
	}

	return files, nil
}

func (s *FileService) generateUniqueFileName(originalName string) (string, error) {
	ext := filepath.Ext(originalName)
	nameWithoutExt := strings.TrimSuffix(originalName, ext)

	timestamp := time.Now().UnixNano()
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%s-%d", nameWithoutExt, timestamp)))
	hashStr := hex.EncodeToString(hash.Sum(nil))[:16]

	return fmt.Sprintf("%s-%s%s", nameWithoutExt, hashStr, ext), nil
}

func (s *FileService) getSubDirectory() string {
	now := time.Now()
	return filepath.Join(
		fmt.Sprintf("%d", now.Year()),
		fmt.Sprintf("%02d", now.Month()),
		fmt.Sprintf("%02d", now.Day()),
	)
}

func (s *FileService) generateFileURL(relativePath string) string {
	cleanPath := filepath.ToSlash(relativePath)
	return fmt.Sprintf("%s/uploads/%s", s.baseURL, cleanPath)
}

func (s *FileService) saveFile(src multipart.File, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return err
	}

	return nil
}

func (s *FileService) GetFileContent(ctx context.Context, fileID uuid.UUID) ([]byte, error) {
	var file models.File
	if err := s.fileRepo.FindByID(ctx, fileID, &file); err != nil {
		return nil, err
	}

	if file.StorageType != models.StorageTypeLocal {
		return nil, errors.New("only local files can be read directly")
	}

	fullPath := filepath.Join(s.uploadDir, file.FilePath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return content, nil
}

func (s *FileService) GetFilePath(ctx context.Context, fileID uuid.UUID) (string, error) {
	var file models.File
	if err := s.fileRepo.FindByID(ctx, fileID, &file); err != nil {
		return "", err
	}

	if file.StorageType != models.StorageTypeLocal {
		return "", errors.New("only local files have physical paths")
	}

	return filepath.Join(s.uploadDir, file.FilePath), nil
}
