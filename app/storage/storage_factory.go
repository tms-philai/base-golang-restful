package storage

import (
	"context"
	"errors"
	"fmt"
)

type StorageType string

const (
	StorageTypeLocal StorageType = "local"
	StorageTypeS3    StorageType = "s3"
	StorageTypeGCS   StorageType = "gcs"
)

var (
	ErrInvalidStorageType = errors.New("invalid storage type")
	ErrMissingConfig      = errors.New("missing storage configuration")
)

type StorageFactory struct {
	defaultType StorageType
	providers   map[StorageType]StorageProvider
}

func NewStorageFactory(defaultType StorageType) *StorageFactory {
	return &StorageFactory{
		defaultType: defaultType,
		providers:   make(map[StorageType]StorageProvider),
	}
}

func (f *StorageFactory) RegisterProvider(storageType StorageType, provider StorageProvider) {
	f.providers[storageType] = provider
}

func (f *StorageFactory) GetProvider(storageType StorageType) (StorageProvider, error) {
	if storageType == "" {
		storageType = f.defaultType
	}

	provider, exists := f.providers[storageType]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrInvalidStorageType, storageType)
	}

	return provider, nil
}

func (f *StorageFactory) GetDefaultProvider() (StorageProvider, error) {
	return f.GetProvider(f.defaultType)
}

type StorageConfig struct {
	Type  StorageType
	Local *LocalConfig
	S3    *S3Config
	GCS   *GCSConfig
}

func InitializeStorage(ctx context.Context, config StorageConfig) (StorageProvider, error) {
	switch config.Type {
	case StorageTypeLocal:
		if config.Local == nil {
			return nil, fmt.Errorf("%w: local config", ErrMissingConfig)
		}
		return NewLocalStorage(*config.Local)

	case StorageTypeS3:
		if config.S3 == nil {
			return nil, fmt.Errorf("%w: S3 config", ErrMissingConfig)
		}
		return NewS3Storage(*config.S3)

	case StorageTypeGCS:
		if config.GCS == nil {
			return nil, fmt.Errorf("%w: GCS config", ErrMissingConfig)
		}
		return NewGCSStorage(ctx, *config.GCS)

	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidStorageType, config.Type)
	}
}

func InitializeMultipleStorages(ctx context.Context, configs []StorageConfig, defaultType StorageType) (*StorageFactory, error) {
	factory := NewStorageFactory(defaultType)

	for _, config := range configs {
		provider, err := InitializeStorage(ctx, config)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize %s storage: %w", config.Type, err)
		}

		factory.RegisterProvider(config.Type, provider)
	}

	return factory, nil
}
