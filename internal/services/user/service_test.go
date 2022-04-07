package user_test

import (
	"context"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/go-auth-service/internal/domain/entity"
	"github.com/uesleicarvalhoo/go-auth-service/internal/infra/repository"
	"github.com/uesleicarvalhoo/go-auth-service/internal/schemas"
	"github.com/uesleicarvalhoo/go-auth-service/internal/services/user"
	"github.com/uesleicarvalhoo/go-auth-service/pkg/database"
)

type Sut struct {
	service      *user.Service
	repository   *repository.UserRepository
	eventChannel chan schemas.Event
}

func newSut() Sut {
	eventChannel := make(chan schemas.Event)

	db, err := database.NewSQLiteMemoryConnection()
	if err != nil {
		panic(err)
	}

	err = repository.AutoMigrate(db)
	if err != nil {
		panic(err)
	}

	userRepository := repository.NewUserRepository(db)

	return Sut{
		repository:   userRepository,
		eventChannel: eventChannel,
		service:      user.NewService(userRepository, eventChannel),
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario      string
		user          entity.User
		expectedError string
	}{
		{
			scenario: "should save user in repository",
			user: entity.User{
				ID:           uuid.New(),
				Name:         gofakeit.Name(),
				Email:        gofakeit.Email(),
				Phone:        gofakeit.Phone(),
				PasswordHash: "fake-hash",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			// Arrange
			sut := newSut()

			// Action
			err := sut.service.Create(context.TODO(), tc.user)

			// Assert
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)

				userByID, err := sut.repository.Get(context.TODO(), tc.user.ID)
				assert.NoError(t, err)
				assert.Equal(t, tc.user.ID, userByID.ID)

				userByEmail, err := sut.repository.GetByEmail(context.TODO(), tc.user.Email)
				assert.NoError(t, err)
				assert.Equal(t, tc.user.Email, userByEmail.Email)

				// Event
				event := <-sut.eventChannel

				assert.Equal(t, event.Action, "create")
				assert.Equal(t, event.Service, "user")
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario      string
		userID        uuid.UUID
		createUser    bool
		payload       schemas.UpdateUserPayload
		expectedError string
	}{
		{
			scenario:      "when id is invalid",
			expectedError: "record not found",
		},
		{
			scenario:      "when payload is empty",
			expectedError: "no data for update",
			createUser:    true,
		},
		{
			scenario: "when has data for update",
			payload: schemas.UpdateUserPayload{
				Name:     gofakeit.Name(),
				Email:    gofakeit.Email(),
				Phone:    gofakeit.Phone(),
				Password: gofakeit.Password(true, true, true, true, true, 10),
			},
			createUser: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			// Arrange
			sut := newSut()

			if tc.createUser {
				user := entity.User{
					ID:           uuid.New(),
					Name:         gofakeit.Name(),
					Email:        gofakeit.Email(),
					Phone:        gofakeit.Phone(),
					PasswordHash: "fake-hash",
				}

				err := sut.repository.Create(context.TODO(), user)
				assert.NoError(t, err)

				tc.userID = user.ID
			}

			// Action
			updatedUser, err := sut.service.Update(context.TODO(), tc.userID, tc.payload)

			// Assert
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, strings.ToTitle(tc.payload.Name), updatedUser.Name)
				assert.Equal(t, tc.payload.Email, updatedUser.Email)
				assert.Equal(t, tc.payload.Phone, updatedUser.Phone)
				assert.True(t, updatedUser.ValidatePassword(tc.payload.Password))

				// Event
				event := <-sut.eventChannel

				assert.Equal(t, event.Action, "update")
				assert.Equal(t, event.Service, "user")
			}
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario      string
		user          entity.User
		userID        uuid.UUID
		expectedError string
	}{
		{
			scenario: "when id not exist in repository",
			user: entity.User{
				ID:           uuid.New(),
				Name:         gofakeit.Name(),
				Email:        gofakeit.Email(),
				Phone:        gofakeit.Phone(),
				PasswordHash: "fake-hash",
			},
			userID:        uuid.New(),
			expectedError: "record not found",
		},
		{
			scenario: "when id exist in repository",
			userID:   uuid.UUID{},
			user: entity.User{
				ID:           uuid.UUID{},
				Name:         gofakeit.Name(),
				Email:        gofakeit.Email(),
				Phone:        gofakeit.Phone(),
				PasswordHash: "fake-hash",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			// Arrange
			sut := newSut()
			err := sut.service.Create(context.TODO(), tc.user)
			assert.NoError(t, err)

			// Action
			returnedUser, err := sut.service.Get(context.TODO(), tc.userID)

			// Assert
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tc.user.ID, returnedUser.ID)
			}
		})
	}
}

func TestGetByEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario      string
		user          entity.User
		email         string
		expectedError string
	}{
		{
			scenario: "when email not exist in repository",
			user: entity.User{
				ID:           uuid.New(),
				Name:         gofakeit.Name(),
				Email:        "user@email.com",
				Phone:        gofakeit.Phone(),
				PasswordHash: "fake-hash",
			},
			email:         "notexisting@email.com",
			expectedError: "record not found",
		},
		{
			scenario: "when id exist in repository",
			email:    "user@email.com",
			user: entity.User{
				ID:           uuid.New(),
				Name:         gofakeit.Name(),
				Email:        "user@email.com",
				Phone:        gofakeit.Phone(),
				PasswordHash: "fake-hash",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			// Arrange
			sut := newSut()
			err := sut.service.Create(context.TODO(), tc.user)
			assert.NoError(t, err)

			// Action
			returnedUser, err := sut.service.GetByEmail(context.TODO(), tc.email)

			// Assert
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tc.user.Email, returnedUser.Email)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario      string
		userID        uuid.UUID
		user          entity.User
		expectedError string
	}{
		{
			scenario: "when id exist in repository",
			userID:   uuid.UUID{},
			user: entity.User{
				ID:           uuid.UUID{},
				Name:         gofakeit.Name(),
				Email:        gofakeit.Email(),
				Phone:        gofakeit.Phone(),
				PasswordHash: "fake-hash",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			// Arrange
			sut := newSut()
			err := sut.service.Create(context.TODO(), tc.user)
			assert.NoError(t, err)
			<-sut.eventChannel

			// Action
			err = sut.service.Delete(context.TODO(), tc.userID)

			// Assert
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)

				_, err := sut.service.Get(context.TODO(), tc.userID)
				assert.Error(t, err)
				assert.EqualError(t, err, "record not found")

				event := <-sut.eventChannel
				assert.Equal(t, event.Action, "delete")
				assert.Equal(t, event.Service, "user")
			}
		})
	}
}
