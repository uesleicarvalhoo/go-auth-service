package user

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/uesleicarvalhoo/go-auth-service/internal/domain/entity"
	"github.com/uesleicarvalhoo/go-auth-service/internal/schemas"
	"github.com/uesleicarvalhoo/go-auth-service/pkg/trace"
)

type Service struct {
	eventChannel chan schemas.Event
	repository   Repository
}

func NewService(repository Repository, eventChannel chan schemas.Event) *Service {
	return &Service{repository: repository, eventChannel: eventChannel}
}

func (s Service) Get(ctx context.Context, id uuid.UUID) (entity.User, error) {
	ctx, span := trace.NewSpan(ctx, "user.get")
	defer span.End()

	return s.repository.Get(ctx, id)
}

func (s Service) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	ctx, span := trace.NewSpan(ctx, "user.get-by-email")
	defer span.End()

	return s.repository.GetByEmail(ctx, email)
}

func (s Service) Create(ctx context.Context, user entity.User) error {
	ctx, span := trace.NewSpan(ctx, "user.create")
	defer span.End()

	err := s.repository.Create(ctx, user)
	if err != nil {
		return err
	}

	go s.sendEvent("create", user)

	return nil
}

func (s Service) Update(ctx context.Context, id uuid.UUID, payload schemas.UpdateUserPayload) (*entity.User, error) {
	ctx, span := trace.NewSpan(ctx, "user.update")
	defer span.End()

	user, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = user.Update(payload)
	if err != nil {
		return nil, err
	}

	err = s.repository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	go s.sendEvent("update", user)

	return &user, nil
}

func (s Service) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := trace.NewSpan(ctx, "user.delete")
	defer span.End()

	err := s.repository.DeleteByID(ctx, id)
	if err != nil {
		return err
	}

	go s.sendEvent("delete", map[string]string{"user_id": id.String(), "deleted_at": time.Now().String()})

	return nil
}

func (s Service) sendEvent(action string, data interface{}) {
	if body, err := json.Marshal(data); err == nil {
		s.eventChannel <- schemas.Event{Service: "user", Action: action, Data: body}
	}
}
