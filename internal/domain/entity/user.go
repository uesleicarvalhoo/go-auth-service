package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uesleicarvalhoo/go-auth-service/internal/schemas"
	"github.com/uesleicarvalhoo/go-auth-service/pkg/auth"
	validators "github.com/uesleicarvalhoo/go-auth-service/pkg/utils/validator"
)

type User struct {
	ID           uuid.UUID `json:"id" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	Email        string    `json:"email" binding:"required"`
	Phone        string    `json:"phone" binding:"required"`
	PasswordHash string    `json:"-" binding:"required"`
	Active       bool      `json:"active" default:"true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

func (u *User) Validate() error {
	validator := NewValidator()

	if strings.TrimSpace(u.Name) == "" {
		validator.AddError("name", "field is required")
	}

	if u.PasswordHash == "" {
		validator.AddError("password", "field is required")
	}

	if u.Email == "" {
		validator.AddError("email", "field is required")
	}

	if u.Phone == "" {
		validator.AddError("phone", "field is required")
	}

	if validator.HasErrors() {
		return validator.GetError()
	}

	return nil
}

func (u *User) Update(payload schemas.UpdateUserPayload) error {
	if payload == (schemas.UpdateUserPayload{}) {
		return ErrInvalidData
	}

	if payload.Name != "" {
		u.Name = strings.ToTitle(payload.Name)
	}

	if payload.Email != "" {
		email, err := validators.NormalizeEmail(payload.Email)
		if err != nil {
			return err
		}

		u.Email = email
	}

	if payload.Phone != "" {
		u.Phone, _ = validators.NormalizePhoneNumber(payload.Phone)
	}

	if payload.Password != "" {
		passwordHash, err := auth.GeneratePasswordHash(payload.Password)
		if err != nil {
			return err
		}

		u.PasswordHash = passwordHash
	}

	u.UpdatedAt = time.Now()

	return u.Validate()
}

func (u *User) ValidatePassword(password string) bool {
	return auth.CheckPasswordHash(password, u.PasswordHash)
}

func NewUser(name, email, phone, password string) (User, error) {
	validator := NewValidator()

	phone, err := validators.NormalizePhoneNumber(phone)
	if err != nil {
		validator.AddError("phone", err.Error())
	}

	passwordHash, err := auth.GeneratePasswordHash(password)
	if err != nil {
		validator.AddError("password", err.Error())
	}

	email, err = validators.NormalizeEmail(email)
	if err != nil {
		validator.AddError("email", err.Error())
	}

	if strings.TrimSpace(name) == "" {
		validator.AddError("name", "field is required")
	}

	if phone == "" {
		validator.AddError("phone", "field is required")
	}

	if validator.HasErrors() {
		return User{}, validator.GetError()
	}

	return User{
		ID:           uuid.New(),
		Name:         strings.ToTitle(name),
		Email:        email,
		Phone:        phone,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		Active:       true,
	}, nil
}
