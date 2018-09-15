package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/getsentry/raven-go"

	"github.com/server-may-cry/bubble-go/container"
	"github.com/server-may-cry/bubble-go/notification"
)

var version = "undefined"

func main() {
	var err error
	log.SetFlags(log.LstdFlags | log.Llongfile)
	dbURL := flag.String("db_credentials", os.Getenv("HEROKU_POSTGRESQL_SILVER_URL"), "in URL")
	pathToMarketConfig := flag.String("market_config", filepath.ToSlash("./config/market.json"), "")
	port := flag.String("port", os.Getenv("PORT"), "port to listen for http server")
	fastShutdown := flag.Bool("fast_shutdown", false, "test that application can be initialized")
	sentryURL := flag.String("sentry_url", os.Getenv("SENTRY_DSN"), "in URL")
	flag.Parse()
	container := container.Get(*pathToMarketConfig, *dbURL, false, version)
	err = container.Invoke(func(worker *notification.VkWorker) {
		go worker.Work()
	})
	if err != nil {
		panic(err.Error() + "\n" + container.String())
	}
	err = raven.SetDSN(*sentryURL)
	if err != nil {
		panic(err)
	}
	raven.DefaultClient.SetRelease(version)
	err = container.Invoke(func(router http.Handler) {
		log.Println("ready")
		if *fastShutdown {
			return
		}
		if err := http.ListenAndServe(":"+*port, router); err != nil {
			log.Fatal("http server failure", err)
		}
	})
	if err != nil {
		panic(err.Error() + "\n" + container.String())
	}
}
