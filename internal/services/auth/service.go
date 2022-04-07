package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/uesleicarvalhoo/go-auth-service/internal/domain/entity"
	"github.com/uesleicarvalhoo/go-auth-service/internal/infra/config"
	"github.com/uesleicarvalhoo/go-auth-service/internal/schemas"
	"github.com/uesleicarvalhoo/go-auth-service/pkg/auth"
	"github.com/uesleicarvalhoo/go-auth-service/pkg/trace"
)

const (
	tokenDuration                 = time.Minute * 15
	recoveryPasswordTokenDuration = time.Hour * 24
)

type Service struct {
	secretKey    string
	eventChannel chan schemas.Event
	userService  UserService
	cacheService CacheService
}

func NewService(userSvc UserService, cacheSvc CacheService, secretKey string, eventCh chan schemas.Event) *Service {
	return &Service{
		eventChannel: eventCh,
		userService:  userSvc,
		cacheService: cacheSvc,
	}
}

func (s Service) GenerateToken(
	ctx context.Context, userID uuid.UUID, prefix TokenPrefix, duration time.Duration,
) (schemas.JwtToken, error) {
	token, err := auth.GenerateJwtToken(s.secretKey, userID, duration)
	if err != nil {
		return schemas.JwtToken{}, err
	}

	err = s.cacheService.Set(
		ctx, fmt.Sprintf("%s-%s", prefix, userID.String()), token.Token, duration,
	)
	if err != nil {
		return schemas.JwtToken{}, err
	}

	return token, nil
}

func (s Service) validateToken(ctx context.Context, token string, prefix TokenPrefix) (uuid.UUID, error) {
	ctx, span := trace.NewSpan(ctx, "validate-token")
	defer span.End()

	userID, err := auth.ValidateJwtToken(token, s.secretKey)
	if err != nil {
		return uuid.UUID{}, errors.Wrap(ErrNotAuthorized, "Invalid Token")
	}

	cachedToken, _ := s.cacheService.Get(ctx, fmt.Sprintf("%s-%s", prefix, userID.String()))
	if cachedToken == "" || cachedToken != "" && cachedToken != token {
		return uuid.UUID{}, errors.Wrap(ErrNotAuthorized, "Token not found")
	}

	user, err := s.userService.Get(ctx, userID)
	if err != nil {
		return uuid.UUID{}, err
	}

	return user.ID, nil
}

func (s Service) invalidateToken(ctx context.Context, userID uuid.UUID, prefix TokenPrefix) error {
	return s.cacheService.Del(ctx, fmt.Sprintf("%s-%s", prefix, userID.String()))
}

func (s Service) ValidateAccessToken(ctx context.Context, token string) (uuid.UUID, error) {
	return s.validateToken(ctx, token, AccessTokenPrefix)
}

func (s Service) SignUp(ctx context.Context, payload schemas.SignUp) (*entity.User, error) {
	ctx, span := trace.NewSpan(ctx, "sign-up")
	defer span.End()

	if u, _ := s.userService.GetByEmail(ctx, payload.Email); u.Email == payload.Email {
		return nil, ErrEmailIsAlreadyUsed
	}

	user, err := entity.NewUser(payload.Name, payload.Email, payload.Phone, payload.Password)
	if err != nil {
		return nil, err
	}

	err = s.userService.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s Service) Login(ctx context.Context, payload schemas.Login) (schemas.LoginResponse, error) {
	ctx, span := trace.NewSpan(ctx, "login")
	defer span.End()

	user, err := s.userService.GetByEmail(ctx, payload.Email)
	if err != nil {
		return schemas.LoginResponse{
			Message: "Invalid email",
		}, err
	}

	if !user.ValidatePassword(payload.Password) {
		return schemas.LoginResponse{
			Message: "Invalid password",
		}, ErrNotAuthorized
	}

	if !user.Active {
		return schemas.LoginResponse{
			Message: "User is inactive",
		}, ErrNotAuthorized
	}

	accessToken, err := s.GenerateToken(ctx, user.ID, AccessTokenPrefix, tokenDuration)
	if err != nil {
		return schemas.LoginResponse{
			Message: "Error on generate access token",
		}, err
	}

	refreshToken, err := s.GenerateToken(ctx, user.ID, RefreshAcessTokenPrefix, config.SessionTime)
	if err != nil {
		return schemas.LoginResponse{
			Message: "Error on generate refresh access token",
		}, err
	}

	go s.sendEvent(
		"login", map[string]string{"user_id": user.ID.String(), "logged_at": time.Now().Format(time.RFC3339Nano)},
	)

	return schemas.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s Service) RefreshAccessToken(
	ctx context.Context, refreshToken schemas.RefreshToken,
) (schemas.LoginResponse, error) {
	ctx, span := trace.NewSpan(ctx, "refresh-token")
	defer span.End()

	userID, err := s.validateToken(ctx, refreshToken.Token, RefreshAcessTokenPrefix)
	if err != nil {
		return schemas.LoginResponse{
			Message: "Invalid token",
		}, err
	}

	newAccessToken, err := s.GenerateToken(ctx, userID, AccessTokenPrefix, tokenDuration)
	if err != nil {
		return schemas.LoginResponse{
			Message: "Error on generate access token",
		}, err
	}

	newRefreshToken, err := s.GenerateToken(ctx, userID, RefreshAcessTokenPrefix, config.SessionTime)
	if err != nil {
		return schemas.LoginResponse{
			Message: "Error on generate refresh access token",
		}, err
	}

	return schemas.LoginResponse{AccessToken: newAccessToken, RefreshToken: newRefreshToken}, nil
}

func (s Service) Logout(ctx context.Context, id uuid.UUID) error {
	ctx, span := trace.NewSpan(ctx, "logout")
	defer span.End()

	return s.invalidateToken(ctx, id, AccessTokenPrefix)
}

func (s Service) SendRecoveryPasswordToken(ctx context.Context, payload schemas.SendRecoveryPasswordPayload) error {
	ctx, span := trace.NewSpan(ctx, "notify-recovery-password-token")
	defer span.End()

	user, err := s.userService.GetByEmail(ctx, payload.Email)
	if err != nil {
		return err
	}

	token, err := s.GenerateToken(ctx, user.ID, RecoveryTokenPrefix, recoveryPasswordTokenDuration)
	if err != nil {
		return err
	}

	go s.sendEvent("recovery-password", map[string]any{
		"user":           map[string]string{"id": user.ID.String(), "name": user.Name, "email": user.Email},
		"recovery_token": token.Token,
		"expires_at":     token.ExpiresAt,
	})

	return nil
}

func (s Service) RecoveryPassword(ctx context.Context, token, password string) error {
	ctx, span := trace.NewSpan(ctx, "recovery-password")
	defer span.End()

	userID, err := s.validateToken(ctx, token, RecoveryTokenPrefix)
	if err != nil {
		return err
	}

	_, err = s.userService.Update(ctx, userID, schemas.UpdateUserPayload{Password: password})
	if err != nil {
		return err
	}

	return s.invalidateToken(ctx, userID, RecoveryTokenPrefix)
}

func (s Service) sendEvent(action string, data interface{}) {
	if body, err := json.Marshal(data); err == nil {
		s.eventChannel <- schemas.Event{Service: "authentication", Action: action, Data: body}
	}
}
