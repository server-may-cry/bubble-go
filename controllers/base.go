package controllers

import (
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
