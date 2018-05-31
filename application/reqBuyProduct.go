package application

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/server-may-cry/bubble-go/market"
)

type buyProductRequest struct {
	baseRequest
	ProductID string `json:"productId"`
}

type buyProductResponse struct {
	ProductID string `json:"productId"`
	Credits   int    `json:"credits,uint16"`
}

// ReqBuyProduct buy product
func ReqBuyProduct(db *gorm.DB, marketInstance *market.Market) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
		db.Save(&user)
		JSON(w, buyProductResponse{
			ProductID: request.ProductID,
			Credits:   user.Credits,
		})
	}
}
