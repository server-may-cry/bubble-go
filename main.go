package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/server-may-cry/bubble-go/models"
)

var i = 0

func main() {
	t := models.Test{
		Ttt: 4,
	}
	log.Println(t)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		i = i + 1
		c.JSON(http.StatusOK, gin.H{
			"status":  "posted",
			"message": "msg",
			"i":       i,
		})
	})
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
