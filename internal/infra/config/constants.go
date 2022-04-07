package config

import "time"

const (
	ServiceName    = "go-auth-service"
	ServiceVersion = "0.0.0"
	SessionTime    = time.Hour * 1
)

const (
	HeaderUserID         = "X-User-ID"
	HeaderAuthentication = "Authorization"
	TokenSchema          = "Bearer"
)
