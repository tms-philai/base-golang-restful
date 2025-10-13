package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	OriginalName string         `gorm:"type:varchar(255);not null" json:"original_name"`
	FileName     string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"file_name"`
	FilePath     string         `gorm:"type:varchar(500);not null" json:"file_path"`
	FileSize     int64          `gorm:"not null" json:"file_size"`
	MimeType     string         `gorm:"type:varchar(100);not null" json:"mime_type"`
	Extension    string         `gorm:"type:varchar(20);not null" json:"extension"`
	StorageType  string         `gorm:"type:varchar(50);not null;default:'local'" json:"storage_type"`
	URL          string         `gorm:"type:varchar(1000)" json:"url,omitempty"`
	ThumbnailURL string         `gorm:"type:varchar(1000)" json:"thumbnail_url,omitempty"`
	UploadedBy   *uuid.UUID     `gorm:"type:uuid;index" json:"uploaded_by,omitempty"`
	User         *User          `gorm:"foreignKey:UploadedBy" json:"user,omitempty"`
	Metadata     string         `gorm:"type:jsonb" json:"metadata,omitempty"`
	IsPublic     bool           `gorm:"default:false;not null" json:"is_public"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (File) TableName() string {
	return "files"
}

func (f *File) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

const (
	StorageTypeLocal = "local"
	StorageTypeS3    = "s3"
	StorageTypeGCS   = "gcs"
	StorageTypeAzure = "azure"
)

func (f *File) IsImage() bool {
	imageTypes := map[string]bool{
		"image/jpeg":    true,
		"image/jpg":     true,
		"image/png":     true,
		"image/gif":     true,
		"image/webp":    true,
		"image/svg+xml": true,
	}
	return imageTypes[f.MimeType]
}

func (f *File) IsVideo() bool {
	videoTypes := map[string]bool{
		"video/mp4":       true,
		"video/mpeg":      true,
		"video/quicktime": true,
		"video/x-msvideo": true,
		"video/webm":      true,
	}
	return videoTypes[f.MimeType]
}

func (f *File) IsDocument() bool {
	docTypes := map[string]bool{
		"application/pdf":    true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/vnd.ms-excel": true,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         true,
		"application/vnd.ms-powerpoint":                                             true,
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
		"text/plain": true,
	}
	return docTypes[f.MimeType]
}

func (f *File) GetSizeInMB() float64 {
	return float64(f.FileSize) / (1024 * 1024)
}

func (f *File) GetSizeInKB() float64 {
	return float64(f.FileSize) / 1024
}
