package application

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
	newrelic "github.com/newrelic/go-agent"
	"github.com/server-may-cry/bubble-go/market"
	"github.com/server-may-cry/bubble-go/mynewrelic"
)

type buyProductRequest struct {
	ProductID string `json:"productId"`
}

type buyProductResponse struct {
	ProductID string `json:"productId"`
	Credits   int    `json:"credits,uint16"`
}

// ReqBuyProduct buy product
func ReqBuyProduct(db *gorm.DB, marketInstance *market.Market) HTTPHandlerContainer {
	handler := HTTPHandler{
		URL: "/ReqBuyProduct",
	}
	handler.HTTPHandler = func(w http.ResponseWriter, r *http.Request) {
		request := buyProductRequest{}
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user := r.Context().Value(userCtxID).(User)
		err = marketInstance.Buy(&user, request.ProductID)
		if err != nil {
			panic(err)
		}

		s := newrelic.DatastoreSegment{
			StartTime:  newrelic.StartSegmentNow(r.Context().Value(mynewrelic.Ctx).(newrelic.Transaction)),
			Product:    newrelic.DatastorePostgres,
			Collection: "user",
			Operation:  "UPDATE",
		}
		db.Save(&user)
		_ = s.End()

		JSON(w, buyProductResponse{
			ProductID: request.ProductID,
			Credits:   user.Credits,
		})
	}

	return HTTPHandlerContainer{
		HTTPHandler: handler,
	}
}
