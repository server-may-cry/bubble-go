package controllers

import (
	"net/http"

	"github.com/server-may-cry/bubble-go/models"
	"github.com/server-may-cry/bubble-go/storage"

	"gopkg.in/gin-gonic/gin.v1"
)

type reduceCreditsRequest struct {
	baseRequest
	Amount uint16 `json:"amount,string" binding:"required"`
}

type reduceCreditsResponse struct {
	ReqMsgID string `json:"reqMsgId"`
	UserID   string `json:"userId"`
	Credits  uint16 `json:"credits"`
}

// ReqReduceCredits reduce user credits. Get amount from request
func ReqReduceCredits(c *gin.Context) {
	request := reduceCreditsRequest{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, getErrBody(err))
		return
	}
	user := c.MustGet("user").(models.User)
	user.Credits -= request.Amount
	storage.Gorm.Save(&user)
	response := reduceCreditsResponse{
		Credits: user.Credits,
	}
	c.JSON(http.StatusOK, response)
}
