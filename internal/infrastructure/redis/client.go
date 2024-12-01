package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrNotFound = errors.New("key not found")
	ErrRedis    = errors.New("redis error")
)

type Client interface {
	Set(ctx context.Context, key string, value interface{}, exp time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	Close() error
	Ping(ctx context.Context) error
}

type client struct {
	cl *redis.Client
}

func NewClient(addr, password string, db int) Client {
	return &client{
		cl: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "",
			DB:       db,
		}),
	}
}

func (c *client) Set(ctx context.Context, key string, value interface{}, exp time.Duration) error {
	err := c.Ping(ctx)
	if err != nil {
		return ErrRedis
	}
	
	err = c.cl.Set(ctx, key, value, exp).Err()
	if err != nil {
		return ErrRedis
	}
	return nil
}

func (c *client) Get(ctx context.Context, key string) (string, error) {
	err := c.Ping(ctx)
	if err != nil {
		return "", ErrRedis
	}
	
	res, err := c.cl.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrNotFound
		} else {
			return "", ErrRedis
		}
	}

	return res, nil
}

func (c *client) Del(ctx context.Context, keys ...string) error {
	err := c.Ping(ctx)
	if err != nil {
		return ErrRedis
	}
	
	err = c.cl.Del(ctx, keys...).Err()
	if err != nil {
		if err == redis.Nil {
			return ErrNotFound
		} else {
			return ErrRedis
		}
	}
	return nil
}

func (c *client) Close() error {
	err := c.cl.Close()
	if err != nil {
		return ErrRedis
	}
	return nil
}

func (c *client) Ping(ctx context.Context) error {
	err := c.cl.Ping(ctx).Err()
	if err != nil {
		return ErrRedis
	}
	return nil
}
