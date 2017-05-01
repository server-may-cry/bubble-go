package application

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/server-may-cry/bubble-go/platforms"
)

type usersProgressRow struct {
	Cnt            uint32
	ReachedStage01 uint8
}

type enterRequest struct {
	baseRequest
	AuthRequestPart
	AppFriends uint8  `json:"appFriends,string"`
	Referer    string `json:"referer"`
	SrcExtID   string `json:"srcExtId"`
}

type enterResponse struct {
	ReqMsgID                  uint64 `json:"reqMsgId"`
	UserID                    uint   `json:"userId"`
	ReachedStage01            int8   `json:"reachedStage01,uint8"` // max user island
	ReachedStage02            int8   `json:"reachedStage02,uint8"`
	ReachedSubStage01         int8   `json:"reachedSubStage01,uint8"` // max user level on max island
	ReachedSubStage02         int8   `json:"reachedSubStage02,uint8"`
	IgnoreSavePointBlock      int8   `json:"ignoreSavePointBlock,bool"`
	RemainingTries            int8   `json:"remainingTries,uint8"`
	TriesMin                  int8   `json:"triesMin"`
	TriesRegenSecondsInterval int    `json:"triesRegenSecondsInterval"`
	SecondsUntilTriesRegen    int64  `json:"secondsUntilTriesRegen"`
	Credits                   int16  `json:"credits,uint16"`
	InfinityExtra00           int8   `json:"inifinityExtra00,uint8"`
	InfinityExtra01           int8   `json:"inifinityExtra01,uint8"`
	InfinityExtra02           int8   `json:"inifinityExtra02,uint8"`
	InfinityExtra03           int8   `json:"inifinityExtra03,uint8"`
	InfinityExtra04           int8   `json:"inifinityExtra04,uint8"`
	InfinityExtra05           int8   `json:"inifinityExtra05,uint8"`
	InfinityExtra06           int8   `json:"inifinityExtra06,uint8"`
	InfinityExtra07           int8   `json:"inifinityExtra07,uint8"`
	InfinityExtra08           int8   `json:"inifinityExtra08,uint8"`
	InfinityExtra09           int8   `json:"inifinityExtra09,uint8"`
	BonusCredits              int16  `json:"bonusCredits,uint16"`
	AppFriendsBonusCredits    int16  `json:"appFriendsBonusCredits,uint16"`
	OfferAvailable            uint8  `json:"offerAvailable"` // bool
	FirstGame                 uint8  `json:"firstGame"`      // bool

	// all players progress
	StagesProgressStat01 [8]uint32 `json:"stagesProgressStat01"` // count users reach that island
	StagesProgressStat02 [8]uint32 `json:"stagesProgressStat02"`

	// current player progress
	SubStagesRecordStats01 [8][]int8 `json:"subStagesRecordStats01"` // user progress in start mode
	SubStagesRecordStats02 [8][]int8 `json:"subStagesRecordStats02"` // casual mode (not more used)
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
	var needUpdate bool
	var triesRestore int64
	var userFriendsBonusCredits int16
	ctx := r.Context()
	value := ctx.Value(userCtxID)
	var user User
	now := time.Now()
	switch value.(type) {
	case nil:
		firstGame = 1
		platformID := platforms.GetByName(request.SysID)
		user = User{
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
			RemainingTries:          defaultConfig.DefaultRemainingTries,
			RestoreTriesAt:          0,
			Credits:                 defaultConfig.DefaultCredits.Vk, // TODO dehardcode `Vk`
			FriendsBonusCreditsTime: now.Unix(),
		}
		user.SetProgresStandart(defaultConfig.InitProgress)
		Gorm.Create(&user) // Gorm.NewRecord check row exists or somehow
	case User:
		user = value.(User)
		if user.FriendsBonusCreditsTime > now.Unix() {
			needUpdate = true
			to := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			user.FriendsBonusCreditsTime = to.Unix()
			userFriendsBonusCredits = int16(request.AppFriends) * defaultConfig.FriendsBonusCreditsMultiplier
		}
		if user.RestoreTriesAt != 0 && now.Unix() >= user.RestoreTriesAt {
			needUpdate = true
			if user.RemainingTries < defaultConfig.DefaultRemainingTries {
				user.RemainingTries = defaultConfig.DefaultRemainingTries
			}
			user.RestoreTriesAt = 0
		} else if user.RestoreTriesAt != 0 {
			triesRestore = user.RestoreTriesAt - now.Unix()
		}
	}
	if needUpdate {
		Gorm.Save(&user)
	}
	var usersProgress [8]uint32
	rows, err := Gorm.Table(
		"users",
	).Select(
		"count(*) as cnt, reached_stage01",
	).Where(
		"reached_stage01 > 0",
	).Group(
		"reached_stage01",
	).Order(
		"reached_stage01 desc",
	).Rows()
	if err != nil {
		log.Print(err)
	} else {
		var rowObject usersProgressRow
		i := -1
		for rows.Next() {
			err = rows.Scan(rowObject)
			if err != nil {
				log.Print(err)
				break
			}
			i++
			if rowObject.ReachedStage01 > uint8(i) {
				usersProgress[i] += rowObject.Cnt
			} else {
				break
			}
		}
	}
	response := enterResponse{
		ReqMsgID:                  request.MsgID,
		UserID:                    user.ID,
		ReachedStage01:            user.ReachedStage01,
		ReachedStage02:            user.ReachedStage02,
		ReachedSubStage01:         user.ReachedSubStage01,
		ReachedSubStage02:         user.ReachedSubStage02,
		IgnoreSavePointBlock:      user.IgnoreSavePointBlock,
		RemainingTries:            user.RemainingTries,
		TriesMin:                  defaultConfig.DefaultRemainingTries,
		TriesRegenSecondsInterval: defaultConfig.IntervalTriesRestoration,
		SecondsUntilTriesRegen:    triesRestore,
		Credits:                   user.Credits,
		InfinityExtra00:           user.InifinityExtra00,
		InfinityExtra01:           user.InifinityExtra01,
		InfinityExtra02:           user.InifinityExtra02,
		InfinityExtra03:           user.InifinityExtra03,
		InfinityExtra04:           user.InifinityExtra04,
		InfinityExtra05:           user.InifinityExtra05,
		InfinityExtra06:           user.InifinityExtra06,
		InfinityExtra07:           user.InifinityExtra07,
		InfinityExtra08:           user.InifinityExtra08,
		InfinityExtra09:           user.InifinityExtra09,
		OfferAvailable:            0,
		FirstGame:                 firstGame,
		BonusCredits:              0, // not used (every 12 hours user get reward. deleted now)
		AppFriendsBonusCredits:    userFriendsBonusCredits,
		StagesProgressStat01:      usersProgress,
		// StagesProgressStat02   TODO
		SubStagesRecordStats01: user.GetProgresStandart(),
		// SubStagesRecordStats02 TODO
	}
	JSON(w, response)
}