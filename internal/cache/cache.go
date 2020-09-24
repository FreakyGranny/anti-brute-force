package cache

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

//go:generate mockgen -source=$GOFILE -destination=../mocks/cache_mock.go -package=mocks Cache

// Cache cache object.
type Cache interface {
	Incr(ctx context.Context, key string, t time.Duration) error
	Get(ctx context.Context, key string) (int, error)
	Del(ctx context.Context, key string) error
	Close() error
}

// RedisCache redis cache storage.
type RedisCache struct {
	client *redis.Client
}

// New returns new redis cache instance.
func New(host string, port int, pass string) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(host, strconv.Itoa(port)),
		Password: pass,
		DB:       0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisCache{client: rdb}, nil
}

// Incr increment value by given key.
func (c *RedisCache) Incr(ctx context.Context, key string, t time.Duration) error {
	err := c.client.Incr(ctx, key).Err()
	if err != nil {
		return err
	}

	return c.client.Expire(ctx, key, t).Err()
}

// Get returns value of given key.
func (c *RedisCache) Get(ctx context.Context, key string) (int, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return 0, ErrKeyNotFound
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		// TODO fix error
		return 0, ErrKeyNotFound
	}

	return intVal, nil
}

// Del drops record by given key.
func (c *RedisCache) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Close closes connection.
func (c *RedisCache) Close() error {
	return c.client.Close()
}
