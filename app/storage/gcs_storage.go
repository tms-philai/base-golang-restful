package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCSStorage struct {
	client    *storage.Client
	bucket    string
	projectID string
	baseURL   string
}

type GCSConfig struct {
	ProjectID       string
	Bucket          string
	CredentialsFile string
	CredentialsJSON []byte
}

func NewGCSStorage(ctx context.Context, config GCSConfig) (*GCSStorage, error) {
	var client *storage.Client
	var err error

	if config.CredentialsFile != "" {
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(config.CredentialsFile))
	} else if len(config.CredentialsJSON) > 0 {
		client, err = storage.NewClient(ctx, option.WithCredentialsJSON(config.CredentialsJSON))
	} else {
		client, err = storage.NewClient(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	baseURL := fmt.Sprintf("https://storage.googleapis.com/%s", config.Bucket)

	return &GCSStorage{
		client:    client,
		bucket:    config.Bucket,
		projectID: config.ProjectID,
		baseURL:   baseURL,
	}, nil
}

func (g *GCSStorage) Upload(ctx context.Context, input UploadInput) (*UploadResult, error) {
	obj := g.client.Bucket(g.bucket).Object(input.Key)
	writer := obj.NewWriter(ctx)

	writer.ContentType = input.ContentType

	if len(input.Metadata) > 0 {
		writer.Metadata = input.Metadata
	}

	if input.ACL == ACLPublicRead {
		writer.PredefinedACL = "publicRead"
	} else if input.ACL == ACLPrivate {
		writer.PredefinedACL = "private"
	}

	if _, err := io.Copy(writer, input.Content); err != nil {
		writer.Close()
		return nil, fmt.Errorf("failed to write to GCS: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close GCS writer: %w", err)
	}

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get object attributes: %w", err)
	}

	url := fmt.Sprintf("%s/%s", g.baseURL, input.Key)

	return &UploadResult{
		Key:      input.Key,
		URL:      url,
		ETag:     attrs.Etag,
		Size:     attrs.Size,
		Metadata: attrs.Metadata,
	}, nil
}

func (g *GCSStorage) Download(ctx context.Context, key string) ([]byte, error) {
	obj := g.client.Bucket(g.bucket).Object(key)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS reader: %w", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read from GCS: %w", err)
	}

	return data, nil
}

func (g *GCSStorage) Delete(ctx context.Context, key string) error {
	obj := g.client.Bucket(g.bucket).Object(key)

	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete from GCS: %w", err)
	}

	return nil
}

func (g *GCSStorage) GetURL(ctx context.Context, key string) (string, error) {
	return fmt.Sprintf("%s/%s", g.baseURL, key), nil
}

func (g *GCSStorage) GetSignedURL(ctx context.Context, key string, expiration int64) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(time.Duration(expiration) * time.Second),
	}

	url, err := storage.SignedURL(g.bucket, key, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return url, nil
}

func (g *GCSStorage) Exists(ctx context.Context, key string) (bool, error) {
	obj := g.client.Bucket(g.bucket).Object(key)

	_, err := obj.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check object existence: %w", err)
	}

	return true, nil
}

func (g *GCSStorage) List(ctx context.Context, prefix string) ([]string, error) {
	query := &storage.Query{}
	if prefix != "" {
		query.Prefix = prefix
	}

	var keys []string
	it := g.client.Bucket(g.bucket).Objects(ctx, query)

	for {
		attrs, err := it.Next()
		if err == storage.ErrObjectNotExist {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", err)
		}

		keys = append(keys, attrs.Name)
	}

	return keys, nil
}

func (g *GCSStorage) CopyObject(ctx context.Context, sourceKey, destKey string) error {
	src := g.client.Bucket(g.bucket).Object(sourceKey)
	dst := g.client.Bucket(g.bucket).Object(destKey)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return fmt.Errorf("failed to copy object: %w", err)
	}

	return nil
}

func (g *GCSStorage) GetObjectMetadata(ctx context.Context, key string) (map[string]string, error) {
	obj := g.client.Bucket(g.bucket).Object(key)
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get object metadata: %w", err)
	}

	return attrs.Metadata, nil
}

func (g *GCSStorage) Close() error {
	return g.client.Close()
}
