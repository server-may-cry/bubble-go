package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/server-may-cry/bubble-go/models"
	"github.com/server-may-cry/bubble-go/protocol"
	"github.com/server-may-cry/bubble-go/storage"
)

var i = 0

// Index control index route
func Index(c *gin.Context) {
	i = i + 1

	c.JSON(http.StatusOK, protocol.IndexResponse{
		Status:  "posted",
		Message: "msg",
		I:       i,
	})
}

// Test used for debug
func Test(c *gin.Context) {
	t := models.Test{
		Ttt: 4,
	}
	log.Println(t)
	c.JSON(http.StatusOK, protocol.TestResponse{
		Test: t,
	})
}

// Redis used for send "ping" and receive "pong" from redis
func Redis(c *gin.Context) {
	pong, err := storage.Redis.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, protocol.RedisResponse{
		Ping: pong,
	})
}
