package controllers

import (
	"log"
	"net/http"

	"github.com/server-may-cry/bubble-go/models"

	"gopkg.in/gin-gonic/gin.v1"
)

type reduceTriesRequest struct {
	baseRequest
}

// ReqReduceTries reduce user tries by one
func ReqReduceTries(c *gin.Context) {
	request := reduceTriesRequest{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, getErrBody(err))
		return
	}
	user := c.MustGet("user").(models.User)
	log.Print(user)
	// logic
	// Response [4]
	response := make([]uint8, 1)
	response[0] = 123
	c.JSON(http.StatusOK, response)
}
