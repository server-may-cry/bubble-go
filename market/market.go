package market

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/server-may-cry/bubble-go/models"
)

// Config struct for market offers and packs
type Config map[string]Pack

// Pack is market element
type Pack struct {
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
		panic(fmt.Sprintf("try buy not existed pack %s", packName))
	}
	for parameter, amount := range pack.Reward.Increase {
		Parameter := strings.Title(parameter)
		r := reflect.ValueOf(user)
		was := reflect.Indirect(r).FieldByName(Parameter).Int()
		reflect.Indirect(r).FieldByName(Parameter).SetInt(was + amount)
	}
	for parameter, amount := range pack.Reward.Set {
		Parameter := strings.Title(parameter)
		r := reflect.ValueOf(user)
		reflect.Indirect(r).FieldByName(Parameter).SetInt(amount)
	}
}

// GetPack return pack description
func GetPack(packName string) Pack {
	pack, exist := marketConfig[packName]
	if !exist {
		panic(fmt.Sprintf("try buy not existed pack %s", packName))
	}
	// TODO add CDN prefix

	return pack
}

// Initialize load market config from market.json
func Initialize(config Config) {
	marketConfig = config

	// check all packs valid
	user := models.User{}
	for packName := range marketConfig {
		Buy(&user, packName)
	}
}
