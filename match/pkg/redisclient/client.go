package redisclient

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrInvalidMode = errors.New("invalid redis mode")

type Client struct {
	redis redis.UniversalClient
}

func New(mode string, addresses []string, password string, db int) (*Client, error) {
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

	return &Client{redis: client}, nil
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.redis.Set(ctx, key, value, expiration).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.redis.Get(ctx, key).Result()
}

func (c *Client) Del(ctx context.Context, key string) error {
	return c.redis.Del(ctx, key).Err()
}
