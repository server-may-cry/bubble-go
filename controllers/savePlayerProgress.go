package controllers

import (
	"net/http"

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
	// logic
	response := "ok"
	c.JSON(http.StatusOK, response)
}
