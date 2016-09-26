package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"github.com/server-may-cry/bubble-go/controllers"
	"github.com/server-may-cry/bubble-go/storage"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	rawMongoURL := os.Getenv("MONGODB_URI")
	if rawRedisURL == "" {
		log.Fatal("$MONGODB_URI must be set")
	}
	mongoURL, err := url.Parse(rawRedisURL)
	user, _ := mongoURL.User.User()
	password, _ := mongoURL.User.Password()
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{mongoURL.Host},
		Timeout:  60 * time.Second,
		Database: user,
		Username: user,
		Password: password,
	}
	if err != nil {
		log.Fatal(err)
	}
	storage.Redis = redis.NewClient(&redis.Options{
		Addr:     mongoURL.Host,
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
