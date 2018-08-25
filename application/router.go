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
	StaticHandler           *StaticHandler
	VkWorker                *notification.VkWorker
	VkPayHandler            http.HandlerFunc
	Test                    bool
	Version                 string
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

	wrapHandlerFunc := func(newrelicApp newrelic.Application) func(route string, handler http.HandlerFunc) (string, http.HandlerFunc) {
		return func(route string, handler http.HandlerFunc) (string, http.HandlerFunc) {
			return newrelic.WrapHandleFunc(newrelicApp, route, handler)
		}
	}(deps.Newrelic)

	router.Get(wrapHandlerFunc("/", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, jsonHelper{
			"foo":     "bar",
			"version": deps.Version,
		})
	}))

	router.Mount("/", func() http.Handler {
		r := chi.NewRouter()
		r.Use(deps.AuthorizationMiddleware)
		for _, handler := range deps.HandlerSecure {
			r.Post(wrapHandlerFunc(handler.URL, handler.HTTPHandler))
		}
		return r
	}())
	router.Post(wrapHandlerFunc("/VkPay", deps.VkPayHandler))

	router.Get(wrapHandlerFunc("/crossdomain.xml", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<?xml version="1.0"?><cross-domain-policy><allow-access-from domain="*" /></cross-domain-policy>`))
	}))
	router.Get(wrapHandlerFunc("/bubble/*filePath", deps.StaticHandler.Serve))
	router.Get(wrapHandlerFunc("/cache-clear", deps.StaticHandler.Clear))

	router.Get(wrapHandlerFunc("/exception", func(w http.ResponseWriter, r *http.Request) {
		panic("test log.Fatal")
	}))

	router.Get(wrapHandlerFunc("/debug-vk", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, jsonHelper{
			"levels": deps.VkWorker.LenLevels(),
			"events": deps.VkWorker.LenEvents(),
		})
	}))

	loaderio := os.Getenv("LOADERIO")
	loaderioRoute := fmt.Sprintf("/loaderio-%s/", loaderio)
	router.Get(wrapHandlerFunc(loaderioRoute, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("loaderio-" + loaderio))
	}))

	return router
}
