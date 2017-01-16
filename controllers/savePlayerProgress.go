package controllers

import (
	"log"
	"net/http"

	"github.com/server-may-cry/bubble-go/models"

	"gopkg.in/gin-gonic/gin.v1"
)

type savePlayerProgressRequest struct {
	baseRequest
	Amount uint16 `json:"productId,string" binding:"required"`
}

// ReqSavePlayerProgress save player progress
func ReqSavePlayerProgress(c *gin.Context) {
	request := savePlayerProgressRequest{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, getErrBody(err))
		return
	}
	user := c.MustGet("user").(models.User)
	log.Print(user)
	// logic
	response := "ok"
	c.JSON(http.StatusOK, response)
}
