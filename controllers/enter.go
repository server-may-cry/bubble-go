package controllers

import (
	"log"
	"net/http"

	"github.com/server-may-cry/bubble-go/models"
	"github.com/server-may-cry/bubble-go/platforms"

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
	ReachedStage01         uint8      `json:"reachedStage01"`
	ReachedStage02         uint8      `json:"reachedStage02"`
	ReachedSubStage01      uint8      `json:"reachedSubStage01"`
	ReachedSubStage02      uint8      `json:"reachedSubStage02"`
	IgnoreSavePointBlock   bool       `json:"ignoreSavePointBlock"`
	RemainingTries         uint8      `json:"remainingTries"`
	Credits                uint16     `json:"credits"`
	InfinityExtra00        bool       `json:"inifinityExtra00"`
	InfinityExtra01        bool       `json:"inifinityExtra01"`
	InfinityExtra02        bool       `json:"inifinityExtra02"`
	InfinityExtra03        bool       `json:"inifinityExtra03"`
	InfinityExtra04        bool       `json:"inifinityExtra04"`
	InfinityExtra05        bool       `json:"inifinityExtra05"`
	InfinityExtra06        bool       `json:"inifinityExtra06"`
	InfinityExtra07        bool       `json:"inifinityExtra07"`
	InfinityExtra08        bool       `json:"inifinityExtra08"`
	InfinityExtra09        bool       `json:"inifinityExtra09"`
	BonusCredits           uint16     `json:"bonusCredits"`
	AppFriendsBonusCredits uint16     `json:"appFriendsBonusCredits"`
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
	}
	log.Println(user)
	// logic
	// add users progress
	response := enterResponse{
		ReqMsgID: request.MsgID,
	}
	c.JSON(http.StatusOK, response)
}
