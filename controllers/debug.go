package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/server-may-cry/bubble-go/models"
)

var i = 0

// Index control index route
func Index(c *gin.Context) {
	i = i + 1
	c.JSON(http.StatusOK, gin.H{
		"status":  "posted",
		"message": "msg",
		"i":       i,
	})
}

// Test used for debug
func Test(c *gin.Context) {
	t := models.Test{
		Ttt: 4,
	}
	log.Println(t)
	c.JSON(http.StatusOK, gin.H{
		"test": t,
	})
}
