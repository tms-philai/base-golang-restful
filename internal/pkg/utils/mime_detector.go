package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"strings"
)

// DetectMimeType detects the MIME type of a file by reading its content
func DetectMimeType(file io.Reader) (string, error) {
	// Read first 512 bytes for MIME detection
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Reset file position for later use
	if seeker, ok := file.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	content := buffer[:n]
	mimeType := detectMimeTypeFromContent(content)

	// Fallback to application/octet-stream if detection fails
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return mimeType, nil
}

// detectMimeTypeFromContent detects MIME type from file content
func detectMimeTypeFromContent(content []byte) string {
	// Check for common file signatures
	if len(content) < 4 {
		return ""
	}

	// Image formats
	if bytes.HasPrefix(content, []byte{0xFF, 0xD8, 0xFF}) {
		return "image/jpeg"
	}
	if bytes.HasPrefix(content, []byte{0x89, 0x50, 0x4E, 0x47}) {
		return "image/png"
	}
	if bytes.HasPrefix(content, []byte{0x47, 0x49, 0x46, 0x38}) {
		return "image/gif"
	}
	if bytes.HasPrefix(content, []byte{0x52, 0x49, 0x46, 0x46}) && len(content) > 8 && bytes.HasPrefix(content[8:], []byte{0x57, 0x45, 0x42, 0x50}) {
		return "image/webp"
	}

	// Document formats
	if bytes.HasPrefix(content, []byte{0x25, 0x50, 0x44, 0x46}) {
		return "application/pdf"
	}
	if bytes.HasPrefix(content, []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}) {
		return "application/msword"
	}
	if bytes.HasPrefix(content, []byte{0x50, 0x4B, 0x03, 0x04}) {
		// Could be various ZIP-based formats, check further
		if bytes.Contains(content, []byte("word/")) {
			return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		}
		if bytes.Contains(content, []byte("xl/")) {
			return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		}
		if bytes.Contains(content, []byte("ppt/")) {
			return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
		}
		return "application/zip"
	}

	// Video formats
	if bytes.HasPrefix(content, []byte{0x00, 0x00, 0x00, 0x18, 0x66, 0x74, 0x79, 0x70}) ||
		bytes.HasPrefix(content, []byte{0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70}) {
		return "video/mp4"
	}
	if bytes.HasPrefix(content, []byte{0x1A, 0x45, 0xDF, 0xA3}) {
		return "video/webm"
	}

	// Audio formats
	if bytes.HasPrefix(content, []byte{0x49, 0x44, 0x33}) {
		return "audio/mpeg"
	}
	if bytes.HasPrefix(content, []byte{0xFF, 0xFB}) || bytes.HasPrefix(content, []byte{0xFF, 0xF3}) || bytes.HasPrefix(content, []byte{0xFF, 0xF2}) {
		return "audio/mpeg"
	}

	// Archive formats
	if bytes.HasPrefix(content, []byte{0x50, 0x4B, 0x03, 0x04}) {
		return "application/zip"
	}
	if bytes.HasPrefix(content, []byte{0x1F, 0x8B}) {
		return "application/gzip"
	}
	if bytes.HasPrefix(content, []byte{0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C}) {
		return "application/x-7z-compressed"
	}

	// Text formats
	if isTextContent(content) {
		return "text/plain"
	}

	return ""
}

// isTextContent checks if content appears to be text
func isTextContent(content []byte) bool {
	// Check for null bytes (binary files usually have them)
	if bytes.Contains(content, []byte{0x00}) {
		return false
	}

	// Check if content is mostly printable ASCII
	printableCount := 0
	for _, b := range content {
		if (b >= 32 && b <= 126) || b == 9 || b == 10 || b == 13 {
			printableCount++
		}
	}

	// If more than 80% of characters are printable, consider it text
	return float64(printableCount)/float64(len(content)) > 0.8
}

// GetExtensionFromMimeType returns the appropriate file extension for a MIME type
func GetExtensionFromMimeType(mimeType string) string {
	ext, _ := mime.ExtensionsByType(mimeType)
	if len(ext) > 0 {
		return ext[0]
	}

	// Fallback mappings
	mimeToExt := map[string]string{
		"image/jpeg":               ".jpg",
		"image/png":                ".png",
		"image/gif":                ".gif",
		"image/webp":               ".webp",
		"application/pdf":          ".pdf",
		"text/plain":               ".txt",
		"application/json":         ".json",
		"application/xml":          ".xml",
		"text/html":                ".html",
		"text/css":                 ".css",
		"text/javascript":          ".js",
		"application/octet-stream": ".bin",
	}

	if ext, exists := mimeToExt[mimeType]; exists {
		return ext
	}

	return ""
}

// IsImageMimeType checks if the MIME type is an image
func IsImageMimeType(mimeType string) bool {
	return strings.HasPrefix(mimeType, "image/")
}

// IsVideoMimeType checks if the MIME type is a video
func IsVideoMimeType(mimeType string) bool {
	return strings.HasPrefix(mimeType, "video/")
}

// IsAudioMimeType checks if the MIME type is audio
func IsAudioMimeType(mimeType string) bool {
	return strings.HasPrefix(mimeType, "audio/")
}

// IsDocumentMimeType checks if the MIME type is a document
func IsDocumentMimeType(mimeType string) bool {
	documentTypes := []string{
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.ms-powerpoint",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"text/plain",
		"application/rtf",
	}

	for _, docType := range documentTypes {
		if mimeType == docType {
			return true
		}
	}
	return false
}
