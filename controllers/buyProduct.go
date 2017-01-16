package controllers

import (
	"log"
	"net/http"

	"github.com/server-may-cry/bubble-go/models"

	"gopkg.in/gin-gonic/gin.v1"
)

type buyProductRequest struct {
	baseRequest
	ProductID string `json:"productId" binding:"required"`

	LevelMode string `json:"levelMode" binding:"required"`
}

type buyProductResponse struct {
	ProductID string `json:"productId"`
	Credits   uint16 `json:"credits"`
}

// ReqBuyProduct buy product
func ReqBuyProduct(c *gin.Context) {
	request := buyProductRequest{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, getErrBody(err))
		return
	}
	user := c.MustGet("user").(models.User)
	log.Print(user)
	// logic
	response := buyProductResponse{
		ProductID: request.ProductID,
		Credits:   0, // TODO
	}
	c.JSON(http.StatusOK, response)
}
