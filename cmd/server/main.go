package main

import (
	"crypto/md5"
	"fmt"
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
	securedGroup.Use(signatureValidatorMiddleware)
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
		c.String(http.StatusOK, loaderioRoute)
	})

	router.Run()
}

func signatureValidatorMiddleware(c *gin.Context) {
	request := controllers.AuthRequestPart{}
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var stringToHash string
	switch request.SysID {
	case "VK":
		appID := os.Getenv("VK_APP_ID")
		secret := os.Getenv("VK_SECRET")
		stringToHash = fmt.Sprintf("%s_%s_%s", appID, request.ExtID, secret)
	case "OK":
		secret := os.Getenv("OK_SECRET")
		stringToHash = fmt.Sprintf("%s%s%s", request.ExtID, request.SessionKey, secret)
	default:
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Unknown platform %s", request.SysID))
	}
	data := []byte(stringToHash)
	expectedMD5 := md5.Sum(data)
	expectedAuthKey := fmt.Sprintf("%x", expectedMD5)
	if expectedAuthKey != request.AuthKey {
		log.Print("authorization failure", " ", stringToHash, " ", expectedAuthKey, " ", request.AuthKey)
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Bad auth key %s", request.AuthKey))
	}

	log.Print("authorization success")
	c.Set("user", request.ExtID) // TODO
	c.Next()
}
