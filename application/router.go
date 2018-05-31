package application

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jinzhu/gorm"
	newrelic "github.com/newrelic/go-agent"
	"github.com/server-may-cry/bubble-go/market"
	"github.com/server-may-cry/bubble-go/notification"
)

var newrelicApp newrelic.Application
var statickHandler StatickHandler

func init() {
	newrelicKey := os.Getenv("NEW_RELIC_LICENSE_KEY")
	if newrelicKey == "" {
		newrelicKey = "1234567890123456789012345678901234567890" // length 40
	}
	config := newrelic.NewConfig("bubble-go", newrelicKey)
	app, err := newrelic.NewApplication(config)
	if err != nil {
		log.Fatal("Can't create newrelic instance", err)
	}
	newrelicApp = app
}

// GetRouter return http.Handler for tests and core
func GetRouter(test bool, db *gorm.DB, marketInstance *market.Market, vkWorker *notification.VkWorker) http.Handler {
	router := chi.NewRouter()

	if !test {
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)
	}

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	router.Mount("/debug", middleware.Profiler())

	router.Get(wrapHandlerFunc("/", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, jsonHelper{
			"foo": "bar",
		})
	}))

	router.Mount("/", func() http.Handler {
		r := chi.NewRouter()
		r.Use(AuthorizationMiddleware(db))
		r.Post(wrapHandlerFunc("/ReqEnter", ReqEnter(db)))
		r.Post(wrapHandlerFunc("/ReqBuyProduct", ReqBuyProduct(db, marketInstance)))
		r.Post(wrapHandlerFunc("/ReqReduceTries", ReqReduceTries(db)))
		r.Post(wrapHandlerFunc("/ReqReduceCredits", ReqReduceCredits(db)))
		r.Post(wrapHandlerFunc("/ReqSavePlayerProgress", ReqSavePlayerProgress(vkWorker, db)))
		r.Post(wrapHandlerFunc("/ReqUsersProgress", ReqUsersProgress(db)))
		return r
	}())
	router.Post(wrapHandlerFunc("/VkPay", VkPay(db, marketInstance)))

	router.Get(wrapHandlerFunc("/crossdomain.xml", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<?xml version="1.0"?><cross-domain-policy><allow-access-from domain="*" /></cross-domain-policy>`))
	}))
	// http://119226.selcdn.ru/bubble/ShootTheBubbleDevVK.html
	// http://bubble-srv-dev.herokuapp.com/bubble/ShootTheBubbleDevVK.html
	statickHandler, err := NewStatickHandler("http://119226.selcdn.ru")
	if err != nil {
		log.Fatal("can't create static handler", err)
	}
	router.Get(wrapHandlerFunc("/bubble/*filePath", statickHandler.Serve))
	router.Get(wrapHandlerFunc("/cache-clear", statickHandler.Clear))

	router.Get(wrapHandlerFunc("/exception", func(w http.ResponseWriter, r *http.Request) {
		panic("test log.Fatal")
	}))

	router.Get(wrapHandlerFunc("/debug-vk", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, jsonHelper{
			"levels": vkWorker.LenLevels(),
			"events": vkWorker.LenEvents(),
		})
	}))

	loaderio := os.Getenv("LOADERIO")
	loaderioRoute := fmt.Sprintf("/loaderio-%s/", loaderio)
	router.Get(wrapHandlerFunc(loaderioRoute, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("loaderio-" + loaderio))
	}))

	return router
}

func wrapHandlerFunc(route string, handler http.HandlerFunc) (string, http.HandlerFunc) {
	return newrelic.WrapHandleFunc(newrelicApp, route, handler)
}
