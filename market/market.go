package market

import (
	"log"
	"reflect"

	"github.com/server-may-cry/bubble-go/models"
)

// Config struct for market offers and packs
type Config map[string]struct {
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

var marketConfig Config

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
}

// InitializeMarket load market config from market.json
func InitializeMarket(config Config) {
	marketConfig = config
}
