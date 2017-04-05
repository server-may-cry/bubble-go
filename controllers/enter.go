package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/server-may-cry/bubble-go/models"
	"github.com/server-may-cry/bubble-go/platforms"
	"github.com/server-may-cry/bubble-go/storage"

	"gopkg.in/gin-gonic/gin.v1"
)

type enterRequest struct {
	baseRequest
	AuthRequestPart
	AppFriends uint8  `json:"appFriends,string" binding:"required"`
	Referer    string `json:"referer" binding:"required"`
	SrcExtID   string `json:"srcExtId" binding:"required"`
}

type enterResponse struct {
	ReqMsgID               uint64     `json:"reqMsgId"`
	UserID                 uint32     `json:"userId"`
	ReachedStage01         int8       `json:"reachedStage01,uint8"`
	ReachedStage02         int8       `json:"reachedStage02,uint8"`
	ReachedSubStage01      int8       `json:"reachedSubStage01,uint8"`
	ReachedSubStage02      int8       `json:"reachedSubStage02,uint8"`
	IgnoreSavePointBlock   bool       `json:"ignoreSavePointBlock"`
	RemainingTries         int8       `json:"remainingTries,uint8"`
	Credits                int16      `json:"credits,uint16"`
	InfinityExtra00        int8       `json:"inifinityExtra00,uint8"`
	InfinityExtra01        int8       `json:"inifinityExtra01,uint8"`
	InfinityExtra02        int8       `json:"inifinityExtra02,uint8"`
	InfinityExtra03        int8       `json:"inifinityExtra03,uint8"`
	InfinityExtra04        int8       `json:"inifinityExtra04,uint8"`
	InfinityExtra05        int8       `json:"inifinityExtra05,uint8"`
	InfinityExtra06        int8       `json:"inifinityExtra06,uint8"`
	InfinityExtra07        int8       `json:"inifinityExtra07,uint8"`
	InfinityExtra08        int8       `json:"inifinityExtra08,uint8"`
	InfinityExtra09        int8       `json:"inifinityExtra09,uint8"`
	BonusCredits           int16      `json:"bonusCredits,uint16"`
	AppFriendsBonusCredits int16      `json:"appFriendsBonusCredits,uint16"`
	OfferAvailable         bool       `json:"offerAvailable"`
	FirstGame              bool       `json:"firstGame"`
	StagesProgressStat01   [8]uint32  `json:"stagesProgressStat01"`
	StagesProgressStat02   [8]uint32  `json:"stagesProgressStat02"`
	SubStagesRecordStats01 [8][]uint8 `json:"subStagesRecordStats01"`
	SubStagesRecordStats02 [8][]uint8 `json:"subStagesRecordStats02"`
}

// ReqEnter first request from client. Return user info and user progress
func ReqEnter(c *gin.Context) {
	request := enterRequest{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, getErrBody(err))
		return
	}
	value, exists := c.Get("user")
	var user models.User
	if exists {
		user = value.(models.User)
		// time rewards logic
	} else {
		platformID := platforms.GetByName(request.SysID)
		user = models.User{
			SysID:                   platformID,
			ExtID:                   request.ExtID,
			ReachedStage01:          0,
			ReachedSubStage01:       0,
			IgnoreSavePointBlock:    0,
			InifinityExtra00:        0,
			InifinityExtra01:        0,
			InifinityExtra02:        0,
			InifinityExtra03:        0,
			InifinityExtra04:        0,
			InifinityExtra05:        0,
			InifinityExtra06:        0,
			InifinityExtra07:        0,
			InifinityExtra08:        0,
			InifinityExtra09:        0,
			RemainingTries:          5,
			RestoreTriesAt:          0, // TODO
			Credits:                 0, // TODO
			FriendsBonusCreditsTime: 0, // TODO
			// TODO ProgressStandart:        [][]int8 // json
		}
		success := storage.Gorm.NewRecord(&user)
		if !success {
			panic(fmt.Sprintf("can`t create user %v", user))
		}
	}
	log.Println(user)
	// logic
	// add users progress
	response := enterResponse{
		ReqMsgID: request.MsgID,
	}
	c.JSON(http.StatusOK, response)
}
