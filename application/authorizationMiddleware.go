package application

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/server-may-cry/bubble-go/platforms"
)

// AuthorizationMiddleware check signature and load user
func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, _ := ioutil.ReadAll(r.Body)
		requestBodyCopy := ioutil.NopCloser(bytes.NewBuffer(buf))
		r.Body = requestBodyCopy

		request := AuthRequestPart{}
		if err := json.Unmarshal(buf, &request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var stringToHash string
		platformID := platforms.GetByName(request.SysID)
		switch request.SysID {
		case "VK":
			stringToHash = fmt.Sprintf(
				"%s_%d_%s",
				os.Getenv("VK_APP_ID"),
				request.ExtID,
				os.Getenv("VK_SECRET"),
			)
		case "OK":
			stringToHash = fmt.Sprintf(
				"%d%s%s",
				request.ExtID,
				request.SessionKey,
				os.Getenv("OK_SECRET"),
			)
		default:
			http.Error(w, fmt.Sprintf("Unknown platform %s", request.SysID), http.StatusBadRequest)
			return
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
			http.Error(w, fmt.Sprintf("Bad auth key %s", request.AuthKey), 403)
			return
		}

		db := Gorm
		var user User
		db.Where("sys_id = ? AND ext_id = ?", platformID, request.ExtID).First(&user)
		if user.ID != 0 { // check user exists
			ctx := context.WithValue(r.Context(), userCtxID, user)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
