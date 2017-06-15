package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/server-may-cry/bubble-go/application"
	"github.com/server-may-cry/bubble-go/market"
	"github.com/server-may-cry/bubble-go/notification"
)

type h map[string]interface{}

func initialize() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	dbURL := os.Getenv("DB_URL")
	u, err := url.Parse(dbURL)
	if err != nil {
		log.Fatalf("can`t parse DB_URL (%s)", dbURL)
	}

	host, port, _ := net.SplitHostPort(u.Host)
	password, _ := u.User.Password()
	db, err := gorm.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		host,
		port,
		u.User.Username(),
		password,
		strings.Trim(u.Path, "/"),
	))
	if err != nil {
		log.Fatalf("failed to connect database: %s", err)
	}
	db.AutoMigrate(&application.User{})
	db.AutoMigrate(&application.Transaction{})
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(17) // 20 actual limit
	application.Gorm = db

	marketConfigFile := "./config/market.json"
	file, err := os.Open(filepath.ToSlash(marketConfigFile))
	if err != nil {
		log.Fatalf("can`t open market.json error: %s", err)
	}
	var marketConfig market.Config
	err = json.NewDecoder(file).Decode(&marketConfig)
	if err != nil {
		log.Fatalf("can`t decode market.json error: %s", err)
	}
	application.Market = market.NewMarket(marketConfig, os.Getenv("CDN_PREFIX"))
	user := application.User{}
	application.Market.Validate(&user)

	application.ConfigInit("./config/user.json")

	application.VkWorker = notification.NewVkWorker(notification.VkConfig{
		AppID:           os.Getenv("VK_APP_ID"),
		Secret:          os.Getenv("VK_SECRET"),
		RequestInterval: time.Millisecond * 300,
	})
}

func main() {
	initialize()
	router := application.GetRouter(false)
	port := os.Getenv("PORT")
	http.ListenAndServe(fmt.Sprint(":", port), router)
}
