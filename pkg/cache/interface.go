package cache

import (
	"context"
	"time"

	"github.com/uesleicarvalhoo/go-auth-service/internal/infra/config"
)

type Config = config.CacheConfig

type Client interface {
	HSet(ctx context.Context, key string, values ...any) error
	Set(ctx context.Context, key, value string, expiration time.Duration) error
	HSetExp(ctx context.Context, key string, expiration time.Duration, values ...any) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}
