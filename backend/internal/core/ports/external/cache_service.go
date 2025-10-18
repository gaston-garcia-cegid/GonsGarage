package external

import (
	"context"
	"time"
)

// CacheService defines caching operations interface
type CacheService interface {
	// Set value with expiration
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Get value by key
	Get(ctx context.Context, key string) (string, error)

	// Get and unmarshal JSON to struct
	GetJSON(ctx context.Context, key string, dest interface{}) error

	// Set JSON value
	SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Delete key
	Delete(ctx context.Context, key string) error

	// Check if key exists
	Exists(ctx context.Context, key string) (bool, error)

	// Set with no expiration
	SetPermanent(ctx context.Context, key string, value interface{}) error

	// Increment counter
	Increment(ctx context.Context, key string) (int64, error)

	// Set expiration for existing key
	Expire(ctx context.Context, key string, expiration time.Duration) error

	// Get multiple keys
	MGet(ctx context.Context, keys ...string) ([]string, error)

	// Set multiple key-value pairs
	MSet(ctx context.Context, pairs map[string]interface{}) error

	// Clear all cache (use with caution)
	FlushAll(ctx context.Context) error
}

// Method implementations should not be provided for interfaces in Go.
// Remove these implementations or move them to a struct that implements CacheService.
