package auth_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/go-auth-service/internal/domain/entity"
	"github.com/uesleicarvalhoo/go-auth-service/internal/infra/repository"
	"github.com/uesleicarvalhoo/go-auth-service/internal/schemas"
	"github.com/uesleicarvalhoo/go-auth-service/internal/services/auth"
	"github.com/uesleicarvalhoo/go-auth-service/internal/services/user"
	"github.com/uesleicarvalhoo/go-auth-service/pkg/cache"
	"github.com/uesleicarvalhoo/go-auth-service/pkg/database"
)

const (
	testSecretKey = "my-test-secret-key"
)

type Sut struct {
	service      *auth.Service
	cache        auth.CacheService
	userSvc      auth.UserService
	userRepo     *repository.UserRepository
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

	cacheClient, err := cache.NewMemoryCacheClient()
	if err != nil {
		panic(err)
	}

	userRepository := repository.NewUserRepository(db)
	userService := user.NewService(userRepository, eventChannel)

	service := auth.NewService(userService, cacheClient, testSecretKey, eventChannel)

	return Sut{
		service:      service,
		cache:        cacheClient,
		userSvc:      userService,
		userRepo:     userRepository,
		eventChannel: eventChannel,
	}
}

func TestSiginUp(t *testing.T) {
	t.Parallel()

	existingUserEmail := gofakeit.Email()

	tests := []struct {
		scenario          string
		existingUserEmail string
		payload           schemas.SignUp
		expectedError     string
	}{
		{
			scenario: "when name is empty",
			payload: schemas.SignUp{
				Email:    gofakeit.Email(),
				Phone:    gofakeit.Phone(),
				Password: gofakeit.Password(true, true, true, true, true, 10),
			},
			expectedError: "name: field is required",
		},
		{
			scenario: "when email is empty",
			payload: schemas.SignUp{
				Name:     gofakeit.Name(),
				Phone:    gofakeit.Phone(),
				Password: gofakeit.Password(true, true, true, true, true, 10),
			},
			expectedError: "email: mail: no address",
		},
		{
			scenario: "when phone is empty",
			payload: schemas.SignUp{
				Name:     gofakeit.Name(),
				Email:    gofakeit.Email(),
				Password: gofakeit.Password(true, true, true, true, true, 10),
			},
			expectedError: "phone: '' is not a valid phone number",
		},
		{
			scenario: "when password is empty",
			payload: schemas.SignUp{
				Name:  gofakeit.Name(),
				Email: gofakeit.Email(),
				Phone: gofakeit.Phone(),
			},
			expectedError: "password: password must be 5 or more caracters",
		},
		{
			scenario: "when password contains 4 caracters",
			payload: schemas.SignUp{
				Name:     gofakeit.Name(),
				Email:    gofakeit.Email(),
				Phone:    gofakeit.Phone(),
				Password: "1234",
			},
			expectedError: "password: password must be 5 or more caracters",
		},
		{
			scenario:          "when user email is in use",
			existingUserEmail: existingUserEmail,
			payload: schemas.SignUp{
				Name:     gofakeit.Name(),
				Email:    existingUserEmail,
				Phone:    gofakeit.Phone(),
				Password: gofakeit.Password(true, true, true, true, true, 10),
			},
			expectedError: "the email is already being used",
		},
		{
			scenario: "when all fields are ok",
			payload: schemas.SignUp{
				Name:     gofakeit.Name(),
				Email:    gofakeit.Email(),
				Phone:    gofakeit.Phone(),
				Password: gofakeit.Password(true, true, true, true, true, 10),
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			// Arrange
			sut := newSut()

			if tc.existingUserEmail != "" {
				existingUser, err := entity.NewUser(
					gofakeit.Name(),
					existingUserEmail,
					gofakeit.Phone(),
					gofakeit.Password(true, true, true, true, true, 10),
				)
				assert.Nil(t, err)

				err = sut.userSvc.Create(context.TODO(), existingUser)
				assert.NoError(t, err)
			}

			// Action
			user, err := sut.service.SignUp(context.TODO(), tc.payload)

			// Assert
			if tc.expectedError != "" {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				repoUser, err := sut.userRepo.Get(context.TODO(), user.ID)

				assert.NoError(t, err)
				assert.Equal(t, user.ID, repoUser.ID)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	t.Parallel()

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, true, 10)

	tests := []struct {
		scenario        string
		payload         schemas.Login
		expectedError   string
		expectedMessage string
	}{
		{
			scenario:        "when email is empty",
			expectedError:   "record not found",
			expectedMessage: "Invalid email",
			payload: schemas.Login{
				Email:    "",
				Password: password,
			},
		},
		{
			scenario:        "when password is empty",
			expectedError:   "not authorized",
			expectedMessage: "Invalid password",
			payload: schemas.Login{
				Email:    email,
				Password: "",
			},
		},
		{
			scenario:        "when password is invalid",
			expectedError:   "not authorized",
			expectedMessage: "Invalid password",
			payload: schemas.Login{
				Email:    email,
				Password: "wrong-password",
			},
		},
		{
			scenario: "when email and passsword is ok",
			payload: schemas.Login{
				Email:    email,
				Password: password,
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		email := email
		password := password

		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			// Arrange
			sut := newSut()

			signUp := schemas.SignUp{
				Email:    email,
				Password: password,
				Name:     gofakeit.Name(),
				Phone:    gofakeit.Phone(),
			}

			user, err := sut.service.SignUp(context.TODO(), signUp)
			assert.NoError(t, err)
			<-sut.eventChannel

			// Action
			response, err := sut.service.Login(context.TODO(), tc.payload)

			// Assert
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
				assert.Equal(t, tc.expectedMessage, response.Message)
			} else {
				// Login
				loggedUserID, err := sut.service.ValidateAccessToken(context.TODO(), response.AccessToken.Token)
				assert.NoError(t, err)
				assert.Equal(t, user.ID, loggedUserID)

				// Event
				event := <-sut.eventChannel
				assert.Equal(t, "login", event.Action)
				assert.Equal(t, "authentication", event.Service)
			}
		})
	}
}

func TestValidateAccessToken(t *testing.T) {
	t.Parallel()

	signUp := schemas.SignUp{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Phone:    gofakeit.Phone(),
		Password: gofakeit.Password(true, true, true, true, true, 10),
	}

	sut := newSut()

	user, err := sut.service.SignUp(context.TODO(), signUp)
	assert.NoError(t, err)
	<-sut.eventChannel

	loginResponse, err := sut.service.Login(context.TODO(), schemas.Login{
		Email:    signUp.Email,
		Password: signUp.Password,
	})
	assert.NoError(t, err)

	tests := []struct {
		scenario      string
		userID        uuid.UUID
		accessToken   string
		expectedError string
	}{
		{
			scenario:      "when token is empty",
			expectedError: "Invalid Token: not authorized",
		},
		{
			userID:      user.ID,
			scenario:    "when token is valid",
			accessToken: loginResponse.AccessToken.Token,
		},
		{
			scenario:      "when token is not a acessToken",
			expectedError: "Token not found: not authorized",
			accessToken:   loginResponse.RefreshToken.Token,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			// Action
			tokenUserID, err := sut.service.ValidateAccessToken(context.TODO(), tc.accessToken)

			// Assert
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.userID, tokenUserID)
			}
		})
	}
}

func TestRefreshAccessToken(t *testing.T) {
	t.Parallel()

	signUp := schemas.SignUp{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Phone:    gofakeit.Phone(),
		Password: gofakeit.Password(true, true, true, true, true, 10),
	}

	sut := newSut()

	_, err := sut.service.SignUp(context.TODO(), signUp)
	assert.NoError(t, err)
	<-sut.eventChannel

	loginResponse, err := sut.service.Login(context.TODO(), schemas.Login{
		Email:    signUp.Email,
		Password: signUp.Password,
	})
	assert.NoError(t, err)

	tests := []struct {
		scenario      string
		accessToken   string
		refreshToken  string
		expectedError string
	}{
		{
			scenario:      "when token is empty should return an error",
			expectedError: "Invalid Token: not authorized",
		},
		{
			scenario:     "when token is valid should return a valid access token and invalidate the last one",
			accessToken:  loginResponse.AccessToken.Token,
			refreshToken: loginResponse.RefreshToken.Token,
		},
		{
			scenario:      "when token is not a refreshToken should return an error",
			expectedError: "Token not found: not authorized",
			refreshToken:  loginResponse.AccessToken.Token,
			accessToken:   loginResponse.AccessToken.Token,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			// Arrange
			payload := schemas.RefreshToken{
				JwtToken: schemas.JwtToken{
					Token: tc.refreshToken,
				},
			}

			// Action
			time.Sleep(time.Second) // Wait 1 second to change accessToken hash
			response, err := sut.service.RefreshAccessToken(context.TODO(), payload)

			// Assert
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				_, err := sut.service.ValidateAccessToken(context.TODO(), response.AccessToken.Token)
				assert.NoError(t, err)

				_, err = sut.service.ValidateAccessToken(context.TODO(), tc.accessToken)
				assert.Error(t, err)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario      string
		userID        uuid.UUID
		expectedError string
	}{
		{
			scenario: "when logout is success, accessToken must be invalidated",
			userID:   uuid.New(),
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			// Arrange
			sut := newSut()

			accessToken, err := sut.service.GenerateToken(context.TODO(), tc.userID, auth.AccessTokenPrefix, time.Hour)
			assert.NoError(t, err)

			// Action
			err = sut.service.Logout(context.TODO(), tc.userID)

			// Assert
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				_, err := sut.service.ValidateAccessToken(context.TODO(), accessToken.Token)
				assert.Error(t, err)
				assert.EqualError(t, err, "Token not found: not authorized")
			}
		})
	}
}

func TestSendRecoveryPasswordToken(t *testing.T) {
	t.Parallel()

	email := gofakeit.Email()

	tests := []struct {
		scenario          string
		expectedError     string
		existingUserEmail string
		payload           schemas.SendRecoveryPasswordPayload
	}{
		{
			scenario:          "when email does not exist",
			existingUserEmail: email,
			payload: schemas.SendRecoveryPasswordPayload{
				Email: email,
			},
		},
		{
			scenario: "when email is empty",
			payload: schemas.SendRecoveryPasswordPayload{
				Email: "",
			},
			expectedError: "record not found",
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			// Arrange
			sut := newSut()

			if tc.existingUserEmail != "" {
				user, err := entity.NewUser(
					gofakeit.Name(),
					tc.existingUserEmail,
					gofakeit.Phone(),
					gofakeit.Password(true, true, true, true, true, 10),
				)
				assert.NoError(t, err)

				err = sut.userSvc.Create(context.TODO(), user)
				assert.NoError(t, err)
				<-sut.eventChannel
			}

			// Action
			err := sut.service.SendRecoveryPasswordToken(context.TODO(), tc.payload)

			// Assert
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)

				event := <-sut.eventChannel
				assert.Equal(t, event.Service, "authentication")
				assert.Equal(t, event.Action, "recovery-password")

				eventData := map[string]any{}
				err = json.Unmarshal(event.Data, &eventData)
				assert.NoError(t, err)

				eventUser, ok := eventData["user"].(map[string]any)
				assert.True(t, ok)

				eventUserEmail, ok := eventUser["email"].(string)
				assert.True(t, ok)
				assert.Equal(t, tc.payload.Email, eventUserEmail)

				_, ok = eventData["recovery_token"]
				assert.True(t, ok)
			}
		})
	}
}

func TestRecoveryPassword(t *testing.T) {
	t.Parallel()

	// Prepare
	tests := []struct {
		scenario      string
		user          *entity.User
		newPassword   string
		recoveryToken string
		expectedError string
		createUser    bool
	}{
		{
			scenario:      "when token is invalid",
			recoveryToken: "invalid-token",
			expectedError: "Invalid Token: not authorized",
		},
		{
			scenario:      "when token is valid and password is empty",
			expectedError: "no data for update",
			createUser:    true,
		},
		{
			scenario:    "when token is valid and password is valid",
			newPassword: gofakeit.Password(true, true, true, true, true, 10),
			createUser:  true,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			// Arrange
			sut := newSut()

			if tc.createUser {
				user, err := sut.service.SignUp(context.Background(), schemas.SignUp{
					Name:     gofakeit.Name(),
					Email:    gofakeit.Email(),
					Phone:    gofakeit.Phone(),
					Password: gofakeit.Password(true, true, true, true, true, 10),
				})
				assert.NoError(t, err)

				assert.NoError(t, err)

				token, err := sut.service.GenerateToken(context.TODO(), user.ID, auth.RecoveryTokenPrefix, time.Hour)
				assert.NoError(t, err)

				tc.user = user
				tc.recoveryToken = token.Token
			}

			// Action
			err := sut.service.RecoveryPassword(context.TODO(), tc.recoveryToken, tc.newPassword)

			// Assert
			if tc.expectedError == "" {
				assert.NoError(t, err)

				updatedUser, err := sut.userRepo.Get(context.TODO(), tc.user.ID)
				assert.NoError(t, err)

				assert.Equal(t, tc.user.ID, updatedUser.ID)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
