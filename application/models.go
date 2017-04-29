package application

import (
	"encoding/json"
)

// User struct in storage
type User struct {
	ID                      uint `gorm:"primary_key"`
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
	ProgressStandart        string // [][]int8 json
}

// GetProgresStandart return user progress as array
func (u *User) GetProgresStandart() [][]uint8 {
	var progress [][]uint8
	err := json.Unmarshal([]byte(u.ProgressStandart), &progress)
	if err != nil {
		panic(err)
	}
	return progress
}

// SetProgresStandart set user progress
func (u *User) SetProgresStandart(progress [][]uint8) {
	r, err := json.Marshal(progress)
	if err != nil {
		panic(err)
	}
	u.ProgressStandart = string(r[:])
}

// Transaction log payment requests
type Transaction struct {
	ID          uint `gorm:"primary_key"`
	OrderID     int64
	CreatedAt   int64
	UserID      uint
	ConfirmedAt int64
}
