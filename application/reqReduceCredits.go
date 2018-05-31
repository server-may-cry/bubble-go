package application

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
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
func ReqReduceCredits(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request := reduceCreditsRequest{}
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if request.Amount < 0 {
			http.Error(w, "amount to low", http.StatusBadRequest)
			return
		}
		user := r.Context().Value(userCtxID).(User)
		user.Credits -= request.Amount
		db.Save(&user)
		JSON(w, reduceCreditsResponse{
			Credits: user.Credits,
		})
	}
}
