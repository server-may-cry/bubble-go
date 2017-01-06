package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
)

type baseRequest struct {
	AuthKey string `json:"authKey" binding:"required"` // some hash
	ExtID   string `json:"extId" binding:"required"`   // "123312693841263"
	SysID   string `json:"sysId" binding:"required"`   // "VK"
	MsgID   uint64 `json:"msgId"`                      // not required. just for back capability
}

func getErrBody(err error) gin.H {
	return gin.H{"bad request": err}
}
