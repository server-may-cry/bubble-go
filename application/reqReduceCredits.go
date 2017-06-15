package application

import (
	"encoding/json"
	"net/http"
)

type reduceCreditsRequest struct {
	baseRequest
	Amount int `json:"amount,string"`
}

type reduceCreditsResponse struct {
	ReqMsgID string `json:"reqMsgId"`
	UserID   string `json:"userId"`
	Credits  int    `json:"credits"`
}

// ReqReduceCredits reduce user credits. Get amount from request
func ReqReduceCredits(w http.ResponseWriter, r *http.Request) {
	request := reduceCreditsRequest{}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if request.Amount < 0 {
		return
	}
	ctx := r.Context()
	user := ctx.Value(userCtxID).(User)
	user.Credits -= request.Amount
	Gorm.Save(&user)
	response := reduceCreditsResponse{
		Credits: user.Credits,
	}
	JSON(w, response)
}
