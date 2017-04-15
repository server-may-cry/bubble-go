package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/server-may-cry/bubble-go/models"
	"github.com/server-may-cry/bubble-go/storage"
)

type reduceCreditsRequest struct {
	baseRequest
	Amount int16 `json:"amount,string" binding:"required"`
}

type reduceCreditsResponse struct {
	ReqMsgID string `json:"reqMsgId"`
	UserID   string `json:"userId"`
	Credits  int16  `json:"credits"`
}

// ReqReduceCredits reduce user credits. Get amount from request
func ReqReduceCredits(w http.ResponseWriter, r *http.Request) {
	request := reduceCreditsRequest{}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, getErrBody(err), http.StatusBadRequest)
		return
	}
	if request.Amount < 0 {
		return
	}
	ctx := r.Context()
	user := ctx.Value(User).(models.User)
	user.Credits -= request.Amount
	storage.Gorm.Save(&user)
	response := reduceCreditsResponse{
		Credits: user.Credits,
	}
	JSON(w, response)
}
