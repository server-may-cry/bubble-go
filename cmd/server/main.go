package main

import (
	"log"
	"net/http"
	"net/url"
	"os"

	redis "gopkg.in/redis.v4"

	"github.com/server-may-cry/bubble-go/controllers"
	"github.com/server-may-cry/bubble-go/storage"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	rawRedisURL := os.Getenv("REDIS_URL")
	if rawRedisURL == "" {
		log.Fatal("$REDIS_URL must be set")
	}
	redisURL, err := url.Parse(rawRedisURL)
	if err != nil {
		log.Fatal(err)
	}
	password, _ := redisURL.User.Password()
	storage.Redis = redis.NewClient(&redis.Options{
		Addr:     redisURL.Host,
		Password: password,
		DB:       0, // default database
	})
}

func main() {
	http.HandleFunc("/echo", controllers.Echo)
	http.HandleFunc("/", controllers.Home)
	http.HandleFunc("/test", controllers.Test)
	http.HandleFunc("/redis", controllers.Redis)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
