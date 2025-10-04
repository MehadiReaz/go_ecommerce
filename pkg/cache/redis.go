package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"ecommerce_project/internal/config"
)

var client *redis.Client

// NewConnection creates a new Redis connection
func NewConnection(cfg config.RedisConfig) (*redis.Client, error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}

// Get returns the Redis client instance
func Get() *redis.Client {
	return client
}

// Close closes the Redis connection
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}

// Set sets a value in Redis with expiration
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value from Redis
func GetValue(ctx context.Context, key string) (string, error) {
	return client.Get(ctx, key).Result()
}

// Delete deletes a key from Redis
func Delete(ctx context.Context, keys ...string) error {
	return client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists in Redis
func Exists(ctx context.Context, keys ...string) (int64, error) {
	return client.Exists(ctx, keys...).Result()
}
