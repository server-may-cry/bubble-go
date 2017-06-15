package application

import (
	"encoding/json"
	"net/http"
	"time"
)

type reduceTriesRequest struct {
	baseRequest
}

// ReqReduceTries reduce user tries by one
func ReqReduceTries(w http.ResponseWriter, r *http.Request) {
	request := reduceTriesRequest{}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	user := ctx.Value(userCtxID).(User)
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
	Gorm.Save(&user)
	response := make([]int8, 1)
	response[0] = user.RemainingTries
	JSON(w, response)
}
