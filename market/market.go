package market

import (
	"fmt"
	"reflect"
	"strings"
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
var cdnPrefix string

// Buy get user and item name (from market config). Change user
func Buy(user interface{}, packName string) {
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
	pack.Photo = fmt.Sprint(cdnPrefix, pack.Photo)

	return pack
}

// Initialize load market config from market.json
func Initialize(config Config, cdn string) {
	marketConfig = config
	cdnPrefix = cdn
}

// Validate check current market configuration for that type of user
func Validate(user interface{}) {
	for packName := range marketConfig {
		Buy(&user, packName)
	}
}
