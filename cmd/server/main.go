package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/server-may-cry/bubble-go/controllers"
	"github.com/server-may-cry/bubble-go/market"
	"github.com/server-may-cry/bubble-go/middleware"
	"github.com/server-may-cry/bubble-go/models"
	"github.com/server-may-cry/bubble-go/notification"
	"github.com/server-may-cry/bubble-go/storage"
	"gopkg.in/gin-gonic/gin.v1"
)

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

	go notification.VkWorkerInit()
}

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"foo": "bar",
		})
	})

	securedGroup := router.Group("")
	securedGroup.Use(middleware.AuthorizationMiddleware)
	{
		securedGroup.POST("/ReqEnter", controllers.ReqEnter)
		securedGroup.POST("/ReqBuyProduct", controllers.ReqBuyProduct)
		securedGroup.POST("/ReqReduceTries", controllers.ReqReduceTries)
		securedGroup.POST("/ReqReduceCredits", controllers.ReqReduceCredits)
		securedGroup.POST("/ReqSavePlayerProgress", controllers.ReqSavePlayerProgress)
		securedGroup.POST("/ReqUsersProgress", controllers.ReqUsersProgress)
	}
	router.POST("/VkPay", controllers.VkPay)

	router.GET("/crossdomain.xml", func(c *gin.Context) {
		c.String(http.StatusOK, "<?xml version=\"1.0\"?><cross-domain-policy><allow-access-from domain=\"*\" /></cross-domain-policy>")
	})
	// http://119226.selcdn.ru/bubble/ShootTheBubbleDevVK.html
	// http://bubble-srv-dev.herokuapp.com/bubble/ShootTheBubbleDevVK.html
	router.GET("/bubble/*filePath", controllers.ServeStatick)
	router.GET("/cache-clear", controllers.ClearStatickCache)

	router.GET("/exception", func(c *gin.Context) {
		panic("test log.Fatal")
	})

	loaderio := os.Getenv("LOADERIO")
	loaderioRoute := fmt.Sprintf("/loaderio-%s", loaderio)
	router.GET(loaderioRoute, func(c *gin.Context) {
		c.String(http.StatusOK, fmt.Sprintf("loaderio-%s", loaderio))
	})

	router.Run()
}
