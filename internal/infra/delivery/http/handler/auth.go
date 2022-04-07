package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uesleicarvalhoo/go-auth-service/internal/infra/config"
	"github.com/uesleicarvalhoo/go-auth-service/internal/schemas"
	"github.com/uesleicarvalhoo/go-auth-service/pkg/trace"
)

// SignUp godoc
// @Summary  Register new user true
// @Param    payload  body  schemas.SignUp  true  "User data"
// @Tags     Auth
// @Accept   json
// @produce  json
// @Success  201  {object}  entity.User
// @Failure  400  {object}  handler.MessageJSON
// @Failure  422  {object}  handler.MessageJSON
// @Router   /api/v1/auth/signup [post].
func (h *Handler) SignUp(c *gin.Context) {
	ctx, span := trace.NewSpan(c.Request.Context(), "handler.signup")
	defer span.End()

	var payload schemas.SignUp
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, MessageJSON{Message: "Invalid Payload"})
		trace.AddSpanError(span, err)
		trace.FailSpan(span, "Unprocessable entity")

		return
	}

	user, err := h.AuthSvc.SignUp(ctx, payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, MessageJSON{Message: err.Error()})
		trace.AddSpanError(span, err)
		trace.FailSpan(span, "Bad request")

		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary      Get user access token
// @Description  Generate a new access token
// @Param        payload  body  schemas.Login  true  "User data"
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  schemas.LoginResponse
// @Failure      401  {object}  handler.MessageJSON
// @Failure      422  {object}  handler.MessageJSON
// @Router       /api/v1/auth/login [post].
func (h *Handler) Login(c *gin.Context) {
	ctx, span := trace.NewSpan(c.Request.Context(), "handler.login")
	defer span.End()

	var payload schemas.Login
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, MessageJSON{Message: "Invalid Payload"})
		trace.AddSpanError(span, err)
		trace.FailSpan(span, "Unprocessable entity")

		return
	}

	res, err := h.AuthSvc.Login(ctx, payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, MessageJSON{Message: res.Message})
		trace.AddSpanError(span, err)
		trace.FailSpan(span, "Login Unauthorized")

		return
	}

	c.JSON(http.StatusOK, res)
}

// Logout godoc
// @Summary      Logout user
// @Description  Logout current user and expire access token
// @Param        Authorization  header  string  true  "Acess token"
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  handler.MessageJSON
// @Failure      401  {object}  handler.MessageJSON
// @Router       /api/v1/auth/logout [post].
func (h *Handler) Logout(c *gin.Context) {
	ctx, span := trace.NewSpan(c.Request.Context(), "handler.logout")
	defer span.End()

	ctxUserID, _ := c.Get("userID")
	userID, _ := ctxUserID.(uuid.UUID)

	trace.AddSpanTags(span, map[string]string{"user_id": userID.String()})

	if err := h.AuthSvc.Logout(ctx, userID); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		trace.AddSpanError(span, err)
		trace.FailSpan(span, "logout error")

		return
	}

	c.JSON(http.StatusOK, MessageJSON{Message: "Success"})
}

// RefreshAccessToken godoc
// @Summary      Refresh user access token
// @Description  Return a new access token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  schemas.JwtToken
// @Failure      401  {object}  handler.MessageJSON
// @Failure      500  {object}  handler.MessageJSON
// @Router       /api/v1/auth/refresh-access-token [post].
func (h *Handler) RefreshAccessToken(c *gin.Context) {
	ctx, span := trace.NewSpan(c.Request.Context(), "handler.refresh-access-token")
	defer span.End()

	var payload schemas.RefreshToken
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, MessageJSON{Message: "Invalid Payload"})
		trace.AddSpanError(span, err)
		trace.FailSpan(span, "Unprocessable entity")

		return
	}

	newToken, err := h.AuthSvc.RefreshAccessToken(ctx, payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, MessageJSON{Message: err.Error()})
		trace.AddSpanError(span, err)
		trace.FailSpan(span, "Unauthorized")

		return
	}

	c.JSON(http.StatusOK, newToken)
}

// Authorize godoc
// @Summary      Check user authentication
// @Description  Check if acess token is valid
// @Param        Authorization  body  schemas.AutorizationPayload  true  "Acess token"
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  handler.MessageJSON
// @Failure      401  {object}  handler.MessageJSON
// @Failure      422  {object}  handler.MessageJSON
// @Router       /api/v1/auth/authorize [post].
func (h *Handler) Authorize(c *gin.Context) {
	ctx, span := trace.NewSpan(c.Request.Context(), "handler.authorize")
	defer span.End()

	var payload schemas.AuthorizationPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, MessageJSON{Message: "Invalid payload"})
		trace.AddSpanError(span, err)
		trace.FailSpan(span, "Unprocessable entity")

		return
	}

	userID, err := h.AuthSvc.ValidateAccessToken(ctx, payload.Token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, MessageJSON{Message: "Invalid token"})
		trace.AddSpanError(span, err)
		trace.FailSpan(span, "Invalid token")

		return
	}

	trace.AddSpanTags(span, map[string]string{"user_id": userID.String()})
	c.Header(config.HeaderUserID, userID.String())
	c.JSON(http.StatusOK, MessageJSON{Message: "ok"})
}

// Login godoc
// @Summary      Send recovery password token
// @Param        payload  body  schemas.SendRecoveryPasswordPayload  true  "User data"
// @Description  Send new token for password recovery
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      201  {object}  handler.MessageJSON
// @Failure      422  {object}  handler.MessageJSON
// @Failure      500  {object}  handler.MessageJSON
// @Router       /api/v1/auth/recovery-password [post].
func (h *Handler) SendRecoveryPasswordToken(c *gin.Context) {
	ctx, span := trace.NewSpan(c.Request.Context(), "handler.send-recovery-password-token")
	defer span.End()

	var payload schemas.SendRecoveryPasswordPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, MessageJSON{Message: "Invalid payload"})
		trace.AddSpanError(span, err)
		trace.FailSpan(span, "Unprocessable entity")

		return
	}

	err := h.AuthSvc.SendRecoveryPasswordToken(ctx, payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, MessageJSON{Message: err.Error()})
		trace.AddSpanError(span, err)
		trace.FailSpan(span, "Error on send recovery password token")

		return
	}

	c.JSON(http.StatusAccepted, MessageJSON{Message: "Accepted"})
}

// Login godoc
// @Summary      Reset password
// @Description  Change user password
// @Param        payload  body  schemas.SendRecoveryPasswordPayload  true  "Recovery data"
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  handler.MessageJSON
// @Success      401  {object}  handler.MessageJSON
// @Failure      422  {object}  handler.MessageJSON
// @Router       /api/v1/auth/reset-password [post].
func (h *Handler) ResetPassword(c *gin.Context) {
	ctx, span := trace.NewSpan(c.Request.Context(), "handler.reset-password")
	defer span.End()

	var payload schemas.RecoveryPassword
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, MessageJSON{Message: "Invalid Payload"})
		trace.AddSpanError(span, err)
		trace.FailSpan(span, "Unprocessable entity")

		return
	}

	err := h.AuthSvc.RecoveryPassword(ctx, payload.Token, payload.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, MessageJSON{Message: err.Error()})
		trace.AddSpanError(span, err)
		trace.FailSpan(span, "Error on reset password")

		return
	}

	c.JSON(http.StatusOK, MessageJSON{Message: "Password updated"})
}
