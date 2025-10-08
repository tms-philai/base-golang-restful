package storage

import (
	"context"
	"io"
)

type StorageProvider interface {
	Upload(ctx context.Context, input UploadInput) (*UploadResult, error)
	Download(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	GetURL(ctx context.Context, key string) (string, error)
	GetSignedURL(ctx context.Context, key string, expiration int64) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
	List(ctx context.Context, prefix string) ([]string, error)
}

type UploadInput struct {
	Key         string
	Content     io.Reader
	ContentType string
	Size        int64
	Metadata    map[string]string
	ACL         string
}

type UploadResult struct {
	Key      string
	URL      string
	ETag     string
	Size     int64
	Metadata map[string]string
}

const (
	ACLPrivate           = "private"
	ACLPublicRead        = "public-read"
	ACLPublicReadWrite   = "public-read-write"
	ACLAuthenticatedRead = "authenticated-read"
)
