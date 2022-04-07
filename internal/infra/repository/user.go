package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/uesleicarvalhoo/go-auth-service/internal/domain/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (ur UserRepository) Get(ctx context.Context, id uuid.UUID) (entity.User, error) {
	var user entity.User

	tx := ur.DB.WithContext(ctx).WithContext(ctx).First(&user, "id = ?", id)

	return user, tx.Error
}

func (ur UserRepository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User

	tx := ur.DB.WithContext(ctx).First(&user, "email = ?", email)

	return user, tx.Error
}

func (ur UserRepository) Create(ctx context.Context, user entity.User) error {
	tx := ur.DB.WithContext(ctx).Create(&user).Omit("updated_at")

	return tx.Error
}

func (ur UserRepository) Update(ctx context.Context, user entity.User) error {
	tx := ur.DB.WithContext(ctx).Save(user)

	return tx.Error
}

func (ur UserRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	tx := ur.DB.WithContext(ctx).Delete(&entity.User{}, "id = ?", id)

	return tx.Error
}
