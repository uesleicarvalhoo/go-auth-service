package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const MinPasswordLength = 5

func GeneratePasswordHash(password string) (string, error) {
	if len(password) < MinPasswordLength {
		return "", fmt.Errorf("password must be %d or more caracters", MinPasswordLength)
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}

func CheckPasswordHash(plain, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))

	return err == nil
}
