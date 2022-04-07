package schemas

type JwtToken struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiration"`
}

type RefreshToken struct {
	JwtToken
}
