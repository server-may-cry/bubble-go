package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/server-may-cry/bubble-go/models"
	"github.com/server-may-cry/bubble-go/platforms"
	"github.com/server-may-cry/bubble-go/storage"
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
	UserID                 uint       `json:"userId"`
	ReachedStage01         int8       `json:"reachedStage01,uint8"`
	ReachedStage02         int8       `json:"reachedStage02,uint8"`
	ReachedSubStage01      int8       `json:"reachedSubStage01,uint8"`
	ReachedSubStage02      int8       `json:"reachedSubStage02,uint8"`
	IgnoreSavePointBlock   int8       `json:"ignoreSavePointBlock,bool"`
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
	OfferAvailable         uint8      `json:"offerAvailable"` // bool
	FirstGame              uint8      `json:"firstGame"`      // bool
	StagesProgressStat01   [8]uint32  `json:"stagesProgressStat01"`
	StagesProgressStat02   [8]uint32  `json:"stagesProgressStat02"`
	SubStagesRecordStats01 [8][]uint8 `json:"subStagesRecordStats01"`
	SubStagesRecordStats02 [8][]uint8 `json:"subStagesRecordStats02"`
}

// ReqEnter first request from client. Return user info and user progress
func ReqEnter(w http.ResponseWriter, r *http.Request) {
	request := enterRequest{}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, getErrBody(err), http.StatusBadRequest)
		return
	}
	var firstGame uint8 // bool
	ctx := r.Context()
	value := ctx.Value(User)
	var user models.User
	switch value.(type) {
	case models.User:
		user = value.(models.User)
		// time rewards logic
	case nil:
		firstGame = 1
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
			RestoreTriesAt:          0,
			Credits:                 0, // TODO
			FriendsBonusCreditsTime: time.Now().Unix(),
			// TODO ProgressStandart:        [][]int8 // json
		}
		storage.Gorm.Create(&user) // Gorm.NewRecord check row exists or somehow
	}
	log.Println(user)
	// TODO logic
	// TODO add users progress
	response := enterResponse{
		ReqMsgID:             request.MsgID,
		UserID:               user.ID,
		ReachedStage01:       user.ReachedStage01,
		ReachedStage02:       user.ReachedStage02,
		ReachedSubStage01:    user.ReachedSubStage01,
		ReachedSubStage02:    user.ReachedSubStage02,
		IgnoreSavePointBlock: user.IgnoreSavePointBlock,
		RemainingTries:       user.RemainingTries,
		Credits:              user.Credits,
		InfinityExtra00:      user.InifinityExtra00,
		InfinityExtra01:      user.InifinityExtra01,
		InfinityExtra02:      user.InifinityExtra02,
		InfinityExtra03:      user.InifinityExtra03,
		InfinityExtra04:      user.InifinityExtra04,
		InfinityExtra05:      user.InifinityExtra05,
		InfinityExtra06:      user.InifinityExtra06,
		InfinityExtra07:      user.InifinityExtra07,
		InfinityExtra08:      user.InifinityExtra08,
		InfinityExtra09:      user.InifinityExtra09,
		OfferAvailable:       0,
		FirstGame:            firstGame,
		// BonusCredits   TODO
		// AppFriendsBonusCredits TODO
		// StagesProgressStat01 TODO
		// StagesProgressStat02   TODO
		// SubStagesRecordStats01 TODO
		// SubStagesRecordStats02 TODO
	}
	JSON(w, response)
}
