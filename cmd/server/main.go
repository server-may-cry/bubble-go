package main

import (
	"log"
	"net/http"
	"net/url"
	"os"

	gorelic "github.com/brandfolder/gin-gorelic"
	"github.com/gin-gonic/gin"
	"github.com/server-may-cry/bubble-go/models"
	"gopkg.in/redis.v4"
)

var i = 0

func main() {
	t := models.Test{
		Ttt: 4,
	}
	log.Println(t)
	r := gin.Default()
	newRelicLicenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY")
	if newRelicLicenseKey == "" {
		log.Fatal("$NEW_RELIC_LICENSE_KEY must be set")
	}
	gorelic.InitNewrelicAgent(newRelicLicenseKey, "bubble-go", true)
	r.Use(gorelic.Handler)
	r.Use(gin.Logger())
	r.GET("/", func(c *gin.Context) {
		i = i + 1
		c.JSON(http.StatusOK, gin.H{
			"status":  "posted",
			"message": "msg",
			"i":       i,
		})
	})
	r.GET("/redis", func(c *gin.Context) {
		redisURI := os.Getenv("REDIS_URL")
		redisURL, err := url.Parse(redisURI)
		if err != nil {
			log.Fatal(err)
		}
		password, _ := redisURL.User.Password()
		client := redis.NewClient(&redis.Options{
			Addr:     redisURL.Host,
			Password: password,
			DB:       0, // default database
		})

		pong, err := client.Ping().Result()
		c.JSON(http.StatusOK, gin.H{
			"ping": pong,
		})
	})
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	r.Run(":" + port)
}
