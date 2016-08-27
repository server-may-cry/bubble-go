package main

import (
	"log"
	"net/url"
	"os"

	gorelic "github.com/brandfolder/gin-gorelic"
	"github.com/gin-gonic/gin"
	"github.com/server-may-cry/bubble-go/controllers"
	"github.com/server-may-cry/bubble-go/storage"
	"gopkg.in/redis.v4"
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
	storage.Redis = redis.NewClient(&redis.Options{
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
	r.GET("/redis", controllers.Redis)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	r.Run(":" + port)
}
