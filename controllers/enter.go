package controllers

import (
	"log"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

type enterRequest struct {
	baseRequest
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
	log.Print("before bind")
	if err := c.BindJSON(&request); err != nil {
		log.Print(err)
		request2 := baseRequest{}
		log.Print("before bind2")
		if err2 := c.BindJSON(&request2); err2 != nil {
			log.Print(err2)
			request3 := AuthRequestPart{}
			log.Print("before bind3")
			if err3 := c.BindJSON(&request3); err3 != nil {
				log.Print(err3)
			}
		}
		c.JSON(http.StatusBadRequest, getErrBody(err))
		return
	}
	// logic
	response := enterResponse{
		ReqMsgID: request.MsgID,
	}
	c.JSON(http.StatusOK, response)
}
