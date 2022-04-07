package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uesleicarvalhoo/go-auth-service/internal/schemas"
	"github.com/uesleicarvalhoo/go-auth-service/pkg/trace"
)

// SignUp godoc
// @Summary  Get current user data
// @Param    Authorization  header  string  true  "Bearer token"
// @Tags     User
// @Accept   json
// @produce  json
// @Success  200  {object}  entity.User
// @Failure  400  {object}  handler.MessageJSON
// @Failure  401  {object}  handler.MessageJSON
// @Router   /api/v1/user/me [get].
func (h *Handler) GetMe(c *gin.Context) {
	ctx, span := trace.NewSpan(c.Request.Context(), "handler.get-me")
	defer span.End()

	ctxUserID, _ := c.Get("userID")
	userID, _ := ctxUserID.(uuid.UUID)

	trace.AddSpanTags(span, map[string]string{"user_id": userID.String()})

	user, err := h.UserSvc.Get(ctx, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, MessageJSON{Message: "Couldn't find the user"})
		trace.AddSpanError(span, err)
		trace.FailSpan(span, "Failed to load user")

		return
	}

	c.JSON(http.StatusOK, user)
}

// SignUp godoc
// @Summary  Update current user data
// @Param    Authorization  header  string  true  "Bearer token"
// @Tags     User
// @Accept   json
// @produce  json
// @Success  200  {object}  entity.User
// @Failure  400  {object}  handler.MessageJSON
// @Failure  401  {object}  handler.MessageJSON
// @Failure  422  {object}  handler.MessageJSON
// @Router   /api/v1/user/me [post].
func (h *Handler) UpdateMe(c *gin.Context) {
	ctx, span := trace.NewSpan(c.Request.Context(), "handler.update-me")
	defer span.End()

	ctxUserID, _ := c.Get("userID")
	userID, _ := ctxUserID.(uuid.UUID)

	trace.AddSpanTags(span, map[string]string{"user_id": userID.String()})

	var payload schemas.UpdateUserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, MessageJSON{Message: "Invalid Payload"})
		trace.AddSpanError(span, err)
		trace.FailSpan(span, "Unprocessable entity")

		return
	}

	user, err := h.UserSvc.Update(ctx, userID, payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, MessageJSON{Message: "Failed to update user"})
		trace.AddSpanError(span, err)
		trace.FailSpan(span, "Failed to update user data")

		return
	}

	c.JSON(http.StatusOK, user)
}

// SignUp godoc
// @Summary  Delete current user
// @Param    Authorization  header  string  true  "Bearer token"
// @Tags     User
// @Accept   json
// @produce  json
// @Success  200  {object}  handler.MessageJSON
// @Failure  401  {object}  handler.MessageJSON
// @Failure  500  {object}  handler.MessageJSON
// @Router   /api/v1/user/me [delete].
func (h *Handler) DeleteMe(c *gin.Context) {
	ctx, span := trace.NewSpan(c.Request.Context(), "handler.delete-me")
	defer span.End()

	ctxUserID, _ := c.Get("userID")
	userID, _ := ctxUserID.(uuid.UUID)

	err := h.UserSvc.Delete(ctx, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, MessageJSON{Message: "Failed to delete user"})
		trace.AddSpanError(span, err)
		trace.FailSpan(span, "Failed to delete user")

		return
	}

	c.JSON(http.StatusOK, MessageJSON{Message: "ok"})
}
