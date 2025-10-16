package validate

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
)

var (
	ErrFileTooLarge     = errors.New("file size exceeds maximum allowed")
	ErrInvalidFileType  = errors.New("file type not allowed")
	ErrInvalidExtension = errors.New("file extension not allowed")
	ErrEmptyFile        = errors.New("file is empty")
	ErrInvalidFileName  = errors.New("invalid file name")
)

type FileValidationConfig struct {
	MaxFileSize       int64
	AllowedMimeTypes  []string
	AllowedExtensions []string
	MinFileSize       int64
}

type FileValidator struct {
	config FileValidationConfig
}

func NewFileValidator(config FileValidationConfig) *FileValidator {
	if config.MaxFileSize == 0 {
		config.MaxFileSize = 10 * 1024 * 1024 // 10MB default
	}
	if config.MinFileSize == 0 {
		config.MinFileSize = 1 // 1 byte minimum
	}

	return &FileValidator{
		config: config,
	}
}

func (v *FileValidator) ValidateFile(fileHeader *multipart.FileHeader) error {
	if fileHeader == nil {
		return ErrEmptyFile
	}

	if fileHeader.Size == 0 {
		return ErrEmptyFile
	}

	if fileHeader.Size < v.config.MinFileSize {
		return fmt.Errorf("file size must be at least %d bytes", v.config.MinFileSize)
	}

	if fileHeader.Size > v.config.MaxFileSize {
		return ErrFileTooLarge
	}

	if err := v.validateFileName(fileHeader.Filename); err != nil {
		return err
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if len(v.config.AllowedExtensions) > 0 {
		if !v.isExtensionAllowed(ext) {
			return ErrInvalidExtension
		}
	}

	return nil
}

func (v *FileValidator) ValidateFileContent(file multipart.File, mimeType string) error {
	if len(v.config.AllowedMimeTypes) > 0 {
		if !v.isMimeTypeAllowed(mimeType) {
			return ErrInvalidFileType
		}
	}

	return nil
}

func (v *FileValidator) validateFileName(filename string) error {
	if filename == "" {
		return ErrInvalidFileName
	}

	invalidChars := []string{"..", "/", "\\", "\x00"}
	for _, char := range invalidChars {
		if strings.Contains(filename, char) {
			return ErrInvalidFileName
		}
	}

	return nil
}

func (v *FileValidator) isExtensionAllowed(ext string) bool {
	ext = strings.ToLower(ext)
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	for _, allowed := range v.config.AllowedExtensions {
		if strings.ToLower(allowed) == ext {
			return true
		}
	}
	return false
}

func (v *FileValidator) isMimeTypeAllowed(mimeType string) bool {
	mimeType = strings.ToLower(mimeType)

	for _, allowed := range v.config.AllowedMimeTypes {
		if strings.ToLower(allowed) == mimeType {
			return true
		}

		if strings.HasSuffix(allowed, "/*") {
			prefix := strings.TrimSuffix(allowed, "/*")
			if strings.HasPrefix(mimeType, prefix+"/") {
				return true
			}
		}
	}
	return false
}

func (v *FileValidator) GetMaxFileSizeInMB() float64 {
	return float64(v.config.MaxFileSize) / (1024 * 1024)
}

var (
	ImageValidator = NewFileValidator(FileValidationConfig{
		MaxFileSize: 5 * 1024 * 1024, // 5MB
		AllowedMimeTypes: []string{
			"image/jpeg",
			"image/jpg",
			"image/png",
			"image/gif",
			"image/webp",
		},
		AllowedExtensions: []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
	})

	DocumentValidator = NewFileValidator(FileValidationConfig{
		MaxFileSize: 20 * 1024 * 1024, // 20MB
		AllowedMimeTypes: []string{
			"application/pdf",
			"application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			"application/vnd.ms-excel",
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			"text/plain",
		},
		AllowedExtensions: []string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".txt"},
	})

	VideoValidator = NewFileValidator(FileValidationConfig{
		MaxFileSize: 100 * 1024 * 1024, // 100MB
		AllowedMimeTypes: []string{
			"video/mp4",
			"video/mpeg",
			"video/quicktime",
			"video/webm",
		},
		AllowedExtensions: []string{".mp4", ".mpeg", ".mov", ".webm"},
	})
)

func DetectMimeType(file multipart.File) (string, error) {
	buffer := make([]byte, 512)

	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	mimeType := detectMimeTypeFromBytes(buffer[:n])
	return mimeType, nil
}

func detectMimeTypeFromBytes(data []byte) string {
	if len(data) < 4 {
		return "application/octet-stream"
	}

	signatures := map[string]string{
		"\xFF\xD8\xFF": "image/jpeg",
		"\x89PNG":      "image/png",
		"GIF87a":       "image/gif",
		"GIF89a":       "image/gif",
		"%PDF":         "application/pdf",
		"PK\x03\x04":   "application/zip",
		"\x1F\x8B":     "application/gzip",
	}

	dataStr := string(data[:min(len(data), 10)])

	for sig, mimeType := range signatures {
		if strings.HasPrefix(dataStr, sig) {
			return mimeType
		}
	}

	if strings.HasPrefix(dataStr, "RIFF") && len(data) >= 12 {
		if string(data[8:12]) == "WEBP" {
			return "image/webp"
		}
	}

	return "application/octet-stream"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
