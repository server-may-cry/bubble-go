package middleware

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/server-may-cry/bubble-go/controllers"
	"github.com/server-may-cry/bubble-go/models"
	"gopkg.in/gin-gonic/gin.v1"
)

// AuthorizationMiddleware check signature and load user
func AuthorizationMiddleware(c *gin.Context) {
	// TODO try decoder := json.NewDecoder(c.Request.Body)
	buf, _ := ioutil.ReadAll(c.Request.Body)
	requestBodyCopy := ioutil.NopCloser(bytes.NewBuffer(buf))
	c.Request.Body = requestBodyCopy

	request := controllers.AuthRequestPart{}
	if err := json.Unmarshal(buf, &request); err != nil {
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
		log.Print("expected ", expectedAuthKey)
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Bad auth key %s", request.AuthKey))
		return
	}

	log.Print("authorization success")
	c.Set("user", models.User{}) // TODO
	c.Next()
}
