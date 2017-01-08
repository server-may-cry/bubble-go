package main

import (
	"log"
	"net/http"
	"os"

	"github.com/server-may-cry/bubble-go/controllers"
	"github.com/server-may-cry/bubble-go/storage"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	rawMongoURL := os.Getenv("MONGODB_URI")
	if rawMongoURL == "" {
		log.Fatal("$MONGODB_URI must be set")
	}
	mongoConnection, err := mgo.Dial(rawMongoURL)
	if err != nil {
		log.Fatal(err)
	}
	storage.MongoDB = mongoConnection
}

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"foo": "bar",
		})
	})

	securedGroup := router.Group("")
	{
		securedGroup.POST("/ReqEnter", controllers.ReqEnter)
		securedGroup.POST("/ReqBuyProduct", controllers.ReqBuyProduct)
		securedGroup.POST("/ReqReduceTries", controllers.ReqReduceTries)
		securedGroup.POST("/ReqReduceCredits", controllers.ReqReduceCredits)
		securedGroup.POST("/ReqSavePlayerProgress", controllers.ReqSavePlayerProgress)
		securedGroup.POST("/ReqUsersProgress", controllers.ReqUsersProgress)
	}
	router.POST("/pay/:platform", controllers.PayPlatform) // vk|ok
	router.POST("/bubble/*filePath", controllers.LoadStatick)
	router.GET("/cache-clear", controllers.ClearStatickCache)

	router.GET("/exception", func(c *gin.Context) {
		log.Fatal("test log.Fatal")
	})
	router.GET("/loaderio-some_hash", func(c *gin.Context) {
		c.String(http.StatusOK, "text/plain", "some_hash")
	})

	router.Run()
}
