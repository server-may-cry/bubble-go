package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/server-may-cry/bubble-go/controllers"
	"github.com/server-may-cry/bubble-go/market"
	"github.com/server-may-cry/bubble-go/models"
	"github.com/server-may-cry/bubble-go/mymiddleware"
	"github.com/server-may-cry/bubble-go/notification"
	"github.com/server-may-cry/bubble-go/storage"
)

type h map[string]interface{}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	sqliteFile, err := ioutil.TempFile("", "bubble.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	db, err := gorm.Open("sqlite3", sqliteFile.Name())
	if err != nil {
		log.Fatalf("failed to connect database: %s", err)
	}
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Transaction{})
	storage.Gorm = db

	marketConfigFile := "./config/market.json"
	file, err := ioutil.ReadFile(filepath.ToSlash(marketConfigFile))
	if err != nil {
		log.Fatal(err)
	}

	var marketConfig market.Config
	json.Unmarshal(file, &marketConfig)
	market.Initialize(marketConfig)

	notification.VkEventChan = make(chan notification.VkEvent)
	go notification.VkWorkerInit(notification.VkEventChan)
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
		controllers.JSON(w, h{
			"foo": "bar",
		})
	})

	router.Mount("/", func() http.Handler {
		r := chi.NewRouter()
		r.Use(mymiddleware.AuthorizationMiddleware)
		r.Post("/ReqEnter", controllers.ReqEnter)
		r.Post("/ReqBuyProduct", controllers.ReqBuyProduct)
		r.Post("/ReqReduceTries", controllers.ReqReduceTries)
		r.Post("/ReqReduceCredits", controllers.ReqReduceCredits)
		r.Post("/ReqSavePlayerProgress", controllers.ReqSavePlayerProgress)
		r.Post("/ReqUsersProgress", controllers.ReqUsersProgress)
		return r
	}())
	router.Post("/VkPay", controllers.VkPay)

	router.Get("/crossdomain.xml", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<?xml version="1.0"?><cross-domain-policy><allow-access-from domain="*" /></cross-domain-policy>`))
	})
	// http://119226.selcdn.ru/bubble/ShootTheBubbleDevVK.html
	// http://bubble-srv-dev.herokuapp.com/bubble/ShootTheBubbleDevVK.html
	router.Get("/bubble/*filePath", controllers.ServeStatick)
	router.Get("/cache-clear", controllers.ClearStatickCache)

	router.Get("/exception", func(w http.ResponseWriter, r *http.Request) {
		panic("test log.Fatal")
	})

	loaderio := os.Getenv("LOADERIO")
	loaderioRoute := fmt.Sprintf("/loaderio-%s", loaderio)
	router.Get(loaderioRoute, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("loaderio-%s", loaderio)))
	})

	port := os.Getenv("PORT")
	http.ListenAndServe(fmt.Sprint(":", port), router)
}
