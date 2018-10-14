package market

import (
	"testing"

	"github.com/server-may-cry/bubble-go/models"
)

type testUser struct {
	Credits int
}

func getMarket() *Market {
	return NewMarket(Config{
		"increase_pack": &Pack{
			Reward: RewardStruct{
				Increase: map[string]int64{
					"credits": 50,
				},
			},
		},
		"set_pack": &Pack{
			Reward: RewardStruct{
				Set: map[string]int64{
					"credits": 800,
				},
			},
		},
	}, "cdn://cdn.cdn/")
}

func TestMarketIncrease(t *testing.T) {
	market := getMarket()
	user := models.User{
		Credits: 100,
	}

	err := market.Buy(&user, "increase_pack")
	if err != nil {
		t.Errorf("market.Buy error: %s", err)
	}
	if user.Credits != 150 {
		t.Errorf("Buy(user, \"increase_pack\"): expected %d, actual %d", 150, user.Credits)
	}
}

func TestMarketSet(t *testing.T) {
	market := getMarket()
	user := models.User{
		Credits: 100,
	}

	err := market.Buy(&user, "set_pack")
	if err != nil {
		t.Errorf("market.Buy error: %s", err)
	}
	if user.Credits != 800 {
		t.Errorf("Buy(user, \"set_pack\"): expected %d, actual %d", 800, user.Credits)
	}
}

func TestMarketBuyNotExistPack(t *testing.T) {
	market := getMarket()
	user := models.User{
		Credits: 100,
	}

	err := market.Buy(&user, "pack_not_exist")
	if err == nil {
		t.Errorf("error expected on pack %s", "pack_not_exist")
	}
}

func TestMarketGetPack(t *testing.T) {
	market := getMarket()
	pack, err := market.GetPack("increase_pack")
	if pack == nil {
		t.Error("pack found expocted on increase_pack")
	}
	if err != nil {
		t.Errorf("no error expected on increase_pack, got %s", err.Error())
	}
}

func TestMarketGetNotExistPack(t *testing.T) {
	market := getMarket()
	pack, err := market.GetPack("pack_not_exist")
	if pack != nil {
		t.Error("no pack expected on pack_not_exist")
	}
	if err == nil {
		t.Error("error expected on pack_not_exist")
	}
}

func TestMarketValidate(t *testing.T) {
	market := getMarket()
	user := models.User{}
	err := market.Validate(&user)
	if err != nil {
		t.Errorf("no error expected on validation, got %s", err.Error())
	}
}
