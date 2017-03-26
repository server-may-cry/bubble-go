package controllers

import (
	"net/http"

	"github.com/server-may-cry/bubble-go/models"
	"github.com/server-may-cry/bubble-go/storage"

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
	user.RemainingTries--
	storage.Gorm.Save(&user)
	response := make([]uint8, 1)
	response[0] = user.RemainingTries
	c.JSON(http.StatusOK, response)
}
