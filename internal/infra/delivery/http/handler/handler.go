package handler

type Handler struct {
	AuthSvc AuthenticationService
	UserSvc UserService
}

func NewHandler(authenticationService AuthenticationService, userService UserService) *Handler {
	return &Handler{
		AuthSvc: authenticationService,
		UserSvc: userService,
	}
}
