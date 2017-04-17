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

func init() {
	var config Config
	json.Unmarshal([]byte(exampleMarketJSON), &config)
	Initialize(config, "cdn://cdn.cdn/")
}

func TestMarketIncrease(t *testing.T) {
	user := testUser{
		Credits: 100,
	}

	Buy(&user, "increase_pack")
	if user.Credits != 150 {
		t.Errorf("Buy(user, \"increase_pack\"): expected %d, actual %d", 150, user.Credits)
	}
}

func TestMarketSet(t *testing.T) {
	user := testUser{
		Credits: 100,
	}

	Buy(&user, "set_pack")
	if user.Credits != 800 {
		t.Errorf("Buy(user, \"set_pack\"): expected %d, actual %d", 800, user.Credits)
	}
}

func TestMarketPanic(t *testing.T) {
	user := testUser{
		Credits: 100,
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("panic expected on pack %s", "pack_not_exist")
		}
	}()
	Buy(&user, "pack_not_exist")
}

func TestMarketGetPack(t *testing.T) {
	GetPack("increase_pack")
}

func TestMarketGetNotExistPack(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("panic expected on GetPack %s", "pack_not_exist")
		}
	}()
	GetPack("pack_not_exist")
}
