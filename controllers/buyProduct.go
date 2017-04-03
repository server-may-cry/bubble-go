package controllers

import (
	"net/http"

	"github.com/server-may-cry/bubble-go/market"
	"github.com/server-may-cry/bubble-go/models"
	"github.com/server-may-cry/bubble-go/storage"

	"gopkg.in/gin-gonic/gin.v1"
)

type buyProductRequest struct {
	baseRequest
	ProductID string `json:"productId" binding:"required"`
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
	market.Buy(&user, request.ProductID)
	response := buyProductResponse{
		ProductID: request.ProductID,
		Credits:   user.Credits,
	}
	storage.Gorm.Save(&user)
	c.JSON(http.StatusOK, response)
}
