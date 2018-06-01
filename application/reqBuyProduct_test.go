package application

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/server-may-cry/bubble-go/market"
)

type buyProductCompleteRequest struct {
	AuthRequestPart
	buyProductRequest
}

func TestBuyProduct(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open("sqlite3", file.Name())
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&User{})
	user := User{
		SysID:   1,
		ExtID:   123,
		Credits: 900,
	}
	db.Create(&user)

	marketInstance := market.NewMarket(market.Config{
		"example_pack": &market.Pack{
			Reward: market.RewardStruct{
				Increase: map[string]int64{
					"credits": 15,
				},
			},
		},
	}, "")
	server := httptest.NewServer(GetRouter(true, db, marketInstance, nil))
	defer server.Close()

	data := []byte("_123_")
	jsonBytes, _ := json.Marshal(buyProductCompleteRequest{
		AuthRequestPart: AuthRequestPart{
			ExtID:   123,
			SysID:   "VK",
			AuthKey: fmt.Sprintf("%x", md5.Sum(data)),
		},
		buyProductRequest: buyProductRequest{
			ProductID: "example_pack",
		},
	})

	reader := bytes.NewReader(jsonBytes)
	resp, err := http.Post(fmt.Sprint(server.URL, "/ReqBuyProduct"), "application/json", reader)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	var response buyProductResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if response.Credits != 915 {
		t.Fatalf("expected 915 credits in response, got %+v", response)
	}

	db.First(&user, user.ID)
	if user.Credits != 915 {
		t.Fatalf("expected 915 credits in db, got %+v", user)
	}
}
