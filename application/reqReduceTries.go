package application

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
)

type reduceTriesRequest struct {
	baseRequest
}

// ReqReduceTries reduce user tries by one
func ReqReduceTries(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request := reduceTriesRequest{}
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user := r.Context().Value(userCtxID).(User)
		ts := time.Now().Unix()
		if user.RestoreTriesAt != 0 && user.RestoreTriesAt <= ts {
			if user.RemainingTries < defaultConfig.DefaultRemainingTries {
				user.RemainingTries = defaultConfig.DefaultRemainingTries
			}
		}
		user.RemainingTries--
		if user.RestoreTriesAt < 0 {
			user.RemainingTries = 0
		}
		if user.RestoreTriesAt == 0 {
			user.RestoreTriesAt = ts + int64(defaultConfig.IntervalTriesRestoration)
		}
		db.Save(&user)
		JSON(w, []int8{user.RemainingTries})
	}
}
