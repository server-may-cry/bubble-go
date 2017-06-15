package market

import (
	"encoding/json"
	"testing"
)

var exampleMarketJSON = `
{
	"increase_pack": {
		"reward": {
			"increase": {
				"credits": 50
			}
		}
	},
	"set_pack": {
		"reward": {
			"set": {
				"credits": 800
			}
		}
	}
}
`

type testUser struct {
	Credits int
}

func getMarket() *Market {
	var config Config
	err := json.Unmarshal([]byte(exampleMarketJSON), &config)
	if err != nil {
		panic(err)
	}
	return NewMarket(config, "cdn://cdn.cdn/")
}

func TestMarketIncrease(t *testing.T) {
	market := getMarket()
	user := testUser{
		Credits: 100,
	}

	market.Buy(&user, "increase_pack")
	if user.Credits != 150 {
		t.Errorf("Buy(user, \"increase_pack\"): expected %d, actual %d", 150, user.Credits)
	}
}

func TestMarketSet(t *testing.T) {
	market := getMarket()
	user := testUser{
		Credits: 100,
	}

	market.Buy(&user, "set_pack")
	if user.Credits != 800 {
		t.Errorf("Buy(user, \"set_pack\"): expected %d, actual %d", 800, user.Credits)
	}
}

func TestMarketBuyNotExistPack(t *testing.T) {
	market := getMarket()
	user := testUser{
		Credits: 100,
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("panic expected on pack %s", "pack_not_exist")
		}
	}()
	market.Buy(&user, "pack_not_exist")
}

func TestMarketGetPack(t *testing.T) {
	market := getMarket()
	market.GetPack("increase_pack")
}

func TestMarketGetNotExistPack(t *testing.T) {
	market := getMarket()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("panic expected on GetPack %s", "pack_not_exist")
		}
	}()
	market.GetPack("pack_not_exist")
}

func TestMarketValidate(t *testing.T) {
	market := getMarket()
	user := testUser{}
	market.Validate(&user)
}
