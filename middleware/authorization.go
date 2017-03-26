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
	"github.com/server-may-cry/bubble-go/storage"
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
	var platformID uint8
	switch request.SysID {
	case "VK":
		platformID = 1
		stringToHash = fmt.Sprintf(
			"%s_%s_%s",
			os.Getenv("VK_APP_ID"),
			request.ExtID,
			os.Getenv("VK_SECRET"),
		)
	case "OK":
		platformID = 2
		stringToHash = fmt.Sprintf(
			"%s%s%s",
			request.ExtID,
			request.SessionKey,
			os.Getenv("OK_SECRET"),
		)
	default:
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Unknown platform %s", request.SysID))
	}
	data := []byte(stringToHash)
	expectedAuthKey := fmt.Sprintf("%x", md5.Sum(data))
	if expectedAuthKey != request.AuthKey {
		log.Printf(
			"authorization failure %s %s %s",
			stringToHash,
			expectedAuthKey,
			request.AuthKey,
		)
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Bad auth key %s", request.AuthKey))
		return
	}

	log.Print("authorization success")
	db := storage.Gorm
	var user models.User
	db.First(&user, "sysID = ? AND extID = ?", platformID, request.ExtID)
	if user.SysID == platformID { // check user exists
		c.Set("user", user)
	}
	c.Next()
}
