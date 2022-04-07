package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uesleicarvalhoo/go-auth-service/internal/infra/config"
)

// HealthCheck godoc
// @Summary  Return status of Service
// @Tags     General
// @Accept   json
// @Produce  json
// @Success  200  {object}  handler.MessageJSON
// @Router   /api/health-check [get].
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":    config.ServiceName,
		"message": "Server running",
		"version": config.ServiceVersion,
	})
}
