package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/gaston-garcia-cegid/gonsgarage/backend/internal/core/ports/external"
)

// RedisCacheRepository implements CacheService interface
type RedisCacheRepository struct {
	client *redis.Client
}

// NewRedisCacheRepository creates a new Redis cache repository
func NewRedisCacheRepository(client *redis.Client) external.CacheService {
	return &RedisCacheRepository{
		client: client,
	}
}

// Set implements CacheService.Set
func (r *RedisCacheRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache key %s: %w", key, err)
	}
	return nil
}

// Get implements CacheService.Get
func (r *RedisCacheRepository) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", external.ErrCacheKeyNotFound
		}
		return "", fmt.Errorf("failed to get cache key %s: %w", key, err)
	}
	return val, nil
}

// GetJSON implements CacheService.GetJSON
func (r *RedisCacheRepository) GetJSON(ctx context.Context, key string, dest interface{}) error {
	val, err := r.Get(ctx, key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("failed to unmarshal JSON for key %s: %w", key, err)
	}

	return nil
}

// SetJSON implements CacheService.SetJSON
func (r *RedisCacheRepository) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON for key %s: %w", key, err)
	}

	return r.Set(ctx, key, jsonData, expiration)
}

// Delete implements CacheService.Delete
func (r *RedisCacheRepository) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete cache key %s: %w", key, err)
	}
	return nil
}

// Exists implements CacheService.Exists
func (r *RedisCacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}
	return count > 0, nil
}

// SetPermanent implements CacheService.SetPermanent
func (r *RedisCacheRepository) SetPermanent(ctx context.Context, key string, value interface{}) error {
	return r.Set(ctx, key, value, 0) // 0 means no expiration
}

// Increment implements CacheService.Increment
func (r *RedisCacheRepository) Increment(ctx context.Context, key string) (int64, error) {
	val, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key %s: %w", key, err)
	}
	return val, nil
}

// Expire implements CacheService.Expire
func (r *RedisCacheRepository) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := r.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiration for key %s: %w", key, err)
	}
	return nil
}

// MGet implements CacheService.MGet
func (r *RedisCacheRepository) MGet(ctx context.Context, keys ...string) ([]string, error) {
	vals, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get multiple keys: %w", err)
	}

	results := make([]string, len(vals))
	for i, val := range vals {
		if val != nil {
			results[i] = val.(string)
		}
	}

	return results, nil
}

// MSet implements CacheService.MSet
func (r *RedisCacheRepository) MSet(ctx context.Context, pairs map[string]interface{}) error {
	// Convert map to slice for Redis MSET
	args := make([]interface{}, 0, len(pairs)*2)
	for key, value := range pairs {
		args = append(args, key, value)
	}

	err := r.client.MSet(ctx, args...).Err()
	if err != nil {
		return fmt.Errorf("failed to set multiple keys: %w", err)
	}
	return nil
}

// FlushAll implements CacheService.FlushAll
func (r *RedisCacheRepository) FlushAll(ctx context.Context) error {
	err := r.client.FlushAll(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to flush all cache: %w", err)
	}
	return nil
}
