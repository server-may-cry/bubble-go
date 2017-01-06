package controllers

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

// LoadStatick load static from file server (crutch for spend less money and not store static files in repo)
func LoadStatick(c *gin.Context) {
	filePath := c.Param("filePath")
	c.String(http.StatusNotFound, "text/plain", "TODO", filePath)
}

// ClearStatickCache remove statick files
func ClearStatickCache(c *gin.Context) {
	c.String(http.StatusOK, "text/plain", "TODO")
}
