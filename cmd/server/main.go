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
	"github.com/pressly/chi"
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
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
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
	market.Validate(user)

	application.VkWorker = notification.NewVkWorker(notification.VkConfig{
		AppID:           os.Getenv("VK_APP_ID"),
		Secret:          os.Getenv("VK_SECRET"),
		RequestInterval: time.Millisecond * 300,
	})
}

func main() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		application.JSON(w, h{
			"foo": "bar",
		})
	})

	router.Mount("/", func() http.Handler {
		r := chi.NewRouter()
		r.Use(application.AuthorizationMiddleware)
		r.Post("/ReqEnter", application.ReqEnter)
		r.Post("/ReqBuyProduct", application.ReqBuyProduct)
		r.Post("/ReqReduceTries", application.ReqReduceTries)
		r.Post("/ReqReduceCredits", application.ReqReduceCredits)
		r.Post("/ReqSavePlayerProgress", application.ReqSavePlayerProgress)
		r.Post("/ReqUsersProgress", application.ReqUsersProgress)
		return r
	}())
	router.Post("/VkPay", application.VkPay)

	router.Get("/crossdomain.xml", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<?xml version="1.0"?><cross-domain-policy><allow-access-from domain="*" /></cross-domain-policy>`))
	})
	// http://119226.selcdn.ru/bubble/ShootTheBubbleDevVK.html
	// http://bubble-srv-dev.herokuapp.com/bubble/ShootTheBubbleDevVK.html
	router.Get("/bubble/*filePath", application.ServeStatick)
	router.Get("/cache-clear", application.ClearStatickCache)

	router.Get("/exception", func(w http.ResponseWriter, r *http.Request) {
		panic("test log.Fatal")
	})

	router.Get("/debug-vk", func(w http.ResponseWriter, r *http.Request) {
		application.JSON(w, h{
			"levels": application.VkWorker.Levels,
			"events": application.VkWorker.Events,
		})
	})

	loaderio := os.Getenv("LOADERIO")
	loaderioRoute := fmt.Sprintf("/loaderio-%s", loaderio)
	router.Get(loaderioRoute, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("loaderio-%s", loaderio)))
	})

	port := os.Getenv("PORT")
	http.ListenAndServe(fmt.Sprint(":", port), router)
}
