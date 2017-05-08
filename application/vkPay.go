package application

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/server-may-cry/bubble-go/market"
	"github.com/server-may-cry/bubble-go/platforms"
)

/*
requests
{
    "app_id":"4890xxx",
    "item":"creditsPacxxx",
    "lang":"ru_RU",
    "notification_type":"get_item_test",
    "order_id":"829xxx",
    "receiver_id":"5523xxx",
    "user_id":"5523xxx",
    "sig":"bd59934272e8xxxx"
}
{
    "app_id":"4890948",
    "date":"1433503962",
    "item":"creditsPack01",
    "item_id":"1",
    "item_photo_url":"http:\\/\\/example.com\\/img.jpg",
    "item_price":"15",
    "item_title":"Extra help pack",
    "notification_type":"order_status_change_test",
    "order_id":"830232",
    "receiver_id":"5523718",
    "status":"chargeable",
    "user_id":"5523718",
    "sig":"bd59934272e8xxxx"
}
*/

// VkPay acept and validate payment request from vk
func VkPay(w http.ResponseWriter, r *http.Request) {
	var rawRequest map[string]string
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rawRequest)
	if err != nil {
		http.Error(w, getErrBody(err), http.StatusBadRequest)
		return
	}
	keys := make([]string, 0, len(rawRequest))
	for k := range rawRequest {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var hashStr bytes.Buffer
	for _, k := range keys {
		if k == "sig" {
			continue
		}
		hashStr.WriteString(fmt.Sprint(k, "=", rawRequest[k]))
	}

	secret := os.Getenv("VK_SECRET")
	if secret == "" {
		panic("VK_SECRET not set")
	}
	data := []byte(fmt.Sprint(hashStr.String(), secret))
	expectedAuthKey := fmt.Sprintf("%x", md5.Sum(data))
	if expectedAuthKey != rawRequest["sig"] {
		JSON(w, h{
			"error": h{
				"error_code": 10,
				"error_msg":  "Несовпадение вычисленной и переданной подписи запроса.",
				"critical":   true,
			},
		})
		return
	}
	switch rawRequest["notification_type"] {
	case "get_item":
		fallthrough
	case "get_item_test":
		pack := market.GetPack(rawRequest["item"])
		JSON(w, h{
			"response": h{
				"item_id":   1,
				"title":     pack.Title.Ru,
				"photo_url": pack.Photo,
				"price":     pack.Price,
			},
		})
		return

	case "order_status_change":
		fallthrough
	case "order_status_change_test":
		if rawRequest["status"] != "chargeable" {
			JSON(w, h{
				"error": h{
					"error_code": 100,
					"error_msg":  "Передано непонятно что вместо chargeable.",
					"critical":   true,
				},
			})
			return
		}
		db := Gorm
		var user User
		db.Where("sys_id = ? AND ext_id = ?", platforms.GetByName("VK"), rawRequest["user_id"]).First(&user)
		if user.ID != 0 { // check user exists
			panic("user not foud. try to buy")
		}
		market.Buy(&user, rawRequest["item"])
		orderID, err := strconv.ParseInt(rawRequest["order_id"], 10, 0)
		if err != nil {
			panic(fmt.Sprintf("cannot convert order id to int64 (%s)", rawRequest["order_id"]))
		}

		ts := time.Now().Unix()
		transaction := Transaction{
			OrderID:     orderID,
			CreatedAt:   ts,
			UserID:      user.ID,
			ConfirmedAt: ts,
		}
		success := Gorm.NewRecord(&transaction)
		if !success {
			panic(fmt.Sprintf("can`t create transaction %v", transaction))
		}
		JSON(w, h{
			"response": h{
				"order_id":     rawRequest["order_id"],
				"app_order_id": transaction.ID,
			},
		})
		return

	default:
		JSON(w, h{
			"error": h{
				"error_code": 100,
				"error_msg":  "Передано непонятно что в notification_type.",
				"critical":   true,
			},
		})
		return
	}
}
