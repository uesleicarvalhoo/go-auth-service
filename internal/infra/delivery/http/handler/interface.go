package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/uesleicarvalhoo/go-auth-service/internal/domain/entity"
	"github.com/uesleicarvalhoo/go-auth-service/internal/schemas"
)

type AuthenticationService interface {
	RefreshAccessToken(ctx context.Context, refreshToken schemas.RefreshToken) (schemas.LoginResponse, error)
	ValidateAccessToken(ctx context.Context, token string) (uuid.UUID, error)
	SignUp(ctx context.Context, payload schemas.SignUp) (*entity.User, error)
	Login(ctx context.Context, payload schemas.Login) (schemas.LoginResponse, error)
	Logout(ctx context.Context, id uuid.UUID) error
	SendRecoveryPasswordToken(ctx context.Context, payload schemas.SendRecoveryPasswordPayload) error
	RecoveryPassword(ctx context.Context, token, password string) error
}

type UserService interface {
	Get(ctx context.Context, id uuid.UUID) (user entity.User, err error)
	GetByEmail(ctx context.Context, email string) (user entity.User, err error)
	Create(ctx context.Context, user entity.User) error
	Update(ctx context.Context, id uuid.UUID, payload schemas.UpdateUserPayload) (*entity.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type MessageJSON struct {
	Message string `json:"message"`
}
