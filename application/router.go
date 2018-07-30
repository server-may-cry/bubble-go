package application

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	newrelic "github.com/newrelic/go-agent"
	"github.com/server-may-cry/bubble-go/notification"
	dig "go.uber.org/dig"
)

// RouterDependencies for uber-go/dig
type RouterDependencies struct {
	dig.In

	HandlerSecure           []HTTPHandler `group:"server"`
	Newrelic                newrelic.Application
	AuthorizationMiddleware Middleware
	StatickHandler          *StatickHandler
	VkWorker                *notification.VkWorker
	VkPayHandler            http.HandlerFunc
	Test                    bool
}

// GetRouter return http.Handler for tests and core
func GetRouter(deps RouterDependencies) http.Handler {
	router := chi.NewRouter()

	if !deps.Test {
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)
	}

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	router.Mount("/debug", middleware.Profiler())

	router.Get(wrapHandlerFunc(deps.Newrelic, "/", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, jsonHelper{
			"foo": "bar",
		})
	}))

	router.Mount("/", func() http.Handler {
		r := chi.NewRouter()
		r.Use(deps.AuthorizationMiddleware)
		for _, handler := range deps.HandlerSecure {
			r.Post(wrapHandlerFunc(deps.Newrelic, handler.URL, handler.HTTPHandler))
		}
		return r
	}())
	router.Post(wrapHandlerFunc(deps.Newrelic, "/VkPay", deps.VkPayHandler))

	router.Get(wrapHandlerFunc(deps.Newrelic, "/crossdomain.xml", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<?xml version="1.0"?><cross-domain-policy><allow-access-from domain="*" /></cross-domain-policy>`))
	}))
	router.Get(wrapHandlerFunc(deps.Newrelic, "/bubble/*filePath", deps.StatickHandler.Serve))
	router.Get(wrapHandlerFunc(deps.Newrelic, "/cache-clear", deps.StatickHandler.Clear))

	router.Get(wrapHandlerFunc(deps.Newrelic, "/exception", func(w http.ResponseWriter, r *http.Request) {
		panic("test log.Fatal")
	}))

	router.Get(wrapHandlerFunc(deps.Newrelic, "/debug-vk", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, jsonHelper{
			"levels": deps.VkWorker.LenLevels(),
			"events": deps.VkWorker.LenEvents(),
		})
	}))

	loaderio := os.Getenv("LOADERIO")
	loaderioRoute := fmt.Sprintf("/loaderio-%s/", loaderio)
	router.Get(wrapHandlerFunc(deps.Newrelic, loaderioRoute, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("loaderio-" + loaderio))
	}))

	return router
}

func wrapHandlerFunc(newrelicApp newrelic.Application, route string, handler http.HandlerFunc) (string, http.HandlerFunc) {
	return newrelic.WrapHandleFunc(newrelicApp, route, handler)
}
