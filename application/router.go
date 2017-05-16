package application

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

// GetRouter return http.Handler for tests and core
func GetRouter(test bool) http.Handler {
	router := chi.NewRouter()

	if !test {
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)
	}
	router.Mount("/debug", middleware.Profiler())

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, h{
			"foo": "bar",
		})
	})

	router.Mount("/", func() http.Handler {
		r := chi.NewRouter()
		r.Use(AuthorizationMiddleware)
		r.Post("/ReqEnter", ReqEnter)
		r.Post("/ReqBuyProduct", ReqBuyProduct)
		r.Post("/ReqReduceTries", ReqReduceTries)
		r.Post("/ReqReduceCredits", ReqReduceCredits)
		r.Post("/ReqSavePlayerProgress", ReqSavePlayerProgress)
		r.Post("/ReqUsersProgress", ReqUsersProgress)
		return r
	}())
	router.Post("/VkPay", VkPay)

	router.Get("/crossdomain.xml", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<?xml version="1.0"?><cross-domain-policy><allow-access-from domain="*" /></cross-domain-policy>`))
	})
	// http://119226.selcdn.ru/bubble/ShootTheBubbleDevVK.html
	// http://bubble-srv-dev.herokuapp.com/bubble/ShootTheBubbleDevVK.html
	router.Get("/bubble/*filePath", ServeStatick)
	router.Get("/cache-clear", ClearStatickCache)

	router.Get("/exception", func(w http.ResponseWriter, r *http.Request) {
		panic("test log.Fatal")
	})

	router.Get("/debug-vk", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, h{
			"levels": VkWorker.Levels,
			"events": VkWorker.Events,
		})
	})

	loaderio := os.Getenv("LOADERIO")
	loaderioRoute := fmt.Sprintf("/loaderio-%s", loaderio)
	router.Get(loaderioRoute, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("loaderio-%s", loaderio)))
	})

	return router
}
