package market_test

import (
	"encoding/json"
	"testing"

	"github.com/server-may-cry/bubble-go/market"
	"github.com/server-may-cry/bubble-go/models"
)

var exampleMarketJson = `
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

func init() {
	var config market.Config
	json.Unmarshal([]byte(exampleMarketJson), &config)
	market.Initialize(config)
}

func TestMarketIncrease(t *testing.T) {
	user := models.User{
		Credits: 100,
	}

	market.Buy(&user, "increase_pack")
	if user.Credits != 150 {
		t.Errorf("market.Buy(user, \"increase_pack\"): expected %d, actual %d", 150, user.Credits)
	}
}

func TestMarketSet(t *testing.T) {
	user := models.User{
		Credits: 100,
	}

	market.Buy(&user, "set_pack")
	if user.Credits != 800 {
		t.Errorf("market.Buy(user, \"set_pack\"): expected %d, actual %d", 800, user.Credits)
	}
}

func TestMarketPanic(t *testing.T) {
	user := models.User{
		Credits: 100,
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("panic expected on pack %s", "pack_not_exist")
		}
	}()
	market.Buy(&user, "pack_not_exist")
}
