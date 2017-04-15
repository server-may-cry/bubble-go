package controllers

import (
	"encoding/json"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

// AuthRequestPart can be used to validate request
type AuthRequestPart struct {
	AuthKey    string `json:"authKey" binding:"required"` // some hash
	ExtID      string `json:"extId" binding:"required"`   // "123312693841263"
	SysID      string `json:"sysId" binding:"required"`   // "VK"
	SessionKey string `json:"sessionKey"`                 // OK only
}

type baseRequest struct {
	// AuthRequestPart not more required
	MsgID uint64 `json:"msgId,string"` // not required. just for back capability
}

func getErrBody(err error) gin.H {
	return gin.H{"bad request": err}
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
