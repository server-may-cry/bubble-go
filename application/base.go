package application

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/server-may-cry/bubble-go/market"
	"github.com/server-may-cry/bubble-go/notification"
)

// Gorm orm
var Gorm *gorm.DB

// VkWorker channel for send vk events
var VkWorker *notification.VkWorker

// Market struct
var Market *market.Market

type ctxID uint

const (
	userCtxID ctxID = iota
)

type h map[string]interface{}

// AuthRequestPart can be used to validate request
type AuthRequestPart struct {
	AuthKey    string `json:"authKey"`      // some hash
	ExtID      int64  `json:"extId,string"` // "123312693841263"
	SysID      string `json:"sysId"`        // "VK"
	SessionKey string `json:"sessionKey"`   // OK only
}

type baseRequest struct {
	// AuthRequestPart not more required
	MsgID uint64 `json:"msgId,string"` // not required. just for back capability
}

// JSON is helper to serve json http response
func JSON(w http.ResponseWriter, obj interface{}) {
	js, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
