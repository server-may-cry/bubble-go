package models

// User struct in storage
type User struct {
	tableName               struct{} `sql:"users"`
	ID                      uint64   `gorm:"primary_key"`
	SysID                   uint8
	ExtID                   string
	ReachedStage01          int8
	ReachedSubStage01       int8
	ReachedStage02          int8 // not used. need in market. TODO remove
	ReachedSubStage02       int8 // not used. need in market. TODO remove
	IgnoreSavePointBlock    int8
	InifinityExtra00        int8
	InifinityExtra01        int8
	InifinityExtra02        int8
	InifinityExtra03        int8
	InifinityExtra04        int8
	InifinityExtra05        int8
	InifinityExtra06        int8
	InifinityExtra07        int8
	InifinityExtra08        int8
	InifinityExtra09        int8
	RemainingTries          int8
	RestoreTriesAt          int64
	Credits                 int16
	FriendsBonusCreditsTime int64
	// ProgressStandart        [][]int8 // json // gorm conflict
}
