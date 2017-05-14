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
	server := httptest.NewServer(GetRouter(true))
	defer server.Close()

	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open("sqlite3", file.Name())
	db.AutoMigrate(&User{})
	user := User{
		SysID:   1,
		ExtID:   123,
		Credits: 900,
	}
	db.Create(&user)
	Gorm = db

	data := []byte("_123_")
	authKey := fmt.Sprintf("%x", md5.Sum(data))
	jsonBytes, _ := json.Marshal(buyProductCompleteRequest{
		AuthRequestPart: AuthRequestPart{
			ExtID:   123,
			SysID:   "VK",
			AuthKey: authKey,
		},
		buyProductRequest: buyProductRequest{
			ProductID: "example_pack",
		},
	})

	var exampleMarketJSON = `
	{
		"example_pack": {
			"reward": {
				"increase": {
					"credits": 15
				}
			}
		}
	}
	`
	var config market.Config
	err = json.Unmarshal([]byte(exampleMarketJSON), &config)
	if err != nil {
		t.Fatal(err)
	}
	market.Initialize(config, "")

	reader := bytes.NewReader(jsonBytes)
	resp, err := http.Post(fmt.Sprint(server.URL, "/ReqBuyProduct"), "application/json", reader)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	decoder := json.NewDecoder(resp.Body)
	var response buyProductResponse
	err = decoder.Decode(&response)
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
