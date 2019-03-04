package market

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/server-may-cry/bubble-go/models"
)

// Config struct for market offers and packs
type Config map[string]*Pack

// Pack is market element
type Pack struct {
	Price  Platforms    `json:"price"`
	Title  Locale       `json:"title"`
	Reward RewardStruct `json:"reward"`
	Photo  string       `json:"photo"`
}

// Platforms to specify price on exact platform
type Platforms struct {
	Vk int `json:"vk"`
	Ok int `json:"ok"`
}

// Locale used for configuring internationalized textes
type Locale struct {
	Ru string `json:"ru"`
}

// RewardStruct describe which user attribute need change and how much
type RewardStruct struct {
	Set      map[string]int64 `json:"set"`
	Increase map[string]int64 `json:"increase"`
}

// Market main market struct
type Market struct {
	packs Config
}

// NewMarket create new market instance
func NewMarket(config Config, cdn string) *Market {
	for i, pack := range config {
		config[i].Photo = fmt.Sprint(cdn, "productIcons/", pack.Photo, ".png")
	}
	return &Market{
		packs: config,
	}
}

// Buy get user and item name (from market config). Change user
func (m *Market) Buy(user *models.User, packName string) error {
	pack, exist := m.packs[packName]
	if !exist {
		return errors.New("try buy not existed pack " + packName)
	}
	indirect := reflect.Indirect(reflect.ValueOf(user))
	for parameter, amount := range pack.Reward.Increase {
		Parameter := strings.Title(parameter)
		field := indirect.FieldByName(Parameter)
		was := field.Int()
		field.SetInt(was + amount)
	}
	for parameter, amount := range pack.Reward.Set {
		Parameter := strings.Title(parameter)
		indirect.FieldByName(Parameter).SetInt(amount)
	}
	return nil
}

// GetPack return pack description
func (m *Market) GetPack(packName string) (*Pack, error) {
	pack, exist := m.packs[packName]
	if !exist {
		return nil, fmt.Errorf("try buy not existed pack %s", packName)
	}

	return pack, nil
}

// Validate check current market configuration for that type of user
func (m *Market) Validate(user *models.User) error {
	for packName, pack := range m.packs {
		if err := m.Buy(user, packName); err != nil {
			return err
		}
		repeatMap := make(map[string]struct{})
		for parameter := range pack.Reward.Increase {
			repeatMap[parameter] = struct{}{}
		}
		for parameter := range pack.Reward.Set {
			if _, exist := repeatMap[parameter]; exist {
				return fmt.Errorf("pack %s got set and increase for %s", packName, parameter)
			}
		}
	}
	return nil
}
