package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFile_IsImage(t *testing.T) {
	tests := []struct {
		name     string
		mimeType string
		expected bool
	}{
		{
			name:     "jpeg image",
			mimeType: "image/jpeg",
			expected: true,
		},
		{
			name:     "png image",
			mimeType: "image/png",
			expected: true,
		},
		{
			name:     "gif image",
			mimeType: "image/gif",
			expected: true,
		},
		{
			name:     "pdf document",
			mimeType: "application/pdf",
			expected: false,
		},
		{
			name:     "video file",
			mimeType: "video/mp4",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &File{MimeType: tt.mimeType}
			result := file.IsImage()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFile_IsVideo(t *testing.T) {
	tests := []struct {
		name     string
		mimeType string
		expected bool
	}{
		{
			name:     "mp4 video",
			mimeType: "video/mp4",
			expected: true,
		},
		{
			name:     "quicktime video",
			mimeType: "video/quicktime",
			expected: true,
		},
		{
			name:     "image file",
			mimeType: "image/jpeg",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &File{MimeType: tt.mimeType}
			result := file.IsVideo()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFile_IsDocument(t *testing.T) {
	tests := []struct {
		name     string
		mimeType string
		expected bool
	}{
		{
			name:     "pdf document",
			mimeType: "application/pdf",
			expected: true,
		},
		{
			name:     "word document",
			mimeType: "application/msword",
			expected: true,
		},
		{
			name:     "excel spreadsheet",
			mimeType: "application/vnd.ms-excel",
			expected: true,
		},
		{
			name:     "image file",
			mimeType: "image/jpeg",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &File{MimeType: tt.mimeType}
			result := file.IsDocument()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFile_GetSizeInMB(t *testing.T) {
	tests := []struct {
		name     string
		fileSize int64
		expected float64
	}{
		{
			name:     "1 MB",
			fileSize: 1024 * 1024,
			expected: 1.0,
		},
		{
			name:     "5 MB",
			fileSize: 5 * 1024 * 1024,
			expected: 5.0,
		},
		{
			name:     "0.5 MB",
			fileSize: 512 * 1024,
			expected: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &File{FileSize: tt.fileSize}
			result := file.GetSizeInMB()
			assert.InDelta(t, tt.expected, result, 0.01)
		})
	}
}

func TestFile_GetSizeInKB(t *testing.T) {
	tests := []struct {
		name     string
		fileSize int64
		expected float64
	}{
		{
			name:     "1 KB",
			fileSize: 1024,
			expected: 1.0,
		},
		{
			name:     "10 KB",
			fileSize: 10 * 1024,
			expected: 10.0,
		},
		{
			name:     "0.5 KB",
			fileSize: 512,
			expected: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &File{FileSize: tt.fileSize}
			result := file.GetSizeInKB()
			assert.InDelta(t, tt.expected, result, 0.01)
		})
	}
}

func TestFile_BeforeCreate(t *testing.T) {
	file := &File{
		OriginalName: "test.jpg",
		FileName:     "test-123.jpg",
	}

	assert.Equal(t, uuid.Nil, file.ID)

	err := file.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, file.ID)
}

func TestFile_TableName(t *testing.T) {
	file := File{}
	assert.Equal(t, "files", file.TableName())
}
