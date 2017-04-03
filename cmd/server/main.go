package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/server-may-cry/bubble-go/controllers"
	"github.com/server-may-cry/bubble-go/market"
	"github.com/server-may-cry/bubble-go/middleware"
	"github.com/server-may-cry/bubble-go/models"
	"github.com/server-may-cry/bubble-go/storage"
	"gopkg.in/gin-gonic/gin.v1"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		log.Fatalf("failed to connect database: %s", err)
	}
	db.AutoMigrate(&models.User{})
	storage.Gorm = db

	market.InitializeMarket()
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
	router.POST("/pay/:platform", controllers.PayPlatform) // vk|ok

	router.GET("/crossdomain.xml", func(c *gin.Context) {
		c.String(http.StatusOK, "<?xml version=\"1.0\"?><cross-domain-policy><allow-access-from domain=\"*\" /></cross-domain-policy>")
	})
	// http://119226.selcdn.ru/bubble/ShootTheBubbleDevVK.html
	// http://bubble-srv-dev.herokuapp.com/bubble/ShootTheBubbleDevVK.html
	router.GET("/bubble/*filePath", controllers.ServeStatick)
	router.GET("/cache-clear", controllers.ClearStatickCache)

	router.GET("/exception", func(c *gin.Context) {
		log.Fatal("test log.Fatal")
	})

	loaderio := os.Getenv("LOADERIO")
	loaderioRoute := fmt.Sprintf("/loaderio-%s", loaderio)
	router.GET(loaderioRoute, func(c *gin.Context) {
		c.String(http.StatusOK, fmt.Sprintf("loaderio-%s", loaderio))
	})

	router.Run()
}
