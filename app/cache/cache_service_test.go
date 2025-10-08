package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCacheProvider struct {
	mock.Mock
}

func (m *MockCacheProvider) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCacheProvider) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCacheProvider) Delete(ctx context.Context, keys ...string) error {
	args := m.Called(ctx, keys)
	return args.Error(0)
}

func (m *MockCacheProvider) Exists(ctx context.Context, keys ...string) (int64, error) {
	args := m.Called(ctx, keys)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCacheProvider) Expire(ctx context.Context, key string, expiration time.Duration) error {
	args := m.Called(ctx, key, expiration)
	return args.Error(0)
}

func (m *MockCacheProvider) TTL(ctx context.Context, key string) (time.Duration, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *MockCacheProvider) Increment(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCacheProvider) Decrement(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCacheProvider) FlushAll(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCacheProvider) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCacheProvider) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestCacheService_SetAndGet(t *testing.T) {
	mockProvider := new(MockCacheProvider)
	service := NewCacheService(CacheServiceConfig{
		Provider: mockProvider,
		Prefix:   "test",
	})

	ctx := context.Background()
	key := "user:1"
	value := map[string]string{"name": "John", "email": "john@example.com"}

	serialized := `{"name":"John","email":"john@example.com"}`
	mockProvider.On("Set", ctx, "test:user:1", []byte(serialized), time.Minute).Return(nil)
	mockProvider.On("Get", ctx, "test:user:1").Return(serialized, nil)

	err := service.Set(ctx, key, value, time.Minute)
	assert.NoError(t, err)

	var result map[string]string
	err = service.Get(ctx, key, &result)
	assert.NoError(t, err)
	assert.Equal(t, value, result)

	mockProvider.AssertExpectations(t)
}

func TestCacheService_Delete(t *testing.T) {
	mockProvider := new(MockCacheProvider)
	service := NewCacheService(CacheServiceConfig{
		Provider: mockProvider,
		Prefix:   "test",
	})

	ctx := context.Background()
	keys := []string{"key1", "key2"}
	fullKeys := []string{"test:key1", "test:key2"}

	mockProvider.On("Delete", ctx, fullKeys).Return(nil)

	err := service.Delete(ctx, keys...)
	assert.NoError(t, err)

	mockProvider.AssertExpectations(t)
}

func TestCacheService_Exists(t *testing.T) {
	mockProvider := new(MockCacheProvider)
	service := NewCacheService(CacheServiceConfig{
		Provider: mockProvider,
		Prefix:   "test",
	})

	ctx := context.Background()
	keys := []string{"key1", "key2"}
	fullKeys := []string{"test:key1", "test:key2"}

	mockProvider.On("Exists", ctx, fullKeys).Return(int64(2), nil)

	exists, err := service.Exists(ctx, keys...)
	assert.NoError(t, err)
	assert.True(t, exists)

	mockProvider.AssertExpectations(t)
}

func TestCacheService_Remember(t *testing.T) {
	mockProvider := new(MockCacheProvider)
	service := NewCacheService(CacheServiceConfig{
		Provider: mockProvider,
		Prefix:   "test",
	})

	ctx := context.Background()
	key := "data"
	value := map[string]string{"result": "from callback"}

	mockProvider.On("Get", ctx, "test:data").Return("", ErrCacheKeyNotFound)
	serialized := `{"result":"from callback"}`
	mockProvider.On("Set", ctx, "test:data", []byte(serialized), time.Minute).Return(nil)

	var result map[string]string
	err := service.Remember(ctx, key, time.Minute, func() (interface{}, error) {
		return value, nil
	}, &result)

	assert.NoError(t, err)
	assert.Equal(t, value, result)

	mockProvider.AssertExpectations(t)
}

func TestCacheService_Remember_CacheHit(t *testing.T) {
	mockProvider := new(MockCacheProvider)
	service := NewCacheService(CacheServiceConfig{
		Provider: mockProvider,
		Prefix:   "test",
	})

	ctx := context.Background()
	key := "data"
	cachedValue := `{"result":"from cache"}`

	mockProvider.On("Get", ctx, "test:data").Return(cachedValue, nil)

	var result map[string]string
	callbackCalled := false
	err := service.Remember(ctx, key, time.Minute, func() (interface{}, error) {
		callbackCalled = true
		return nil, nil
	}, &result)

	assert.NoError(t, err)
	assert.False(t, callbackCalled)
	assert.Equal(t, "from cache", result["result"])

	mockProvider.AssertExpectations(t)
}

func TestCacheService_Remember_CallbackError(t *testing.T) {
	mockProvider := new(MockCacheProvider)
	service := NewCacheService(CacheServiceConfig{
		Provider: mockProvider,
		Prefix:   "test",
	})

	ctx := context.Background()
	key := "data"

	mockProvider.On("Get", ctx, "test:data").Return("", ErrCacheKeyNotFound)

	var result map[string]string
	err := service.Remember(ctx, key, time.Minute, func() (interface{}, error) {
		return nil, errors.New("callback error")
	}, &result)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "callback failed")

	mockProvider.AssertExpectations(t)
}

func TestCacheService_Flush(t *testing.T) {
	mockProvider := new(MockCacheProvider)
	service := NewCacheService(CacheServiceConfig{
		Provider: mockProvider,
		Prefix:   "test",
	})

	ctx := context.Background()

	mockProvider.On("FlushAll", ctx).Return(nil)

	err := service.Flush(ctx)
	assert.NoError(t, err)

	mockProvider.AssertExpectations(t)
}

func TestCacheService_Increment(t *testing.T) {
	mockProvider := new(MockCacheProvider)
	service := NewCacheService(CacheServiceConfig{
		Provider: mockProvider,
		Prefix:   "test",
	})

	ctx := context.Background()
	key := "counter"

	mockProvider.On("Increment", ctx, "test:counter").Return(int64(5), nil)

	val, err := service.Increment(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), val)

	mockProvider.AssertExpectations(t)
}

func TestJSONSerializer(t *testing.T) {
	serializer := &JSONSerializer{}

	data := map[string]interface{}{
		"name":  "John",
		"age":   30,
		"email": "john@example.com",
	}

	serialized, err := serializer.Serialize(data)
	assert.NoError(t, err)
	assert.NotEmpty(t, serialized)

	var result map[string]interface{}
	err = serializer.Deserialize(serialized, &result)
	assert.NoError(t, err)
	assert.Equal(t, "John", result["name"])
	assert.Equal(t, float64(30), result["age"])
}

func TestQueryCache_GenerateKey(t *testing.T) {
	mockProvider := new(MockCacheProvider)
	cacheService := NewCacheService(CacheServiceConfig{
		Provider: mockProvider,
		Prefix:   "test",
	})
	queryCache := NewQueryCache(cacheService, time.Minute)

	query := "SELECT * FROM users WHERE id = ?"
	params := []interface{}{123}

	key1 := queryCache.generateKey(query, params)
	key2 := queryCache.generateKey(query, params)

	assert.Equal(t, key1, key2)

	params2 := []interface{}{456}
	key3 := queryCache.generateKey(query, params2)

	assert.NotEqual(t, key1, key3)
}
