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

// Market main market struct
type Market struct {
	packs     Config
	cdnPrefix string
}

// Buy get user and item name (from market config). Change user
func (m *Market) Buy(user interface{}, packName string) {
	pack, exist := m.packs[packName]
	if !exist {
		panic(fmt.Sprintf("try buy not existed pack %s", packName))
	}
	indirect := reflect.Indirect(reflect.ValueOf(user))
	for parameter, amount := range pack.Reward.Increase {
		Parameter := strings.Title(parameter)
		was := indirect.FieldByName(Parameter).Int()
		indirect.FieldByName(Parameter).SetInt(was + amount)
	}
	for parameter, amount := range pack.Reward.Set {
		Parameter := strings.Title(parameter)
		indirect.FieldByName(Parameter).SetInt(amount)
	}
}

// GetPack return pack description
func (m *Market) GetPack(packName string) Pack {
	pack, exist := m.packs[packName]
	if !exist {
		panic(fmt.Sprintf("try buy not existed pack %s", packName))
	}
	pack.Photo = fmt.Sprint(m.cdnPrefix, "productIcons/", pack.Photo, ".png")

	return pack
}

// NewMarket create new market instance
func NewMarket(config Config, cdn string) *Market {
	return &Market{
		packs:     config,
		cdnPrefix: cdn,
	}
}

// Validate check current market configuration for that type of user
func (m *Market) Validate(user interface{}) {
	for packName := range m.packs {
		m.Buy(user, packName)
	}
}
