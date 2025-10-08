package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedisClient_SetAndGet(t *testing.T) {
	t.Skip("Requires Redis server")

	client, err := NewRedisClient(RedisConfig{
		Host: "localhost",
		Port: 6379,
	})
	assert.NoError(t, err)
	defer client.Close()

	ctx := context.Background()
	key := "test:key"
	value := "test value"

	err = client.Set(ctx, key, value, time.Minute)
	assert.NoError(t, err)

	result, err := client.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, value, result)

	err = client.Delete(ctx, key)
	assert.NoError(t, err)
}

func TestRedisClient_Increment(t *testing.T) {
	t.Skip("Requires Redis server")

	client, err := NewRedisClient(RedisConfig{
		Host: "localhost",
		Port: 6379,
	})
	assert.NoError(t, err)
	defer client.Close()

	ctx := context.Background()
	key := "test:counter"

	val, err := client.Increment(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), val)

	val, err = client.Increment(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), val)

	err = client.Delete(ctx, key)
	assert.NoError(t, err)
}

func TestRedisClient_Hash(t *testing.T) {
	t.Skip("Requires Redis server")

	client, err := NewRedisClient(RedisConfig{
		Host: "localhost",
		Port: 6379,
	})
	assert.NoError(t, err)
	defer client.Close()

	ctx := context.Background()
	key := "test:hash"

	err = client.HSet(ctx, key, "field1", "value1", "field2", "value2")
	assert.NoError(t, err)

	val, err := client.HGet(ctx, key, "field1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)

	all, err := client.HGetAll(ctx, key)
	assert.NoError(t, err)
	assert.Len(t, all, 2)
	assert.Equal(t, "value1", all["field1"])
	assert.Equal(t, "value2", all["field2"])

	err = client.Delete(ctx, key)
	assert.NoError(t, err)
}

func TestRedisClient_Set(t *testing.T) {
	t.Skip("Requires Redis server")

	client, err := NewRedisClient(RedisConfig{
		Host: "localhost",
		Port: 6379,
	})
	assert.NoError(t, err)
	defer client.Close()

	ctx := context.Background()
	key := "test:set"

	err = client.SAdd(ctx, key, "member1", "member2", "member3")
	assert.NoError(t, err)

	members, err := client.SMembers(ctx, key)
	assert.NoError(t, err)
	assert.Len(t, members, 3)

	isMember, err := client.SIsMember(ctx, key, "member1")
	assert.NoError(t, err)
	assert.True(t, isMember)

	err = client.Delete(ctx, key)
	assert.NoError(t, err)
}

func TestRedisClient_Exists(t *testing.T) {
	t.Skip("Requires Redis server")

	client, err := NewRedisClient(RedisConfig{
		Host: "localhost",
		Port: 6379,
	})
	assert.NoError(t, err)
	defer client.Close()

	ctx := context.Background()
	key := "test:exists"

	count, err := client.Exists(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)

	err = client.Set(ctx, key, "value", time.Minute)
	assert.NoError(t, err)

	count, err = client.Exists(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	err = client.Delete(ctx, key)
	assert.NoError(t, err)
}

func TestRedisClient_TTL(t *testing.T) {
	t.Skip("Requires Redis server")

	client, err := NewRedisClient(RedisConfig{
		Host: "localhost",
		Port: 6379,
	})
	assert.NoError(t, err)
	defer client.Close()

	ctx := context.Background()
	key := "test:ttl"

	err = client.Set(ctx, key, "value", 10*time.Second)
	assert.NoError(t, err)

	ttl, err := client.TTL(ctx, key)
	assert.NoError(t, err)
	assert.Greater(t, ttl, time.Duration(0))
	assert.LessOrEqual(t, ttl, 10*time.Second)

	err = client.Delete(ctx, key)
	assert.NoError(t, err)
}

func TestRedisClient_SetNX(t *testing.T) {
	t.Skip("Requires Redis server")

	client, err := NewRedisClient(RedisConfig{
		Host: "localhost",
		Port: 6379,
	})
	assert.NoError(t, err)
	defer client.Close()

	ctx := context.Background()
	key := "test:setnx"

	ok, err := client.SetNX(ctx, key, "value1", time.Minute)
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = client.SetNX(ctx, key, "value2", time.Minute)
	assert.NoError(t, err)
	assert.False(t, ok)

	val, err := client.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)

	err = client.Delete(ctx, key)
	assert.NoError(t, err)
}

func TestRedisClient_Ping(t *testing.T) {
	t.Skip("Requires Redis server")

	client, err := NewRedisClient(RedisConfig{
		Host: "localhost",
		Port: 6379,
	})
	assert.NoError(t, err)
	defer client.Close()

	err = client.Ping(context.Background())
	assert.NoError(t, err)
}
