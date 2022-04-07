package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cfg Config) (*RedisClient, error) {
	var client *redis.Client

	if cfg.Password != "" {
		client = redis.NewClient(&redis.Options{
			Addr:     cfg.GetURL(),
			DB:       0,
			Password: cfg.Password,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr: cfg.GetURL(),
			DB:   0,
		})
	}

	client.AddHook(redisotel.TracingHook{})

	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisClient{
		Client: client,
	}, nil
}

func (cache *RedisClient) HSet(ctx context.Context, key string, values ...any) error {
	return cache.Client.HSet(ctx, key, values).Err()
}

func (cache *RedisClient) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	return cache.Client.Set(ctx, key, value, expiration).Err()
}

func (cache *RedisClient) HSetExp(ctx context.Context, key string, expiration time.Duration, values ...any) error {
	pipe := cache.Client.Pipeline()
	pipe.HSet(ctx, key, values)
	pipe.Expire(ctx, key, expiration)
	_, err := pipe.Exec(ctx)

	return err
}

func (cache *RedisClient) Get(ctx context.Context, key string) (string, error) {
	data, err := cache.Client.Get(ctx, key).Result()
	if err != nil && err.Error() == "redis: nil" {
		return data, nil
	}

	return data, err
}

func (cache *RedisClient) Del(ctx context.Context, key string) error {
	res := cache.Client.Del(ctx, key)
	_, err := res.Result()

	return err
}
