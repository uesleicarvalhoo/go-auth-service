package validators

import (
	"net/mail"
)

func NormalizeEmail(email string) (string, error) {
	e, err := mail.ParseAddress(email)
	if err != nil {
		return "", err
	}

	return e.Address, nil
}
