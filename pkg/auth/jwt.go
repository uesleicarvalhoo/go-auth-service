package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/uesleicarvalhoo/go-auth-service/internal/schemas"
)

func GenerateJwtToken(secret string, userID uuid.UUID, expiration time.Duration) (schemas.JwtToken, error) {
	exp := time.Now().Add(expiration).Unix()
	claims := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.StandardClaims{
			Subject:   userID.String(),
			ExpiresAt: exp,
		})

	token, err := claims.SignedString([]byte(secret))
	if err != nil {
		return schemas.JwtToken{}, nil
	}

	return schemas.JwtToken{Token: token, ExpiresAt: exp}, err
}

func ValidateJwtToken(token, secret string) (uuid.UUID, error) {
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	claims, ok := tokenObj.Claims.(jwt.MapClaims)
	if ok && tokenObj.Valid {
		sub, _ := claims["sub"].(string)

		return uuid.Parse(sub)
	}

	return uuid.UUID{}, nil
}
