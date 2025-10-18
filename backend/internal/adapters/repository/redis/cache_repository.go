package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports/repositories"
	"github.com/redis/go-redis/v9"
)

type cacheRepository struct {
	client *redis.Client
}

func NewRedisCacheRepository(client *redis.Client) repositories.CacheRepository {
	return &cacheRepository{client: client}
}

func (r *cacheRepository) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil // Chave não existe
		}
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *cacheRepository) Set(ctx context.Context, key string, value interface{}, ttl int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, time.Duration(ttl)*time.Second).Err()
}

func (r *cacheRepository) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
