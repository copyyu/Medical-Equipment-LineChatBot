package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck returns server health status
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "medical-equipment-webhook",
		"version": "1.0.0",
	})
}
