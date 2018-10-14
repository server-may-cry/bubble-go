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

	"github.com/jinzhu/gorm"
	newrelic "github.com/newrelic/go-agent"
	"github.com/server-may-cry/bubble-go/market"
	"github.com/server-may-cry/bubble-go/models"
	"github.com/server-may-cry/bubble-go/mynewrelic"
	"github.com/server-may-cry/bubble-go/platforms"
	dig "go.uber.org/dig"
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
	OrderID    string `json:"order_id"`
	AppOrderID uint64 `json:"app_order_id"`
}

type itemResponse struct {
	ItemID   int    `json:"item_id"`
	Title    string `json:"title"`
	PhotoURL string `json:"photo_url"`
	Price    int    `json:"price"`
}

// VkPayForContainer for uber-go/dig
type VkPayForContainer struct {
	dig.Out

	HTTPHandler http.HandlerFunc
}

// VkPay acept and validate payment request from vk
func VkPay(db *gorm.DB, marketInstance *market.Market) VkPayForContainer {
	handler := VkPayForContainer{}
	handler.HTTPHandler = func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
			hashStr.WriteString(fmt.Sprint(k, "=", r.PostForm.Get(k)))
		}

		secret := os.Getenv("VK_SECRET")
		if secret == "" {
			panic("VK_SECRET not set")
		}
		data := []byte(fmt.Sprint(hashStr.String(), secret))
		expectedAuthKey := fmt.Sprintf("%x", md5.Sum(data))
		if expectedAuthKey != r.PostForm.Get("sig") {
			JSON(w, jsonHelper{
				"error": errorResponse{
					ErrorCode: 10,
					ErrorMsg:  "Invalid signature.",
					Critical:  true,
				},
			})
			return
		}
		switch r.PostForm.Get("notification_type") {
		case "get_item":
			fallthrough
		case "get_item_test":
			pack, err := marketInstance.GetPack(r.PostForm.Get("item"))
			if err != nil {
				panic(err)
			}
			JSON(w, jsonHelper{
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
			if r.PostForm.Get("status") != "chargeable" {
				JSON(w, jsonHelper{
					"error": errorResponse{
						ErrorCode: 100,
						ErrorMsg:  "Unexpected status. Expected - chargeable.",
						Critical:  true,
					},
				})
				return
			}
			var user models.User
			platformID := platforms.VK

			s := newrelic.DatastoreSegment{
				StartTime:  newrelic.StartSegmentNow(r.Context().Value(mynewrelic.Ctx).(newrelic.Transaction)),
				Product:    newrelic.DatastorePostgres,
				Collection: "user",
				Operation:  "SELECT",
			}
			db.Where("sys_id = ? AND ext_id = ?", platformID, r.PostForm.Get("user_id")).First(&user)
			_ = s.End()

			if user.ID == 0 { // check user exists
				panic("user not foud. try to buy")
			}
			err := marketInstance.Buy(&user, r.PostForm.Get("item"))
			if err != nil {
				panic(err)
			}
			orderID, err := strconv.ParseInt(r.PostForm.Get("order_id"), 10, 0)
			if err != nil {
				panic(fmt.Sprintf("cannot convert order id to int64 (%s)", r.PostForm.Get("order_id")))
			}

			ts := time.Now().Unix()
			transaction := models.Transaction{
				OrderID:     orderID,
				CreatedAt:   ts,
				UserID:      user,
				ConfirmedAt: ts,
			}

			s = newrelic.DatastoreSegment{
				StartTime:  newrelic.StartSegmentNow(r.Context().Value(mynewrelic.Ctx).(newrelic.Transaction)),
				Product:    newrelic.DatastorePostgres,
				Collection: "transaction",
				Operation:  "INSERT",
			}
			err = db.Create(&transaction).Error
			_ = s.End()
			if err != nil {
				panic(err)
			}

			s = newrelic.DatastoreSegment{
				StartTime:  newrelic.StartSegmentNow(r.Context().Value(mynewrelic.Ctx).(newrelic.Transaction)),
				Product:    newrelic.DatastorePostgres,
				Collection: "user",
				Operation:  "UPDATE",
			}
			db.Save(&user)
			_ = s.End()

			JSON(w, jsonHelper{
				"response": orderResponse{
					OrderID:    r.PostForm.Get("order_id"),
					AppOrderID: transaction.ID,
				},
			})
			return

		default:
			JSON(w, jsonHelper{
				"error": errorResponse{
					ErrorCode: 100,
					ErrorMsg:  "Unexpected type in 'notification_type'.",
					Critical:  true,
				},
			})
			return
		}
	}
	return handler
}
