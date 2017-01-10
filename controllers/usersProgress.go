package controllers

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

type usersProgressRequest struct {
	baseRequest
	SocIDs []string `json:"socIds,[]uint64" binding:"required"`
}

type userProgress struct {
	UserID            uint64 `json:"userId"`
	SocID             string `json:"socId"`
	ReachedStage01    uint8  `json:"reachedStage01"`
	ReachedStage02    uint8  `json:"reachedStage02"`
	ReachedSubStage01 uint8  `json:"reachedSubStage01"`
	ReachedSubStage02 uint8  `json:"reachedSubStage02"`
}
type usersProgressResponse struct {
	ReqMsgID []userProgress `json:"usersProgress"`
}

// ReqUsersProgress return progres of recieved users
func ReqUsersProgress(c *gin.Context) {
	request := usersProgressRequest{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, getErrBody(err))
		return
	}
	// logic
	response := usersProgressResponse{}
	c.JSON(http.StatusOK, response)
}
