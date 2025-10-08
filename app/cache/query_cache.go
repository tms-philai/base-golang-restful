package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type QueryCache struct {
	cache      *CacheService
	expiration time.Duration
}

func NewQueryCache(cache *CacheService, expiration time.Duration) *QueryCache {
	if expiration == 0 {
		expiration = 5 * time.Minute
	}

	return &QueryCache{
		cache:      cache,
		expiration: expiration,
	}
}

func (q *QueryCache) Remember(ctx context.Context, query string, params []interface{}, callback func() (interface{}, error), dest interface{}) error {
	key := q.generateKey(query, params)
	return q.cache.Remember(ctx, key, q.expiration, callback, dest)
}

func (q *QueryCache) Forget(ctx context.Context, query string, params []interface{}) error {
	key := q.generateKey(query, params)
	return q.cache.Delete(ctx, key)
}

func (q *QueryCache) ForgetByPattern(ctx context.Context, pattern string) error {
	return q.cache.Delete(ctx, pattern)
}

func (q *QueryCache) generateKey(query string, params []interface{}) string {
	hash := sha256.New()
	hash.Write([]byte(query))
	
	for _, param := range params {
		hash.Write([]byte(fmt.Sprintf("%v", param)))
	}
	
	hashStr := hex.EncodeToString(hash.Sum(nil))
	return fmt.Sprintf("query:%s", hashStr)
}

func (q *QueryCache) SetExpiration(expiration time.Duration) {
	q.expiration = expiration
}

func (q *QueryCache) GetExpiration() time.Duration {
	return q.expiration
}
