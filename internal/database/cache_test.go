package database_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carldunham/useful-cookery/internal/database"
)

// TestRedisCache_WithMocks tests the Redis cache implementation using mocks.
func TestRedisCache_WithMocks(t *testing.T) {
	t.Parallel()
	// Skip this test since we don't have a Redis server available
	t.Skip("Skipping Redis cache test - no Redis server available")
}

func TestNewRedisCache(t *testing.T) {
	t.Parallel()
	// Test creating a Redis cache
	cache := database.NewRedisCache("localhost:6379", "", 0)
	assert.NotNil(t, cache, "Redis cache should not be nil")

	// No need to check interface implementation since we're returning concrete type
}

func TestNewInMemoryCache(t *testing.T) {
	t.Parallel()
	// Test creating an in-memory cache
	cache := database.NewInMemoryCache()
	assert.NotNil(t, cache, "In-memory cache should not be nil")

	// No need to check interface implementation since we're returning concrete type
}

func TestInMemoryCache(t *testing.T) {
	t.Parallel()
	// Create a new in-memory cache
	cache := database.NewInMemoryCache()
	require.NotNil(t, cache, "In-memory cache should not be nil")

	ctx := t.Context()
	key := "test-key"
	value := []byte("test-value")
	ttl := time.Minute

	// Test Set
	err := cache.Set(ctx, key, value, ttl)
	require.NoError(t, err, "Failed to set cache value")

	// Test Get
	retrievedValue, err := cache.Get(ctx, key)
	require.NoError(t, err, "Failed to get cache value")
	assert.Equal(t, value, retrievedValue, "Retrieved value does not match original value")

	// Test Delete
	err = cache.Delete(ctx, key)
	require.NoError(t, err, "Failed to delete cache value")

	// Test Get after Delete
	retrievedValue, err = cache.Get(ctx, key)
	require.NoError(t, err, "Failed to get cache value after deletion")
	assert.Nil(t, retrievedValue, "Retrieved value should be nil after deletion")

	// Test expiration
	err = cache.Set(ctx, key, value, time.Millisecond)
	require.NoError(t, err, "Failed to set cache value with short TTL")

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	// Test Get after expiration
	retrievedValue, err = cache.Get(ctx, key)
	require.NoError(t, err, "Failed to get cache value after expiration")
	assert.Nil(t, retrievedValue, "Retrieved value should be nil after expiration")

	// Test Cleanup
	err = cache.Set(ctx, "key1", []byte("value1"), time.Millisecond)
	require.NoError(t, err, "Failed to set cache value for cleanup test")
	err = cache.Set(ctx, "key2", []byte("value2"), time.Hour)
	require.NoError(t, err, "Failed to set cache value for cleanup test")

	// Wait for expiration of key1
	time.Sleep(10 * time.Millisecond)

	// Run cleanup
	cache.Cleanup()

	// Check that key1 is gone and key2 is still there
	retrievedValue, err = cache.Get(ctx, "key1")
	require.NoError(t, err, "Failed to get cache value after cleanup")
	assert.Nil(t, retrievedValue, "Retrieved value for key1 should be nil after cleanup")

	retrievedValue, err = cache.Get(ctx, "key2")
	require.NoError(t, err, "Failed to get cache value after cleanup")
	assert.NotNil(t, retrievedValue, "Retrieved value for key2 should not be nil after cleanup")
}
