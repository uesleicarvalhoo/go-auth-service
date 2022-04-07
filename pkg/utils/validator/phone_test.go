package validators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	validators "github.com/uesleicarvalhoo/go-auth-service/pkg/utils/validator"
)

func TestNormalizePhoneNumber(t *testing.T) {
	t.Parallel()

	phone := "5511999999999"
	_, err := validators.NormalizePhoneNumber(phone)
	assert.Nil(t, err)
}

func TestNormalizePhoneNumberReturnErrorWhenPhoneIsInvalid(t *testing.T) {
	t.Parallel()

	phone := ""
	_, err := validators.NormalizePhoneNumber(phone)
	assert.NotNil(t, err)
}
