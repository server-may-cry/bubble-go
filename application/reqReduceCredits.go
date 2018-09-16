package application

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
	newrelic "github.com/newrelic/go-agent"
	"github.com/server-may-cry/bubble-go/models"
	"github.com/server-may-cry/bubble-go/mynewrelic"
)

type reduceCreditsRequest struct {
	Amount int `json:"amount,string"`
}

type reduceCreditsResponse struct {
	ReqMsgID string `json:"reqMsgId"`
	UserID   string `json:"userId"`
	Credits  int    `json:"credits"`
}

// ReqReduceCredits reduce user credits. Get amount from request
func ReqReduceCredits(db *gorm.DB) HTTPHandlerContainer {
	handler := HTTPHandler{
		URL: "/ReqReduceCredits",
	}
	handler.HTTPHandler = func(w http.ResponseWriter, r *http.Request) {
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
		user := r.Context().Value(userCtxID).(models.User)
		user.Credits -= request.Amount

		s := newrelic.DatastoreSegment{
			StartTime:  newrelic.StartSegmentNow(r.Context().Value(mynewrelic.Ctx).(newrelic.Transaction)),
			Product:    newrelic.DatastorePostgres,
			Collection: "user",
			Operation:  "UPDATE",
		}
		db.Save(&user)
		_ = s.End()

		JSON(w, reduceCreditsResponse{
			Credits: user.Credits,
		})
	}

	return HTTPHandlerContainer{
		HTTPHandler: handler,
	}
}
