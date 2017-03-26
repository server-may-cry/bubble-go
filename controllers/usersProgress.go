package controllers

import (
	"net/http"

	"github.com/server-may-cry/bubble-go/models"
	"github.com/server-may-cry/bubble-go/storage"

	"gopkg.in/gin-gonic/gin.v1"
)

type usersProgressRequest struct {
	baseRequest
	SocIDs []uint64 `json:"socIds,[]uint64" binding:"required"`
}

type userProgres struct {
	UserID            uint64 `json:"userId"`
	SocID             string `json:"socId"`
	ReachedStage01    uint8  `json:"reachedStage01"`
	ReachedStage02    uint8  `json:"reachedStage02"`
	ReachedSubStage01 uint8  `json:"reachedSubStage01"`
	ReachedSubStage02 uint8  `json:"reachedSubStage02"`
}
type usersProgressResponse struct {
	UsersProgress []userProgres `json:"usersProgress"`
}

// ReqUsersProgress return progres of recieved users
func ReqUsersProgress(c *gin.Context) {
	request := usersProgressRequest{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, getErrBody(err))
		return
	}
	user := c.MustGet("user").(models.User)
	usersLen := len(request.SocIDs)
	users := make([]models.User, usersLen)
	storage.Gorm.Find(&users, "sysId = ? and extId in (?)", user.SysID, request.SocIDs)
	response := usersProgressResponse{
		UsersProgress: make([]userProgres, usersLen),
	}
	for i, friend := range users {
		response.UsersProgress[i] = userProgres{
			UserID:            friend.ID,
			SocID:             friend.ExtID,
			ReachedStage01:    friend.ReachedStage01,
			ReachedSubStage01: friend.ReachedSubStage01,
		}
	}
	c.JSON(http.StatusOK, response)
}
