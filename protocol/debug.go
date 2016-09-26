package protocol

import "github.com/server-may-cry/bubble-go/models"

type IndexResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	I       int    `json:"i"`
}

type TestResponse struct {
	Test models.Test `json:"test"`
}

type RedisResponse struct {
	Ping string `json:"ping"`
}

type EnterRequest struct {
	AppFriends uint8  `json:"appFriends,string"`
	AuthKey    string `json:"authKey"`
	ExtId      string `json:"extId"`
	MsgId      uint   `json:"msgId,string"`
	Referer    string `json:"referer"`
	SrcExtId   string `json:"srcExtId"`
	SysId      string `json:"sysId"`
}

type EnterResponse struct {
	reqMsgId               uint
	userId                 uint32
	reachedStage01         uint8
	reachedStage02         uint8
	reachedSubStage01      uint8
	reachedSubStage02      uint8
	ignoreSavePointBlock   bool
	remainingTries         uint8
	credits                uint16
	inifinityExtra00       bool
	inifinityExtra01       bool
	inifinityExtra02       bool
	inifinityExtra03       bool
	inifinityExtra04       bool
	inifinityExtra05       bool
	inifinityExtra06       bool
	inifinityExtra07       bool
	inifinityExtra08       bool
	inifinityExtra09       bool
	bonusCredits           uint16
	appFriendsBonusCredits uint16
	offerAvailable         bool
	firstGame              bool
	stagesProgressStat01   [8]uint32
	stagesProgressStat02   [8]uint32
	subStagesRecordStats01 [8][]uint8
	subStagesRecordStats02 [8][]uint8
}
