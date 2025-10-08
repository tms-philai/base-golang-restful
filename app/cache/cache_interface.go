package cache

import (
	"context"
	"errors"
	"time"
)

var (
	ErrCacheKeyNotFound = errors.New("cache key not found")
	ErrCacheMiss        = errors.New("cache miss")
	ErrInvalidValue     = errors.New("invalid cache value")
)

type CacheProvider interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	Increment(ctx context.Context, key string) (int64, error)
	Decrement(ctx context.Context, key string) (int64, error)
	FlushAll(ctx context.Context) error
	Ping(ctx context.Context) error
	Close() error
}

type CacheSerializer interface {
	Serialize(value interface{}) ([]byte, error)
	Deserialize(data []byte, dest interface{}) error
}

const (
	DefaultExpiration = 5 * time.Minute
	NoExpiration      = 0
	ForeverExpiration = -1
)
