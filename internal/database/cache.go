package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Cache interface defines common caching operations.
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

// RedisCache implements the Cache interface using Redis.
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache client.
func NewRedisCache(addr, password string, db int) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{
		client: client,
	}
}

// Get retrieves a value from the cache.
func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // Key not found, return nil value
		}
		return nil, fmt.Errorf("failed to get value from Redis: %w", err)
	}
	return val, nil
}

// Set stores a value in the cache with the given TTL.
func (c *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	err := c.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set value in Redis: %w", err)
	}
	return nil
}

// Delete removes a value from the cache.
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete value from Redis: %w", err)
	}
	return nil
}

// Close closes the Redis connection.
func (c *RedisCache) Close() error {
	err := c.client.Close()
	if err != nil {
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}
	return nil
}

// Ping checks if the Redis server is available.
func (c *RedisCache) Ping(ctx context.Context) error {
	err := c.client.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to ping Redis server: %w", err)
	}
	return nil
}

// InMemoryCache implements the Cache interface using a map.
// Used for testing or when Redis is not available.
type InMemoryCache struct {
	data map[string]cacheItem
}

type cacheItem struct {
	value      []byte
	expiration time.Time
}

// NewInMemoryCache creates a new in-memory cache.
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]cacheItem),
	}
}

// Get retrieves a value from the in-memory cache.
func (c *InMemoryCache) Get(_ context.Context, key string) ([]byte, error) {
	item, ok := c.data[key]
	if !ok {
		return nil, nil
	}

	// Check if expired
	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		delete(c.data, key)
		return nil, nil
	}

	return item.value, nil
}

// Set stores a value in the in-memory cache.
func (c *InMemoryCache) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	var expiration time.Time
	if ttl > 0 {
		expiration = time.Now().Add(ttl)
	}

	c.data[key] = cacheItem{
		value:      value,
		expiration: expiration,
	}

	return nil
}

// Delete removes a value from the in-memory cache.
func (c *InMemoryCache) Delete(_ context.Context, key string) error {
	delete(c.data, key)
	return nil
}

// Cleanup removes expired items from the cache.
func (c *InMemoryCache) Cleanup() {
	now := time.Now()
	for key, item := range c.data {
		if !item.expiration.IsZero() && now.After(item.expiration) {
			delete(c.data, key)
		}
	}
}
