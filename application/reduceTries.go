package application

import (
	"encoding/json"
	"net/http"

	"github.com/server-may-cry/bubble-go/storage"
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
		http.Error(w, getErrBody(err), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	user := ctx.Value(userCtxID).(User)
	user.RemainingTries--
	storage.Gorm.Save(&user)
	response := make([]int8, 1)
	response[0] = user.RemainingTries
	JSON(w, response)
}
