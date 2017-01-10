package models

// User struct in storage
type User struct {
	tableName               struct{} `sql:"users"`
	ID                      uint64
	SysID                   uint8
	ExtID                   string
	ReachedStage01          uint8
	ReachedSubStage01       uint8
	IgnoreSavePointBlock    uint8
	InifinityExtra00        uint8
	InifinityExtra01        uint8
	InifinityExtra02        uint8
	InifinityExtra03        uint8
	InifinityExtra04        uint8
	InifinityExtra05        uint8
	InifinityExtra06        uint8
	InifinityExtra07        uint8
	InifinityExtra08        uint8
	InifinityExtra09        uint8
	RemainingTries          uint8
	RestoreTriesAt          int64
	Credits                 uint16
	FriendsBonusCreditsTime int64
	ProgressStandart        [][]int8 // json
}
