package application

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var defaultConfig struct {
	DefaultRemainingTries         int8  `json:"default_remaining_tries"`
	IntervalTriesRestoration      int   `json:"interval_tries_restoration"`
	FriendsBonusCreditsMultiplier int16 `json:"friends_bonus_credits_multiplier"`
	DefaultCredits                struct {
		Vk int16 `json:"vk"`
		Ok int16 `json:"ok"`
	} `json:"default_credits"`
}

func init() {
	configFile := "./config/user.json"
	file, err := os.Open(filepath.ToSlash(configFile))
	if err != nil {
		log.Fatal(err)
	}
	err = json.NewDecoder(file).Decode(&defaultConfig)
	if err != nil {
		log.Fatal(err)
	}
}
