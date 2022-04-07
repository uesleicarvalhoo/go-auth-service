package schemas

type Login struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  JwtToken `json:"access_token,omitempty"`
	RefreshToken JwtToken `json:"refresh_token,omitempty"`
	Message      string   `json:"message,omitempty"`
}

type SendRecoveryPasswordPayload struct {
	Email string `json:"email" binding:"required"`
}

type RecoveryPassword struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthorizationPayload struct {
	Token string `json:"token" binding:"required"`
}
