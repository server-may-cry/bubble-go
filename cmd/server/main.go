package main

import (
	"log"
	"net/url"
	"os"

	gorelic "github.com/brandfolder/gin-gorelic"
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"github.com/server-may-cry/bubble-go/controllers"
	"github.com/server-may-cry/bubble-go/storage"
	"gopkg.in/redis.v4"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	rawRedisURL := os.Getenv("REDIS_URL")
	if rawRedisURL == "" {
		log.Fatal("$REDIS_URL must be set")
	}
	redisURL, err := url.Parse(rawRedisURL)
	if err != nil {
		log.Fatal(err)
	}
	password, _ := redisURL.User.Password()
	storage.Redis = redis.NewClient(&redis.Options{
		Addr:     redisURL.Host,
		Password: password,
		DB:       0, // default database
	})
}

// GetEngine return gin engine instance same for tests and server
func GetEngine() *gin.Engine {
	r := gin.Default()

	socketioServer, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	socketioServer.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		so.Join("chat")
		so.On("chat message", func(msg string) {
			log.Println("emit:", so.Emit("chat message", msg))
			so.BroadcastTo("chat", "chat message", msg)
		})
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	socketioServer.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	if "release" == os.Getenv("GIN_MODE") {
		newRelicLicenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY")
		if newRelicLicenseKey == "" {
			log.Fatal("$NEW_RELIC_LICENSE_KEY must be set")
		}
		gorelic.InitNewrelicAgent(newRelicLicenseKey, "bubble-go", true)
		r.Use(gorelic.Handler)
	}
	r.Use(gin.Logger())

	r.GET("/", controllers.Index)
	r.GET("/test", controllers.Test)
	r.GET("/redis", controllers.Redis)
	r.GET("/socket.io/", gin.WrapH(socketioServer))

	return r
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	GetEngine().Run(":" + port)
}
