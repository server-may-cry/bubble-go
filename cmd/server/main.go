package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pressly/chi/middleware"
	"github.com/server-may-cry/bubble-go/application"
	"github.com/server-may-cry/bubble-go/market"
	"github.com/server-may-cry/bubble-go/notification"
)

type h map[string]interface{}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	dbURL := os.Getenv("DB_URL")
	u, err := url.Parse(dbURL)
	if err != nil {
		log.Fatalf("can`t parse DB_URL (%s)", dbURL)
	}

	hostPort := strings.Split(u.Host, ":")
	password, _ := u.User.Password()
	db, err := gorm.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		hostPort[0],
		hostPort[1],
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
	market.Initialize(marketConfig, os.Getenv("CDN_PREFIX"))
	user := application.User{}
	market.Validate(&user)

	application.ConfigInit("./config/user.json")

	application.VkWorker = notification.NewVkWorker(notification.VkConfig{
		AppID:           os.Getenv("VK_APP_ID"),
		Secret:          os.Getenv("VK_SECRET"),
		RequestInterval: time.Millisecond * 300,
	})
}

func main() {
	router := application.GetRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	port := os.Getenv("PORT")
	http.ListenAndServe(fmt.Sprint(":", port), router)
}
