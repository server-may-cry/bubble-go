package container

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgresql driver
	newrelic "github.com/newrelic/go-agent"
	"github.com/pkg/errors"
	"github.com/server-may-cry/bubble-go/application"
	"github.com/server-may-cry/bubble-go/market"
	"github.com/server-may-cry/bubble-go/notification"
	"go.uber.org/dig"
)

// Configuration initialization configuration
type Configuration struct {
	dbURL              string
	pathToMarketConfig string
}

// Get create DI container
func Get(pathToMarketConfig string, dbURL string, test bool, version string) *dig.Container {
	container := dig.New()

	_ = container.Provide(func() bool {
		return test
	})

	_ = container.Provide(func() string {
		return version
	})

	_ = container.Provide(application.AuthorizationMiddleware)
	_ = container.Provide(application.ReqBuyProduct)
	_ = container.Provide(application.ReqEnter)
	_ = container.Provide(application.ReqReduceCredits)
	_ = container.Provide(application.ReqReduceTries)
	_ = container.Provide(application.ReqSavePlayerProgress)
	_ = container.Provide(application.ReqUsersProgress)

	_ = container.Provide(func() Configuration {
		return Configuration{
			dbURL:              dbURL,
			pathToMarketConfig: pathToMarketConfig,
		}
	})
	_ = container.Provide(func(config Configuration) (*gorm.DB, error) {
		u, err := url.Parse(config.dbURL)
		if err != nil {
			return nil, errors.Wrapf(err, "can`t parse DB_URL (%s)", dbURL)
		}

		host, port, _ := net.SplitHostPort(u.Host)
		password, _ := u.User.Password()
		ssl := "enable"
		if u.Scheme == "http" {
			ssl = "disable"
		}
		connect := func() (*gorm.DB, error) {
			return gorm.Open("postgres", fmt.Sprintf(
				"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
				host,
				port,
				u.User.Username(),
				password,
				strings.Trim(u.Path, "/"),
				ssl,
			))
		}
		db, err := connect()
		if err != nil {
			log.Println("DB connection failed. Wait 5 seconds before repeat.")
			time.Sleep(5 * time.Second)
			db, err = connect()
			if err != nil {
				return nil, errors.Wrap(err, "failed to connect database")
			}
		}
		db.AutoMigrate(&application.User{})
		db.AutoMigrate(&application.Transaction{})
		db.DB().SetMaxIdleConns(5)
		db.DB().SetMaxOpenConns(18) // 20 actual limit

		return db, nil
	})

	_ = container.Provide(func(config Configuration) (*market.Market, error) {
		file, err := os.Open(config.pathToMarketConfig)
		if err != nil {
			return nil, errors.Wrap(err, "can`t open market config file")
		}
		var marketConfig market.Config
		err = json.NewDecoder(file).Decode(&marketConfig)
		if err != nil {
			return nil, errors.Wrap(err, "can`t decode market config")
		}
		marketInstance := market.NewMarket(marketConfig, os.Getenv("CDN_PREFIX"))
		user := application.User{}
		err = marketInstance.Validate(&user)
		if err != nil {
			return nil, errors.Wrap(err, "invalid marker configuration")
		}
		return marketInstance, nil
	})

	_ = container.Provide(func() (*notification.VkWorker, error) {
		vkWorker := notification.NewVkWorker(notification.VkConfig{
			AppID:           os.Getenv("VK_APP_ID"),
			Secret:          os.Getenv("VK_SECRET"),
			RequestInterval: time.Millisecond * 300,
		}, &http.Client{})
		err := vkWorker.Initialize()
		if err != nil {
			return nil, errors.Wrap(err, "can`t initialize vk worker")
		}

		return vkWorker, nil
	})

	_ = container.Provide(func() (*application.StatickHandler, error) {
		return application.NewStatickHandler("http://119226.selcdn.ru")
	})

	_ = container.Provide(func() (newrelic.Application, error) {
		newrelicKey := os.Getenv("NEW_RELIC_LICENSE_KEY")
		if newrelicKey == "" {
			newrelicKey = "1234567890123456789012345678901234567890" // length 40
		}
		config := newrelic.NewConfig("bubble-go", newrelicKey)
		app, err := newrelic.NewApplication(config)
		if err != nil {
			return nil, errors.Wrap(err, "Can't create newrelic instance")
		}

		return app, err
	})

	_ = container.Provide(application.GetRouter)
	_ = container.Provide(application.VkPay)

	return container
}
