package redisclient

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

var (
	ErrInvalidMode = errors.New("invalid redis mode")
	sfGroup        singleflight.Group
)

type RedisClient struct {
	redis redis.UniversalClient
}

// New creates a new Redis client based on the specified mode.
// mode can be "single" for a single Redis instance or "cluster" for a Redis cluster.
func New(mode string, addresses []string, password string, db int) (*RedisClient, error) {
	var client redis.UniversalClient

	switch mode {
	case "single":
		client = redis.NewClient(&redis.Options{
			Addr:     addresses[0],
			Password: password,
			DB:       db,
		})
	case "cluster":
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    addresses,
			Password: password,
		})
	default:
		return nil, ErrInvalidMode
	}

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{redis: client}, nil
}

func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	log.Printf("RedisClient Set: key=%s", key)
	return c.redis.Set(ctx, key, value, expiration).Err()
}

func (c *RedisClient) Get(ctx context.Context, key string) (string, error) {
	log.Printf("RedisClient Get: key=%s", key)
	return c.redis.Get(ctx, key).Result()
}

func (c *RedisClient) Del(ctx context.Context, key string) error {
	log.Printf("RedisClient Del: key=%s", key)
	return c.redis.Del(ctx, key).Err()
}

// GetWithFallback tries to get value from Redis. If not found, it uses fallback to load the data,
// and caches it. Uses singleflight.DoChan to avoid cache penetration.
func (c *RedisClient) GetWithFallback(ctx context.Context, key string, ttl time.Duration, fallback func() (string, error)) (string, error) {
	// Step 1: Try from Redis
	val, err := c.Get(ctx, key)
	if err == nil {
		return val, nil
	}
	if err != redis.Nil {
		return "", err
	}

	// Step 2: Use singleflight to avoid cache penetration
	ch := sfGroup.DoChan(key, func() (interface{}, error) {
		// fallback to DB or other data source
		val, err := fallback()
		if err != nil {
			return "", err
		}
		// cache result
		_ = c.redis.Set(ctx, key, val, ttl).Err()
		return val, nil
	})

	// Step 3: Wait for result
	res := <-ch
	if res.Err != nil {
		return "", res.Err
	}
	return res.Val.(string), nil
}
