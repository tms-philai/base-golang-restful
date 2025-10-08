package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

func NewRedisClient(config RedisConfig) (*RedisClient, error) {
	if config.Host == "" {
		config.Host = "localhost"
	}
	if config.Port == 0 {
		config.Port = 6379
	}
	if config.PoolSize == 0 {
		config.PoolSize = 10
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{
		client: client,
	}, nil
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrCacheKeyNotFound
	}
	if err != nil {
		return "", fmt.Errorf("failed to get key: %w", err)
	}
	return val, nil
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key: %w", err)
	}
	return nil
}

func (r *RedisClient) Delete(ctx context.Context, keys ...string) error {
	err := r.client.Del(ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("failed to delete keys: %w", err)
	}
	return nil
}

func (r *RedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	count, err := r.client.Exists(ctx, keys...).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to check existence: %w", err)
	}
	return count, nil
}

func (r *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := r.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiration: %w", err)
	}
	return nil
}

func (r *RedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get TTL: %w", err)
	}
	return ttl, nil
}

func (r *RedisClient) Increment(ctx context.Context, key string) (int64, error) {
	val, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment: %w", err)
	}
	return val, nil
}

func (r *RedisClient) Decrement(ctx context.Context, key string) (int64, error) {
	val, err := r.client.Decr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to decrement: %w", err)
	}
	return val, nil
}

func (r *RedisClient) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	val, err := r.client.IncrBy(ctx, key, value).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment by value: %w", err)
	}
	return val, nil
}

func (r *RedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	ok, err := r.client.SetNX(ctx, key, value, expiration).Result()
	if err != nil {
		return false, fmt.Errorf("failed to set if not exists: %w", err)
	}
	return ok, nil
}

func (r *RedisClient) GetSet(ctx context.Context, key string, value interface{}) (string, error) {
	oldVal, err := r.client.GetSet(ctx, key, value).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get and set: %w", err)
	}
	return oldVal, nil
}

func (r *RedisClient) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	vals, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get multiple keys: %w", err)
	}
	return vals, nil
}

func (r *RedisClient) MSet(ctx context.Context, pairs ...interface{}) error {
	err := r.client.MSet(ctx, pairs...).Err()
	if err != nil {
		return fmt.Errorf("failed to set multiple keys: %w", err)
	}
	return nil
}

func (r *RedisClient) Keys(ctx context.Context, pattern string) ([]string, error) {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys: %w", err)
	}
	return keys, nil
}

func (r *RedisClient) FlushDB(ctx context.Context) error {
	err := r.client.FlushDB(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to flush database: %w", err)
	}
	return nil
}

func (r *RedisClient) FlushAll(ctx context.Context) error {
	err := r.client.FlushAll(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to flush all databases: %w", err)
	}
	return nil
}

func (r *RedisClient) HSet(ctx context.Context, key string, values ...interface{}) error {
	err := r.client.HSet(ctx, key, values...).Err()
	if err != nil {
		return fmt.Errorf("failed to set hash field: %w", err)
	}
	return nil
}

func (r *RedisClient) HGet(ctx context.Context, key, field string) (string, error) {
	val, err := r.client.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return "", ErrCacheKeyNotFound
	}
	if err != nil {
		return "", fmt.Errorf("failed to get hash field: %w", err)
	}
	return val, nil
}

func (r *RedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	vals, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get all hash fields: %w", err)
	}
	return vals, nil
}

func (r *RedisClient) HDel(ctx context.Context, key string, fields ...string) error {
	err := r.client.HDel(ctx, key, fields...).Err()
	if err != nil {
		return fmt.Errorf("failed to delete hash fields: %w", err)
	}
	return nil
}

func (r *RedisClient) HExists(ctx context.Context, key, field string) (bool, error) {
	exists, err := r.client.HExists(ctx, key, field).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check hash field existence: %w", err)
	}
	return exists, nil
}

func (r *RedisClient) SAdd(ctx context.Context, key string, members ...interface{}) error {
	err := r.client.SAdd(ctx, key, members...).Err()
	if err != nil {
		return fmt.Errorf("failed to add to set: %w", err)
	}
	return nil
}

func (r *RedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	members, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get set members: %w", err)
	}
	return members, nil
}

func (r *RedisClient) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	isMember, err := r.client.SIsMember(ctx, key, member).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check set membership: %w", err)
	}
	return isMember, nil
}

func (r *RedisClient) SRem(ctx context.Context, key string, members ...interface{}) error {
	err := r.client.SRem(ctx, key, members...).Err()
	if err != nil {
		return fmt.Errorf("failed to remove from set: %w", err)
	}
	return nil
}

func (r *RedisClient) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	err := r.client.ZAdd(ctx, key, members...).Err()
	if err != nil {
		return fmt.Errorf("failed to add to sorted set: %w", err)
	}
	return nil
}

func (r *RedisClient) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	members, err := r.client.ZRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get sorted set range: %w", err)
	}
	return members, nil
}

func (r *RedisClient) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	members, err := r.client.ZRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get sorted set range with scores: %w", err)
	}
	return members, nil
}

func (r *RedisClient) ZRem(ctx context.Context, key string, members ...interface{}) error {
	err := r.client.ZRem(ctx, key, members...).Err()
	if err != nil {
		return fmt.Errorf("failed to remove from sorted set: %w", err)
	}
	return nil
}

func (r *RedisClient) Publish(ctx context.Context, channel string, message interface{}) error {
	err := r.client.Publish(ctx, channel, message).Err()
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}

func (r *RedisClient) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return r.client.Subscribe(ctx, channels...)
}

func (r *RedisClient) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}

func (r *RedisClient) TxPipeline() redis.Pipeliner {
	return r.client.TxPipeline()
}

func (r *RedisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}
