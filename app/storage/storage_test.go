package storage

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalStorage_Upload(t *testing.T) {
	tmpDir := t.TempDir()

	storage, err := NewLocalStorage(LocalConfig{
		BasePath: tmpDir,
		BaseURL:  "http://localhost:8080",
	})
	assert.NoError(t, err)

	content := []byte("test content")
	input := UploadInput{
		Key:         "test/file.txt",
		Content:     bytes.NewReader(content),
		ContentType: "text/plain",
		Size:        int64(len(content)),
	}

	result, err := storage.Upload(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test/file.txt", result.Key)
	assert.Contains(t, result.URL, "test/file.txt")

	fullPath := filepath.Join(tmpDir, "test", "file.txt")
	assert.FileExists(t, fullPath)

	data, err := os.ReadFile(fullPath)
	assert.NoError(t, err)
	assert.Equal(t, content, data)
}

func TestLocalStorage_Download(t *testing.T) {
	tmpDir := t.TempDir()

	storage, err := NewLocalStorage(LocalConfig{
		BasePath: tmpDir,
	})
	assert.NoError(t, err)

	content := []byte("test content")
	key := "test/file.txt"

	input := UploadInput{
		Key:     key,
		Content: bytes.NewReader(content),
	}
	_, err = storage.Upload(context.Background(), input)
	assert.NoError(t, err)

	downloaded, err := storage.Download(context.Background(), key)
	assert.NoError(t, err)
	assert.Equal(t, content, downloaded)
}

func TestLocalStorage_Delete(t *testing.T) {
	tmpDir := t.TempDir()

	storage, err := NewLocalStorage(LocalConfig{
		BasePath: tmpDir,
	})
	assert.NoError(t, err)

	content := []byte("test content")
	key := "test/file.txt"

	input := UploadInput{
		Key:     key,
		Content: bytes.NewReader(content),
	}
	_, err = storage.Upload(context.Background(), input)
	assert.NoError(t, err)

	err = storage.Delete(context.Background(), key)
	assert.NoError(t, err)

	fullPath := filepath.Join(tmpDir, key)
	assert.NoFileExists(t, fullPath)
}

func TestLocalStorage_Exists(t *testing.T) {
	tmpDir := t.TempDir()

	storage, err := NewLocalStorage(LocalConfig{
		BasePath: tmpDir,
	})
	assert.NoError(t, err)

	key := "test/file.txt"

	exists, err := storage.Exists(context.Background(), key)
	assert.NoError(t, err)
	assert.False(t, exists)

	content := []byte("test content")
	input := UploadInput{
		Key:     key,
		Content: bytes.NewReader(content),
	}
	_, err = storage.Upload(context.Background(), input)
	assert.NoError(t, err)

	exists, err = storage.Exists(context.Background(), key)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestLocalStorage_List(t *testing.T) {
	tmpDir := t.TempDir()

	storage, err := NewLocalStorage(LocalConfig{
		BasePath: tmpDir,
	})
	assert.NoError(t, err)

	files := []string{
		"test/file1.txt",
		"test/file2.txt",
		"other/file3.txt",
	}

	for _, file := range files {
		input := UploadInput{
			Key:     file,
			Content: bytes.NewReader([]byte("content")),
		}
		_, err := storage.Upload(context.Background(), input)
		assert.NoError(t, err)
	}

	allFiles, err := storage.List(context.Background(), "")
	assert.NoError(t, err)
	assert.Len(t, allFiles, 3)

	testFiles, err := storage.List(context.Background(), "test")
	assert.NoError(t, err)
	assert.Len(t, testFiles, 2)
}

func TestLocalStorage_GetURL(t *testing.T) {
	storage, err := NewLocalStorage(LocalConfig{
		BasePath: "./uploads",
		BaseURL:  "http://localhost:8080",
	})
	assert.NoError(t, err)

	url, err := storage.GetURL(context.Background(), "test/file.txt")
	assert.NoError(t, err)
	assert.Equal(t, "http://localhost:8080/uploads/test/file.txt", url)
}

func TestLocalStorage_CopyObject(t *testing.T) {
	tmpDir := t.TempDir()

	storage, err := NewLocalStorage(LocalConfig{
		BasePath: tmpDir,
	})
	assert.NoError(t, err)

	content := []byte("test content")
	sourceKey := "source/file.txt"
	destKey := "dest/file.txt"

	input := UploadInput{
		Key:     sourceKey,
		Content: bytes.NewReader(content),
	}
	_, err = storage.Upload(context.Background(), input)
	assert.NoError(t, err)

	err = storage.CopyObject(context.Background(), sourceKey, destKey)
	assert.NoError(t, err)

	downloaded, err := storage.Download(context.Background(), destKey)
	assert.NoError(t, err)
	assert.Equal(t, content, downloaded)
}

func TestLocalStorage_GetObjectMetadata(t *testing.T) {
	tmpDir := t.TempDir()

	storage, err := NewLocalStorage(LocalConfig{
		BasePath: tmpDir,
	})
	assert.NoError(t, err)

	content := []byte("test content")
	key := "test/file.txt"

	input := UploadInput{
		Key:     key,
		Content: bytes.NewReader(content),
	}
	_, err = storage.Upload(context.Background(), input)
	assert.NoError(t, err)

	metadata, err := storage.GetObjectMetadata(context.Background(), key)
	assert.NoError(t, err)
	assert.NotNil(t, metadata)
	assert.Contains(t, metadata, "size")
	assert.Contains(t, metadata, "modified")
}

func TestStorageFactory(t *testing.T) {
	factory := NewStorageFactory(StorageTypeLocal)

	localStorage, err := NewLocalStorage(LocalConfig{
		BasePath: t.TempDir(),
	})
	assert.NoError(t, err)

	factory.RegisterProvider(StorageTypeLocal, localStorage)

	provider, err := factory.GetProvider(StorageTypeLocal)
	assert.NoError(t, err)
	assert.NotNil(t, provider)

	defaultProvider, err := factory.GetDefaultProvider()
	assert.NoError(t, err)
	assert.NotNil(t, defaultProvider)

	_, err = factory.GetProvider("invalid")
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidStorageType)
}

func TestInitializeStorage(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		config      StorageConfig
		expectError bool
	}{
		{
			name: "local storage",
			config: StorageConfig{
				Type: StorageTypeLocal,
				Local: &LocalConfig{
					BasePath: tmpDir,
				},
			},
			expectError: false,
		},
		{
			name: "missing local config",
			config: StorageConfig{
				Type: StorageTypeLocal,
			},
			expectError: true,
		},
		{
			name: "invalid storage type",
			config: StorageConfig{
				Type: "invalid",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := InitializeStorage(context.Background(), tt.config)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, provider)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, provider)
			}
		})
	}
}
