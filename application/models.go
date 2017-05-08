package application

import (
	"encoding/json"
)

// User struct in storage
type User struct {
	ID                      uint64 `gorm:"primary_key"`
	SysID                   uint8  `gorm:"unique_index:user_search_index"`
	ExtID                   int64  `gorm:"unique_index:user_search_index"`
	ReachedStage01          int8
	ReachedSubStage01       int8
	ReachedStage02          int8 // not used. need in market
	ReachedSubStage02       int8 // not used. need in market
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
	RestoreTriesAt          int64 `sql:"not null"`
	Credits                 int16
	FriendsBonusCreditsTime int64
	ProgressStandart        string // [][]int8 json
}

// GetProgresStandart return user progress as array
func (u *User) GetProgresStandart() (progress [7][]int8) {
	err := json.Unmarshal([]byte(u.ProgressStandart), &progress)
	if err != nil {
		panic(err)
	}
	return progress
}

// SetProgresStandart set user progress
func (u *User) SetProgresStandart(progress [7][]int8) {
	r, err := json.Marshal(progress)
	if err != nil {
		panic(err)
	}
	u.ProgressStandart = string(r[:])
}

// Transaction log payment requests
type Transaction struct {
	ID          uint64 `gorm:"primary_key"`
	OrderID     int64  `sql:"not null"`
	CreatedAt   int64  `sql:"not null"`
	UserID      uint64 `sql:"not null"`
	ConfirmedAt int64
}
