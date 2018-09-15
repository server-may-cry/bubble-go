package application

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	newrelic "github.com/newrelic/go-agent"
	"github.com/server-may-cry/bubble-go/errorlisteners/sentry"
	"github.com/server-may-cry/bubble-go/mynewrelic"
	"github.com/server-may-cry/bubble-go/notification"
	dig "go.uber.org/dig"
)

// RouterDependencies for uber-go/dig
type RouterDependencies struct {
	dig.In

	HandlerSecure           []HTTPHandler `group:"server"`
	Newrelic                newrelic.Application
	AuthorizationMiddleware Middleware
	NewrelicMiddleware      mynewrelic.Middleware
	StaticHandler           *StaticHandler
	VkWorker                *notification.VkWorker
	VkPayHandler            http.HandlerFunc
	Test                    bool
	Version                 string
}

// GetRouter return http.Handler for tests and core
func GetRouter(deps RouterDependencies) http.Handler {
	router := chi.NewRouter()
	router.Use(deps.NewrelicMiddleware)

	if !deps.Test {
		router.Use(middleware.Logger)

		// improved middleware.Recoverer to send error into newrelic
		router.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if rvr := recover(); rvr != nil {
						var err error
						switch v := rvr.(type) {
						case error:
							err = v
						case *net.OpError:
							err = errors.New(v.Error())
						default:
							err = errors.New(rvr.(string))
						}
						r.Context().Value(mynewrelic.Ctx).(newrelic.Transaction).NoticeError(err)
						sentry.HandleError(err)
						logEntry := middleware.GetLogEntry(r)
						if logEntry != nil {
							logEntry.Panic(rvr, debug.Stack())
						} else {
							fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
							debug.PrintStack()
						}

						http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					}
				}()
				next.ServeHTTP(w, r)
			})
		})
	}

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	router.Mount("/debug", middleware.Profiler())

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, jsonHelper{
			"foo":     "bar",
			"version": deps.Version,
		})
	})

	router.Mount("/", func() http.Handler {
		r := chi.NewRouter()
		r.Use(deps.AuthorizationMiddleware)
		for _, handler := range deps.HandlerSecure {
			r.Post(handler.URL, handler.HTTPHandler)
		}
		return r
	}())
	router.Post("/VkPay", deps.VkPayHandler)

	router.Get("/crossdomain.xml", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<?xml version="1.0"?><cross-domain-policy><allow-access-from domain="*" /></cross-domain-policy>`))
	})
	router.Get("/bubble/*filePath", deps.StaticHandler.Serve)
	router.Get("/cache-clear", deps.StaticHandler.Clear)

	router.Get("/exception", func(w http.ResponseWriter, r *http.Request) {
		panic("test log.Fatal")
	})

	router.Get("/debug-vk", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, jsonHelper{
			"levels": deps.VkWorker.LenLevels(),
			"events": deps.VkWorker.LenEvents(),
		})
	})

	loaderio := os.Getenv("LOADERIO")
	loaderioRoute := fmt.Sprintf("/loaderio-%s/", loaderio)
	router.Get(loaderioRoute, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("loaderio-" + loaderio))
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		_ = r.Context().Value(mynewrelic.Ctx).(newrelic.Transaction).SetName("/404")
		http.NotFound(w, r)
	})

	return router
}
