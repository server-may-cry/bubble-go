package application

import (
	"encoding/json"
	"net/http"
)

type ctxID uint

const (
	userCtxID ctxID = iota
)

type jsonHelper map[string]interface{}

// AuthRequestPart can be used to validate request
type AuthRequestPart struct {
	AuthKey    string `json:"authKey"`      // some hash
	ExtID      int64  `json:"extId,string"` // id on platform
	SysID      string `json:"sysId"`        // "VK"
	SessionKey string `json:"sessionKey"`   // OK only
}

type baseRequest struct {
	// AuthRequestPart not more required
	MsgID uint64 `json:"msgId,string"` // not required. just for back capability
}

// JSON is helper to serve json http response
func JSON(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
