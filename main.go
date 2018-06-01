package main

import (
	"encoding/json"
	"flag"
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
	"github.com/pkg/errors"
	"github.com/server-may-cry/bubble-go/application"
	"github.com/server-may-cry/bubble-go/market"
	"github.com/server-may-cry/bubble-go/notification"
)

func initialize(pathToMarketConfig string, dbURL string) (
	*gorm.DB,
	*market.Market,
	*notification.VkWorker,
	error,
) {
	fail := func(err error) (*gorm.DB, *market.Market, *notification.VkWorker, error) {
		return nil, nil, nil, err
	}
	u, err := url.Parse(dbURL)
	if err != nil {
		return fail(errors.Wrapf(err, "can`t parse DB_URL (%s)", dbURL))
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
		return fail(errors.Wrap(err, "failed to connect database"))
	}
	db.AutoMigrate(&application.User{})
	db.AutoMigrate(&application.Transaction{})
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(17) // 20 actual limit

	file, err := os.Open(pathToMarketConfig)
	if err != nil {
		return fail(errors.Wrap(err, "can`t open market config file"))
	}
	var marketConfig market.Config
	err = json.NewDecoder(file).Decode(&marketConfig)
	if err != nil {
		return fail(errors.Wrap(err, "can`t decode market config"))
	}
	marketInstance := market.NewMarket(marketConfig, os.Getenv("CDN_PREFIX"))
	user := application.User{}
	err = marketInstance.Validate(&user)
	if err != nil {
		return fail(errors.Wrap(err, "invalid marker configuration"))
	}

	vkWorker := notification.NewVkWorker(notification.VkConfig{
		AppID:           os.Getenv("VK_APP_ID"),
		Secret:          os.Getenv("VK_SECRET"),
		RequestInterval: time.Millisecond * 300,
	}, &http.Client{})
	err = vkWorker.Initialize()
	if err != nil {
		return fail(errors.Wrap(err, "can`t initialize vk worker"))
	}
	return db, marketInstance, vkWorker, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	dbURL := flag.String("db_credentials", os.Getenv("DB_URL"), "in URL")
	pathToMarketConfig := flag.String("market_config", filepath.ToSlash("./config/market.json"), "")
	port := flag.String("port", os.Getenv("PORT"), "port to listen for http server")
	flag.Parse()
	db, marketInstance, vkWorker, err := initialize(*pathToMarketConfig, *dbURL)
	if err != nil {
		log.Fatal("can't initialize application", err)
	}
	go vkWorker.Work()
	router := application.GetRouter(false, db, marketInstance, vkWorker)
	log.Println("ready")
	if err := http.ListenAndServe(":"+*port, router); err != nil {
		log.Fatal("http server failure", err)
	}
}
