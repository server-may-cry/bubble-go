package controllers

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

// PayPlatform acept and validate payment request from platforms
func PayPlatform(c *gin.Context) {
	platform := c.Param("platform")
	c.String(http.StatusNotFound, "text/plain", "TODO", platform)
}
