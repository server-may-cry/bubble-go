package application

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
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
func ConfigInit(configFilePath string) error {
	file, err := os.Open(filepath.ToSlash(configFilePath))
	if err != nil {
		return errors.Wrap(err, "can`t open user.json")
	}
	err = json.NewDecoder(file).Decode(&defaultConfig)
	if err != nil {
		return errors.Wrap(err, "can`t decode user.json")
	}
	return nil
}
