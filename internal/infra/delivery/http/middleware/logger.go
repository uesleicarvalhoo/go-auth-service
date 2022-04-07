package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uesleicarvalhoo/go-auth-service/internal/infra/config"
	"github.com/uesleicarvalhoo/go-auth-service/pkg/logger"
)

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.URL.Path, "/health-check") {
			return
		}

		c.Next()
		statusCode := c.Writer.Status()
		entry := logger.WithFields(logger.Fields{
			"log_version": "1.0.0",
			"date_time":   time.Now(),
			"product": map[string]interface{}{
				"name":        config.ServiceName,
				"application": config.ServiceName,
				"version":     config.ServiceVersion,
				"http": map[string]string{
					"method": c.Request.Method,
					"path":   c.Request.URL.Path,
				},
			},
			"origin": map[string]interface{}{
				"application": config.ServiceName,
				"ip":          c.ClientIP(),
				"headers": map[string]string{
					"user_agent": c.Request.UserAgent(),
					"origin":     c.GetHeader("Origin"),
					"refer":      c.Request.Referer(),
				},
			},
			"context": map[string]interface{}{
				"service":     config.ServiceName,
				"status_code": statusCode,
				"request_id":  "",
				"user_id":     c.GetHeader(config.HeaderUserID),
			},
		})

		switch {
		case len(c.Errors) > 0:
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		case statusCode >= http.StatusInternalServerError:
			entry.Error()
		case statusCode >= http.StatusBadRequest:
			entry.Warn()
		default:
			entry.Info()
		}
	}
}
