package services

import (
	"bytes"
	"mime/multipart"
	"net/textproto"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createMockFileHeader(filename string, size int64, content []byte) *multipart.FileHeader {
	return &multipart.FileHeader{
		Filename: filename,
		Size:     size,
		Header:   textproto.MIMEHeader{},
	}
}

func TestFileValidator_ValidateFile(t *testing.T) {
	validator := NewFileValidator(FileValidationConfig{
		MaxFileSize:       5 * 1024 * 1024, // 5MB
		MinFileSize:       1,
		AllowedExtensions: []string{".jpg", ".png", ".pdf"},
	})

	tests := []struct {
		name        string
		fileHeader  *multipart.FileHeader
		expectError bool
		expectedErr error
	}{
		{
			name:        "valid file",
			fileHeader:  createMockFileHeader("test.jpg", 1024*1024, nil),
			expectError: false,
		},
		{
			name:        "file too large",
			fileHeader:  createMockFileHeader("large.jpg", 10*1024*1024, nil),
			expectError: true,
			expectedErr: ErrFileTooLarge,
		},
		{
			name:        "empty file",
			fileHeader:  createMockFileHeader("empty.jpg", 0, nil),
			expectError: true,
			expectedErr: ErrEmptyFile,
		},
		{
			name:        "invalid extension",
			fileHeader:  createMockFileHeader("test.exe", 1024, nil),
			expectError: true,
			expectedErr: ErrInvalidExtension,
		},
		{
			name:        "invalid filename with path traversal",
			fileHeader:  createMockFileHeader("../../../etc/passwd", 1024, nil),
			expectError: true,
			expectedErr: ErrInvalidFileName,
		},
		{
			name:        "nil file header",
			fileHeader:  nil,
			expectError: true,
			expectedErr: ErrEmptyFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateFile(tt.fileHeader)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.ErrorIs(t, err, tt.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileValidator_ValidateFileContent(t *testing.T) {
	validator := NewFileValidator(FileValidationConfig{
		AllowedMimeTypes: []string{"image/jpeg", "image/png", "application/pdf"},
	})

	tests := []struct {
		name        string
		mimeType    string
		expectError bool
	}{
		{
			name:        "allowed mime type",
			mimeType:    "image/jpeg",
			expectError: false,
		},
		{
			name:        "disallowed mime type",
			mimeType:    "application/x-executable",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateFileContent(nil, tt.mimeType)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileValidator_IsExtensionAllowed(t *testing.T) {
	validator := NewFileValidator(FileValidationConfig{
		AllowedExtensions: []string{".jpg", ".png", ".pdf"},
	})

	tests := []struct {
		name     string
		ext      string
		expected bool
	}{
		{
			name:     "allowed extension",
			ext:      ".jpg",
			expected: true,
		},
		{
			name:     "allowed extension uppercase",
			ext:      ".JPG",
			expected: true,
		},
		{
			name:     "disallowed extension",
			ext:      ".exe",
			expected: false,
		},
		{
			name:     "extension without dot",
			ext:      "png",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.isExtensionAllowed(tt.ext)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFileValidator_IsMimeTypeAllowed(t *testing.T) {
	validator := NewFileValidator(FileValidationConfig{
		AllowedMimeTypes: []string{"image/*", "application/pdf"},
	})

	tests := []struct {
		name     string
		mimeType string
		expected bool
	}{
		{
			name:     "exact match",
			mimeType: "application/pdf",
			expected: true,
		},
		{
			name:     "wildcard match",
			mimeType: "image/jpeg",
			expected: true,
		},
		{
			name:     "wildcard match png",
			mimeType: "image/png",
			expected: true,
		},
		{
			name:     "no match",
			mimeType: "video/mp4",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.isMimeTypeAllowed(tt.mimeType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFileValidator_GetMaxFileSizeInMB(t *testing.T) {
	validator := NewFileValidator(FileValidationConfig{
		MaxFileSize: 10 * 1024 * 1024, // 10MB
	})

	result := validator.GetMaxFileSizeInMB()
	assert.Equal(t, 10.0, result)
}

func TestDetectMimeTypeFromBytes(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected string
	}{
		{
			name:     "JPEG image",
			data:     []byte{0xFF, 0xD8, 0xFF, 0xE0},
			expected: "image/jpeg",
		},
		{
			name:     "PNG image",
			data:     []byte{0x89, 0x50, 0x4E, 0x47},
			expected: "image/png",
		},
		{
			name:     "PDF document",
			data:     []byte("%PDF-1.4"),
			expected: "application/pdf",
		},
		{
			name:     "GIF image",
			data:     []byte("GIF89a"),
			expected: "image/gif",
		},
		{
			name:     "unknown type",
			data:     []byte{0x00, 0x01, 0x02, 0x03},
			expected: "application/octet-stream",
		},
		{
			name:     "too short",
			data:     []byte{0x00},
			expected: "application/octet-stream",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectMimeTypeFromBytes(tt.data)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPredefinedValidators(t *testing.T) {
	t.Run("ImageValidator", func(t *testing.T) {
		assert.NotNil(t, ImageValidator)
		assert.Equal(t, float64(5), ImageValidator.GetMaxFileSizeInMB())
		assert.True(t, ImageValidator.isMimeTypeAllowed("image/jpeg"))
		assert.False(t, ImageValidator.isMimeTypeAllowed("application/pdf"))
	})

	t.Run("DocumentValidator", func(t *testing.T) {
		assert.NotNil(t, DocumentValidator)
		assert.Equal(t, float64(20), DocumentValidator.GetMaxFileSizeInMB())
		assert.True(t, DocumentValidator.isMimeTypeAllowed("application/pdf"))
		assert.False(t, DocumentValidator.isMimeTypeAllowed("image/jpeg"))
	})

	t.Run("VideoValidator", func(t *testing.T) {
		assert.NotNil(t, VideoValidator)
		assert.Equal(t, float64(100), VideoValidator.GetMaxFileSizeInMB())
		assert.True(t, VideoValidator.isMimeTypeAllowed("video/mp4"))
		assert.False(t, VideoValidator.isMimeTypeAllowed("image/jpeg"))
	})
}

func TestFileValidator_ValidateFileName(t *testing.T) {
	validator := NewFileValidator(FileValidationConfig{})

	tests := []struct {
		name        string
		filename    string
		expectError bool
	}{
		{
			name:        "valid filename",
			filename:    "document.pdf",
			expectError: false,
		},
		{
			name:        "empty filename",
			filename:    "",
			expectError: true,
		},
		{
			name:        "path traversal",
			filename:    "../../../etc/passwd",
			expectError: true,
		},
		{
			name:        "null byte",
			filename:    "file\x00.txt",
			expectError: true,
		},
		{
			name:        "forward slash",
			filename:    "path/to/file.txt",
			expectError: true,
		},
		{
			name:        "backslash",
			filename:    "path\\to\\file.txt",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateFileName(tt.filename)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMin(t *testing.T) {
	assert.Equal(t, 5, min(5, 10))
	assert.Equal(t, 5, min(10, 5))
	assert.Equal(t, 5, min(5, 5))
}
