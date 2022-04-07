package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uesleicarvalhoo/go-auth-service/internal/domain/entity"
	"github.com/uesleicarvalhoo/go-auth-service/internal/schemas"
)

type TokenPrefix string

const (
	AccessTokenPrefix       TokenPrefix = "acess-token"
	RefreshAcessTokenPrefix TokenPrefix = "refresh-acess-token"
	RecoveryTokenPrefix     TokenPrefix = "recovery-token"
)

type UserService interface {
	Get(ctx context.Context, id uuid.UUID) (user entity.User, err error)
	GetByEmail(ctx context.Context, email string) (user entity.User, err error)
	Create(ctx context.Context, user entity.User) error
	Update(ctx context.Context, id uuid.UUID, payload schemas.UpdateUserPayload) (*entity.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CacheService interface {
	HSet(ctx context.Context, key string, values ...any) error
	Set(ctx context.Context, key, value string, expiration time.Duration) error
	HSetExp(ctx context.Context, key string, expiration time.Duration, values ...any) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}
