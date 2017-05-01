package application

import (
	"encoding/json"
	"net/http"
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
	Gorm.Save(&user)
	response := make([]int8, 1)
	response[0] = user.RemainingTries
	JSON(w, response)
}
