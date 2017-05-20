package application

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var defaultConfig struct {
	DefaultRemainingTries         int8 `json:"default_remaining_tries"`
	IntervalTriesRestoration      int  `json:"interval_tries_restoration"`
	FriendsBonusCreditsMultiplier int  `json:"friends_bonus_credits_multiplier"`
	DefaultCredits                struct {
		Vk int `json:"vk"`
		Ok int `json:"ok"`
	} `json:"default_credits"`
	InitProgress [7][]int8 `json:"init_progress"`
}

// ConfigInit pass config file
func ConfigInit(configFilePath string) {
	file, err := os.Open(filepath.ToSlash(configFilePath))
	if err != nil {
		log.Fatalf("can`t open user.json error: %s", err)
	}
	err = json.NewDecoder(file).Decode(&defaultConfig)
	if err != nil {
		log.Fatalf("cant decode user.json error: %s", err)
	}
}
