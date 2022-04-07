package validators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	validators "github.com/uesleicarvalhoo/go-auth-service/pkg/utils/validator"
)

func TestNormalizeEmailReturnErrorWhenEmailIsInvalid(t *testing.T) {
	t.Parallel()

	email := "@email.com"

	_, err := validators.NormalizeEmail(email)
	assert.NotNil(t, err)
}

func TestNormalizeEmail(t *testing.T) {
	t.Parallel()

	email := "user@email.com"

	_, err := validators.NormalizeEmail(email)
	assert.Nil(t, err)
}
