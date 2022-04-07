package auth

import "errors"

var (
	ErrNotAuthorized      = errors.New("not authorized")
	ErrEmailIsAlreadyUsed = errors.New("the email is already being used")
)
