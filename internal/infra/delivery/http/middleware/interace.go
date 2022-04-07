package middleware

import (
	"context"

	"github.com/google/uuid"
)

type TokenService interface {
	ValidateAccessToken(ctx context.Context, token string) (uuid.UUID, error)
}
