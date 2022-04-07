package entity

import (
	"time"

	"github.com/google/uuid"
)

type AccessHistory struct {
	UserID   uuid.UUID `json:"user_id"`
	LoggedAt time.Time `json:"logged_at"`
}
