package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	basePath string
	baseURL  string
}

type LocalConfig struct {
	BasePath string
	BaseURL  string
}

func NewLocalStorage(config LocalConfig) (*LocalStorage, error) {
	if config.BasePath == "" {
		config.BasePath = "./uploads"
	}
	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:8080"
	}

	if err := os.MkdirAll(config.BasePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}

	return &LocalStorage{
		basePath: config.BasePath,
		baseURL:  config.BaseURL,
	}, nil
}

func (l *LocalStorage) Upload(ctx context.Context, input UploadInput) (*UploadResult, error) {
	fullPath := filepath.Join(l.basePath, input.Key)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	written, err := io.Copy(file, input.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	url := l.generateURL(input.Key)

	return &UploadResult{
		Key:      input.Key,
		URL:      url,
		Size:     written,
		Metadata: input.Metadata,
	}, nil
}

func (l *LocalStorage) Download(ctx context.Context, key string) ([]byte, error) {
	fullPath := filepath.Join(l.basePath, key)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", key)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

func (l *LocalStorage) Delete(ctx context.Context, key string) error {
	fullPath := filepath.Join(l.basePath, key)

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	dir := filepath.Dir(fullPath)
	l.cleanEmptyDirs(dir)

	return nil
}

func (l *LocalStorage) GetURL(ctx context.Context, key string) (string, error) {
	return l.generateURL(key), nil
}

func (l *LocalStorage) GetSignedURL(ctx context.Context, key string, expiration int64) (string, error) {
	return l.generateURL(key), nil
}

func (l *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
	fullPath := filepath.Join(l.basePath, key)

	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}

	return true, nil
}

func (l *LocalStorage) List(ctx context.Context, prefix string) ([]string, error) {
	var keys []string

	searchPath := l.basePath
	if prefix != "" {
		searchPath = filepath.Join(l.basePath, prefix)
	}

	err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(l.basePath, path)
			if err != nil {
				return err
			}

			relPath = filepath.ToSlash(relPath)
			keys = append(keys, relPath)
		}

		return nil
	})

	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return keys, nil
}

func (l *LocalStorage) CopyObject(ctx context.Context, sourceKey, destKey string) error {
	sourcePath := filepath.Join(l.basePath, sourceKey)
	destPath := filepath.Join(l.basePath, destKey)

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

func (l *LocalStorage) GetObjectMetadata(ctx context.Context, key string) (map[string]string, error) {
	fullPath := filepath.Join(l.basePath, key)

	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	metadata := map[string]string{
		"size":         fmt.Sprintf("%d", info.Size()),
		"modified":     info.ModTime().Format("2006-01-02 15:04:05"),
		"is_directory": fmt.Sprintf("%t", info.IsDir()),
	}

	return metadata, nil
}

func (l *LocalStorage) generateURL(key string) string {
	cleanKey := filepath.ToSlash(key)
	return fmt.Sprintf("%s/uploads/%s", l.baseURL, cleanKey)
}

func (l *LocalStorage) cleanEmptyDirs(dir string) {
	if dir == l.basePath {
		return
	}

	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) > 0 {
		return
	}

	os.Remove(dir)

	parentDir := filepath.Dir(dir)
	l.cleanEmptyDirs(parentDir)
}
