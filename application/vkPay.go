package application

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/server-may-cry/bubble-go/platforms"
)

/*
requests post paameters actually
    "app_id":"4890xxx",
    "item":"creditsPacxxx",
    "lang":"ru_RU",
    "notification_type":"get_item_test",
    "order_id":"829xxx",
    "receiver_id":"5523xxx",
    "user_id":"5523xxx",
    "sig":"bd59934272e8xxxx"

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
*/

type errorResponse struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
	Critical  bool   `json:"critical"`
}

type orderResponse struct {
	OrderID    interface{} `json:"order_id"`
	AppOrderID uint64      `json:"app_order_id"`
}

type itemResponse struct {
	ItemID   int    `json:"item_id"`
	Title    string `json:"title"`
	PhotoURL string `json:"photo_url"`
	Price    int    `json:"price"`
}

// VkPay acept and validate payment request from vk
func VkPay(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if err := r.ParseForm(); err != nil {
		http.Error(w, getErrBody(err), http.StatusBadRequest)
	}
	keys := make([]string, 0, len(r.PostForm))
	for k := range r.PostForm {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var hashStr bytes.Buffer
	for _, k := range keys {
		if k == "sig" {
			continue
		}
		hashStr.WriteString(fmt.Sprint(k, "=", r.PostFormValue(k)))
	}

	secret := os.Getenv("VK_SECRET")
	if secret == "" {
		panic("VK_SECRET not set")
	}
	data := []byte(fmt.Sprint(hashStr.String(), secret))
	expectedAuthKey := fmt.Sprintf("%x", md5.Sum(data))
	if expectedAuthKey != r.PostFormValue("sig") {
		JSON(w, h{
			"error": errorResponse{
				ErrorCode: 10,
				ErrorMsg:  "Несовпадение вычисленной и переданной подписи запроса.",
				Critical:  true,
			},
		})
		return
	}
	switch r.PostFormValue("notification_type") {
	case "get_item":
		fallthrough
	case "get_item_test":
		pack := Market.GetPack(r.PostFormValue("item"))
		JSON(w, h{
			"response": itemResponse{
				ItemID:   1,
				Title:    pack.Title.Ru,
				PhotoURL: pack.Photo,
				Price:    pack.Price.Vk,
			},
		})
		return

	case "order_status_change":
		fallthrough
	case "order_status_change_test":
		if r.PostFormValue("status") != "chargeable" {
			JSON(w, h{
				"error": errorResponse{
					ErrorCode: 100,
					ErrorMsg:  "Передано непонятно что вместо chargeable.",
					Critical:  true,
				},
			})
			return
		}
		db := Gorm
		var user User
		db.Where("sys_id = ? AND ext_id = ?", platforms.GetByName("VK"), r.PostFormValue("user_id")).First(&user)
		if user.ID == 0 { // check user exists
			panic("user not foud. try to buy")
		}
		Market.Buy(&user, r.PostFormValue("item"))
		orderID, err := strconv.ParseInt(r.PostFormValue("order_id"), 10, 0)
		if err != nil {
			panic(fmt.Sprintf("cannot convert order id to int64 (%s)", r.PostFormValue("order_id")))
		}

		ts := time.Now().Unix()
		transaction := Transaction{
			OrderID:     orderID,
			CreatedAt:   ts,
			UserID:      user,
			ConfirmedAt: ts,
		}
		Gorm.Create(&transaction)
		Gorm.Save(&user)
		JSON(w, h{
			"response": orderResponse{
				OrderID:    r.PostFormValue("order_id"),
				AppOrderID: transaction.ID,
			},
		})
		return

	default:
		JSON(w, h{
			"error": errorResponse{
				ErrorCode: 100,
				ErrorMsg:  "Передано непонятно что в notification_type.",
				Critical:  true,
			},
		})
		return
	}
}
