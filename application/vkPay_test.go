package application

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/server-may-cry/bubble-go/market"
)

func TestVkBadSignature(t *testing.T) {
	server := httptest.NewServer(GetRouter(true))
	defer server.Close()

	jsonBytes, _ := json.Marshal(map[string]string{
		"app_id": "1",
		"sig":    "invalid_sig",
	})
	reader := bytes.NewReader(jsonBytes)

	os.Setenv("VK_SECRET", "secret")
	resp, err := http.Post(fmt.Sprint(server.URL, "/VkPay"), "application/json", reader)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	decoder := json.NewDecoder(resp.Body)
	var response map[string]interface{}
	err = decoder.Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	_, exist := response["error"]
	if !exist {
		t.Fatalf("expected error, got %+v", response)
	}
}

func TestVkGetItem(t *testing.T) {
	server := httptest.NewServer(GetRouter(true))
	defer server.Close()

	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open("sqlite3", file.Name())
	db.AutoMigrate(&User{})
	Gorm = db

	data := []byte("app_id=1item=creditsPack01lang=ru_RUnotification_type=get_itemorder_id=1receiver_id=123user_id=123secret")
	sig := fmt.Sprintf("%x", md5.Sum(data))
	jsonBytes, _ := json.Marshal(map[string]string{
		"app_id":            "1",
		"item":              "creditsPack01",
		"lang":              "ru_RU",
		"notification_type": "get_item",
		"order_id":          "1",
		"receiver_id":       "123",
		"user_id":           "123",
		"sig":               sig,
	})
	reader := bytes.NewReader(jsonBytes)

	market.Initialize(market.Config{
		"creditsPack01": market.Pack{},
	}, "")
	os.Setenv("VK_SECRET", "secret")
	resp, err := http.Post(fmt.Sprint(server.URL, "/VkPay"), "application/json", reader)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	decoder := json.NewDecoder(resp.Body)
	var response map[string]interface{}
	err = decoder.Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	_, exist := response["response"]
	if !exist {
		t.Fatalf("expected success, got %+v", response)
	}
}

func TestVkBuyItem(t *testing.T) {
	server := httptest.NewServer(GetRouter(true))
	defer server.Close()

	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open("sqlite3", file.Name())
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Transaction{})
	user := User{
		SysID:   1,
		ExtID:   123,
		Credits: 900,
	}
	db.Create(&user)
	Gorm = db

	data := []byte("app_id=1item=creditsPack01notification_type=order_status_changeorder_id=1receiver_id=123status=chargeableuser_id=123secret")
	sig := fmt.Sprintf("%x", md5.Sum(data))
	jsonBytes, _ := json.Marshal(map[string]string{
		"app_id":            "1",
		"item":              "creditsPack01",
		"notification_type": "order_status_change",
		"order_id":          "1",
		"receiver_id":       "123",
		"status":            "chargeable",
		"user_id":           "123",
		"sig":               sig,
	})
	reader := bytes.NewReader(jsonBytes)

	var exampleMarketJSON = `
	{
		"creditsPack01": {
			"reward": {
				"increase": {
					"credits": 140
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
	os.Setenv("VK_SECRET", "secret")
	resp, err := http.Post(fmt.Sprint(server.URL, "/VkPay"), "application/json", reader)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	decoder := json.NewDecoder(resp.Body)
	var response map[string]interface{}
	err = decoder.Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	_, exist := response["response"]
	if !exist {
		t.Fatalf("expected success, got %+v", response)
	}

	db.First(&user, user.ID)
	if user.Credits != 1040 {
		t.Fatalf("expected 1040 credits in db, got %+v", user)
	}
}
