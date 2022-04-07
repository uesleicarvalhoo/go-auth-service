package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uesleicarvalhoo/go-auth-service/internal/infra/config"
	"github.com/uesleicarvalhoo/go-auth-service/pkg/trace"
)

func AuthenticationMiddlware(service TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := trace.NewSpan(c.Request.Context(), "Middleware.Authentication")

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "Authorization not found"})
			trace.FailSpan(span, "Unauthorized")

			return
		}

		token := authHeader[len(config.TokenSchema)+1:]

		userID, err := service.ValidateAccessToken(ctx, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
			trace.AddSpanError(span, err)
			trace.FailSpan(span, "Unauthorized")

			return
		}

		c.Header(config.HeaderUserID, userID.String())
		c.Set("userID", userID)
	}
}
