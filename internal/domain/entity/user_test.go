package entity_test

import (
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v4"
	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/go-auth-service/internal/domain/entity"
	"github.com/uesleicarvalhoo/go-auth-service/internal/schemas"
)

func TestNewUserErrorMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario         string
		name             string
		email            string
		phone            string
		password         string
		expectedMessages []string
	}{
		{
			scenario:         "when name is empty",
			name:             "",
			email:            gofakeit.Email(),
			phone:            gofakeit.Phone(),
			password:         gofakeit.Password(true, true, true, true, true, 10),
			expectedMessages: []string{"name: field is required"},
		},
		{
			scenario:         "when email is empty",
			name:             gofakeit.Name(),
			email:            "",
			phone:            gofakeit.Phone(),
			password:         gofakeit.Password(true, true, true, true, true, 10),
			expectedMessages: []string{"email: mail: no address"},
		},
		{
			scenario:         "when phone is empty",
			name:             gofakeit.Name(),
			email:            gofakeit.Email(),
			phone:            "",
			password:         gofakeit.Password(true, true, true, true, true, 10),
			expectedMessages: []string{"phone: '' is not a valid phone number"},
		},
		{
			scenario:         "when password is empty",
			name:             gofakeit.Name(),
			email:            gofakeit.Email(),
			phone:            gofakeit.Phone(),
			password:         "",
			expectedMessages: []string{"password: password must be 5 or more caracters"},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			_, err := entity.NewUser(tc.name, tc.email, tc.phone, tc.password)
			assert.Error(t, err)

			for _, expectedMessage := range tc.expectedMessages {
				assert.ErrorContains(t, err, expectedMessage)
			}
		})
	}
}

func TestNewUser(t *testing.T) {
	t.Parallel()

	user, err := entity.NewUser(
		gofakeit.Name(),
		gofakeit.Email(),
		gofakeit.Phone(),
		gofakeit.Password(true, true, true, true, true, 10),
	)

	assert.Nil(t, err)

	assert.True(t, user.Active)
	assert.True(t, user.DeletedAt.IsZero())
	assert.True(t, user.UpdatedAt.IsZero())
	assert.False(t, user.CreatedAt.IsZero())
}

func TestUserUpdate(t *testing.T) {
	t.Parallel()

	user, err := entity.NewUser(
		gofakeit.Name(),
		gofakeit.Email(),
		gofakeit.Phone(),
		gofakeit.Password(true, true, true, true, true, 10),
	)
	assert.Nil(t, err)

	oldEmail := user.Email
	oldPhone := user.Phone
	oldPasswordHash := user.PasswordHash

	newName := gofakeit.Name()
	newEmail := gofakeit.Email()
	newPassword := gofakeit.Password(true, true, true, true, true, 10)
	newPhone := gofakeit.Phone()

	err = user.Update(schemas.UpdateUserPayload{Name: newName, Email: newEmail, Phone: newPhone, Password: newPassword})
	assert.Nil(t, err)
	assert.Equal(t, user.Name, strings.ToTitle(newName))
	assert.NotEqual(t, user.Email, oldEmail)
	assert.NotEqual(t, user.Phone, oldPhone)
	assert.NotEqual(t, user.PasswordHash, newPassword)
	assert.NotEqual(t, user.PasswordHash, oldPasswordHash)

	assert.Nil(t, err)
}

func TestUserUpdateSetUpdatedAt(t *testing.T) {
	t.Parallel()

	user, err := entity.NewUser(
		gofakeit.Name(),
		gofakeit.Email(),
		gofakeit.Phone(),
		gofakeit.Password(true, true, true, true, true, 10),
	)
	assert.Nil(t, err)

	err = user.Update(schemas.UpdateUserPayload{Name: "New user name"})
	assert.Nil(t, err)

	assert.False(t, user.UpdatedAt.IsZero())
}

func TestUserValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario       string
		user           entity.User
		shouldHasError bool
	}{
		{
			scenario:       "when name is empty",
			shouldHasError: true,
			user: entity.User{
				Email:        gofakeit.Email(),
				Phone:        gofakeit.Phone(),
				PasswordHash: gofakeit.Password(true, true, true, true, true, 10),
			},
		},
		{
			scenario:       "when email is empty",
			shouldHasError: true,
			user: entity.User{
				Name:         gofakeit.Name(),
				Phone:        gofakeit.Phone(),
				PasswordHash: gofakeit.Password(true, true, true, true, true, 10),
			},
		},
		{
			scenario:       "when phone is empty",
			shouldHasError: true,
			user: entity.User{
				Name:         gofakeit.Name(),
				Email:        gofakeit.Email(),
				PasswordHash: gofakeit.Password(true, true, true, true, true, 10),
			},
		},
		{
			scenario:       "when password hash is empty",
			shouldHasError: true,
			user: entity.User{
				Name:  gofakeit.Name(),
				Email: gofakeit.Email(),
				Phone: gofakeit.Phone(),
			},
		},
		{
			scenario:       "when all fields are ok",
			shouldHasError: false,
			user: entity.User{
				Name:         gofakeit.Name(),
				Email:        gofakeit.Email(),
				Phone:        gofakeit.Phone(),
				PasswordHash: gofakeit.Password(true, true, true, true, true, 10),
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()
			err := tc.user.Validate()

			if tc.shouldHasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserValidatePassword(t *testing.T) {
	t.Parallel()

	password := "iamsecret"

	user, err := entity.NewUser(
		gofakeit.Name(),
		gofakeit.Email(),
		gofakeit.Phone(),
		password,
	)
	assert.Nil(t, err)

	assert.True(t, user.ValidatePassword(password))
	assert.False(t, user.ValidatePassword("wrongpassword"))
}
