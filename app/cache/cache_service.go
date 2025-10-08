package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type CacheService struct {
	provider   CacheProvider
	serializer CacheSerializer
	prefix     string
}

type CacheServiceConfig struct {
	Provider   CacheProvider
	Serializer CacheSerializer
	Prefix     string
}

func NewCacheService(config CacheServiceConfig) *CacheService {
	if config.Serializer == nil {
		config.Serializer = &JSONSerializer{}
	}
	if config.Prefix == "" {
		config.Prefix = "app"
	}

	return &CacheService{
		provider:   config.Provider,
		serializer: config.Serializer,
		prefix:     config.Prefix,
	}
}

func (s *CacheService) buildKey(key string) string {
	return fmt.Sprintf("%s:%s", s.prefix, key)
}

func (s *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	fullKey := s.buildKey(key)
	
	data, err := s.provider.Get(ctx, fullKey)
	if err != nil {
		return err
	}

	if err := s.serializer.Deserialize([]byte(data), dest); err != nil {
		return fmt.Errorf("failed to deserialize: %w", err)
	}

	return nil
}

func (s *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	fullKey := s.buildKey(key)

	data, err := s.serializer.Serialize(value)
	if err != nil {
		return fmt.Errorf("failed to serialize: %w", err)
	}

	return s.provider.Set(ctx, fullKey, data, expiration)
}

func (s *CacheService) Delete(ctx context.Context, keys ...string) error {
	fullKeys := make([]string, len(keys))
	for i, key := range keys {
		fullKeys[i] = s.buildKey(key)
	}

	return s.provider.Delete(ctx, fullKeys...)
}

func (s *CacheService) Exists(ctx context.Context, keys ...string) (bool, error) {
	fullKeys := make([]string, len(keys))
	for i, key := range keys {
		fullKeys[i] = s.buildKey(key)
	}

	count, err := s.provider.Exists(ctx, fullKeys...)
	if err != nil {
		return false, err
	}

	return count == int64(len(keys)), nil
}

func (s *CacheService) Remember(ctx context.Context, key string, expiration time.Duration, callback func() (interface{}, error), dest interface{}) error {
	err := s.Get(ctx, key, dest)
	if err == nil {
		return nil
	}

	if err != ErrCacheKeyNotFound {
		return err
	}

	value, err := callback()
	if err != nil {
		return fmt.Errorf("callback failed: %w", err)
	}

	if err := s.Set(ctx, key, value, expiration); err != nil {
		return err
	}

	data, err := s.serializer.Serialize(value)
	if err != nil {
		return err
	}

	return s.serializer.Deserialize(data, dest)
}

func (s *CacheService) RememberForever(ctx context.Context, key string, callback func() (interface{}, error), dest interface{}) error {
	return s.Remember(ctx, key, NoExpiration, callback, dest)
}

func (s *CacheService) Forget(ctx context.Context, keys ...string) error {
	return s.Delete(ctx, keys...)
}

func (s *CacheService) Flush(ctx context.Context) error {
	return s.provider.FlushAll(ctx)
}

func (s *CacheService) TTL(ctx context.Context, key string) (time.Duration, error) {
	fullKey := s.buildKey(key)
	return s.provider.TTL(ctx, fullKey)
}

func (s *CacheService) Expire(ctx context.Context, key string, expiration time.Duration) error {
	fullKey := s.buildKey(key)
	return s.provider.Expire(ctx, fullKey, expiration)
}

func (s *CacheService) Increment(ctx context.Context, key string) (int64, error) {
	fullKey := s.buildKey(key)
	return s.provider.Increment(ctx, fullKey)
}

func (s *CacheService) Decrement(ctx context.Context, key string) (int64, error) {
	fullKey := s.buildKey(key)
	return s.provider.Decrement(ctx, fullKey)
}

func (s *CacheService) Tags(tags ...string) *TaggedCache {
	return &TaggedCache{
		service: s,
		tags:    tags,
	}
}

type TaggedCache struct {
	service *CacheService
	tags    []string
}

func (t *TaggedCache) buildTaggedKey(key string) string {
	tagStr := ""
	for _, tag := range t.tags {
		tagStr += tag + ":"
	}
	return tagStr + key
}

func (t *TaggedCache) Get(ctx context.Context, key string, dest interface{}) error {
	return t.service.Get(ctx, t.buildTaggedKey(key), dest)
}

func (t *TaggedCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return t.service.Set(ctx, t.buildTaggedKey(key), value, expiration)
}

func (t *TaggedCache) Delete(ctx context.Context, keys ...string) error {
	taggedKeys := make([]string, len(keys))
	for i, key := range keys {
		taggedKeys[i] = t.buildTaggedKey(key)
	}
	return t.service.Delete(ctx, taggedKeys...)
}

func (t *TaggedCache) Flush(ctx context.Context) error {
	return t.service.Delete(ctx, t.tags...)
}

type JSONSerializer struct{}

func (j *JSONSerializer) Serialize(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (j *JSONSerializer) Deserialize(data []byte, dest interface{}) error {
	return json.Unmarshal(data, dest)
}
