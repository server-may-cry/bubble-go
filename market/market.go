package market

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"

	"github.com/server-may-cry/bubble-go/models"
)

type marketConfigStruct map[string]struct {
	Price struct {
		Vk int `json:"vk"`
		Ok int `json:"ok"`
	} `json:"price"`
	Title struct {
		Ru string `json:"ru"`
	} `json:"title"`
	Reward struct {
		Set      map[string]int64 `json:"set"`
		Increase map[string]int64 `json:"increase"`
	} `json:"reward"`
	Photo string `json:"photo"`
}

var marketConfig marketConfigStruct

// Buy get user and item name (from market config). Change user
func Buy(user *models.User, packName string) {
	pack, exist := marketConfig[packName]
	if !exist {
		log.Fatalf("try buy not existed pack %s", packName)
	}
	for parameter, amount := range pack.Reward.Increase {
		r := reflect.ValueOf(user)
		was := reflect.Indirect(r).FieldByName(parameter).Int()
		reflect.Indirect(r).FieldByName(parameter).SetInt(was + amount)
	}
	for parameter, amount := range pack.Reward.Set {
		r := reflect.ValueOf(user)
		reflect.Indirect(r).FieldByName(parameter).SetInt(amount)
	}
	// TODO wat?
}

// InitializeMarket load market config from market.json
func InitializeMarket() {
	if len(marketConfig) > 1 {
		return
	}
	file, err := ioutil.ReadFile(filepath.ToSlash("./market/market.json"))
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(file, &marketConfig)
}
