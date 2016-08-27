package main

import (
	"log"
	"net/http"
	"net/url"
	"os"

	gorelic "github.com/brandfolder/gin-gorelic"
	"github.com/gin-gonic/gin"
	"github.com/server-may-cry/bubble-go/controllers"
	"gopkg.in/redis.v4"
)

var (
	redisClient *redis.Client
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	redisURI := os.Getenv("REDIS_URL")
	if redisURI == "" {
		log.Fatal("$REDIS_URL must be set")
	}
	redisURL, err := url.Parse(redisURI)
	if err != nil {
		log.Fatal(err)
	}
	password, _ := redisURL.User.Password()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisURL.Host,
		Password: password,
		DB:       0, // default database
	})

	newRelicLicenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY")
	if newRelicLicenseKey == "" {
		log.Fatal("$NEW_RELIC_LICENSE_KEY must be set")
	}
	gorelic.InitNewrelicAgent(newRelicLicenseKey, "bubble-go", true)
}

func main() {
	r := gin.Default()

	r.Use(gorelic.Handler)
	r.Use(gin.Logger())

	r.GET("/", controllers.Index)
	r.GET("/test", controllers.Test)
	r.GET("/redis", func(c *gin.Context) {
		pong, err := redisClient.Ping().Result()
		if err != nil {
			log.Fatal(err)
		}
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
