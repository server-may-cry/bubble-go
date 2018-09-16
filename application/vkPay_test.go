package application

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	newrelic "github.com/newrelic/go-agent"
	"github.com/server-may-cry/bubble-go/market"
	"github.com/server-may-cry/bubble-go/models"
	"github.com/server-may-cry/bubble-go/mynewrelic"
)

func TestVkBadSignature(t *testing.T) {
	form := url.Values{}
	form.Add("app_id", "1")
	form.Add("sig", "invalid_sig")
	reader := strings.NewReader(form.Encode())

	os.Setenv("VK_SECRET", "secret")
	req, err := http.NewRequest("POST", "/VkPay", reader)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	vkPayContainer := VkPay(nil, nil)

	vkPayContainer.HTTPHandler.ServeHTTP(resp, req)
	if resp.Code != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.Code)
	}
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	_, exist := response["error"]
	if !exist {
		t.Fatalf("expected error, got %+v", response)
	}
}

func TestVkGetItem(t *testing.T) {
	marketInstance := market.NewMarket(market.Config{
		"creditsPack01": &market.Pack{},
	}, "")

	data := []byte(
		"app_id=1item=creditsPack01lang=ru_RUnotification_type=get_itemorder_id=1receiver_id=123user_id=123secret",
	)
	form := url.Values{}
	form.Add("app_id", "1")
	form.Add("item", "creditsPack01")
	form.Add("lang", "ru_RU")
	form.Add("notification_type", "get_item")
	form.Add("order_id", "1")
	form.Add("receiver_id", "123")
	form.Add("user_id", "123")
	form.Add("sig", fmt.Sprintf("%x", md5.Sum(data)))
	reader := strings.NewReader(form.Encode())

	os.Setenv("VK_SECRET", "secret")
	req, err := http.NewRequest("POST", "/VkPay", reader)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp := httptest.NewRecorder()
	vkPayContainer := VkPay(nil, marketInstance)

	vkPayContainer.HTTPHandler.ServeHTTP(resp, req)
	if resp.Code != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.Code)
	}
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	_, exist := response["response"]
	if !exist {
		t.Fatalf("expected success, got %+v", response)
	}
}

func TestVkBuyItem(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open("sqlite3", file.Name())
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Transaction{})
	user := models.User{
		SysID:   1,
		ExtID:   123,
		Credits: 900,
	}
	db.Create(&user)

	marketInstance := market.NewMarket(market.Config{
		"creditsPack01": &market.Pack{
			Reward: market.RewardStruct{
				Increase: map[string]int64{
					"credits": 140,
				},
			},
		},
	}, "")

	data := []byte(
		"app_id=1item=creditsPack01notification_type=order_status_changeorder_id=1" +
			"receiver_id=123status=chargeableuser_id=123secret",
	)
	form := url.Values{}
	form.Add("app_id", "1")
	form.Add("item", "creditsPack01")
	form.Add("notification_type", "order_status_change")
	form.Add("order_id", "1")
	form.Add("receiver_id", "123")
	form.Add("status", "chargeable")
	form.Add("user_id", "123")
	form.Add("sig", fmt.Sprintf("%x", md5.Sum(data)))
	reader := strings.NewReader(form.Encode())

	os.Setenv("VK_SECRET", "secret")
	req, err := http.NewRequest("POST", "/VkPay", reader)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	config := newrelic.NewConfig("bubble-go", "1234567890123456789012345678901234567890")
	app, err := newrelic.NewApplication(config)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), mynewrelic.Ctx, app.StartTransaction("test", nil, nil))
	req = req.WithContext(ctx)

	resp := httptest.NewRecorder()
	vkPayContainer := VkPay(db, marketInstance)

	vkPayContainer.HTTPHandler.ServeHTTP(resp, req)
	if resp.Code != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.Code)
	}
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
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
