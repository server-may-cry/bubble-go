package application

import (
	"encoding/json"
	"net/http"

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
func ReqBuyProduct(w http.ResponseWriter, r *http.Request) {
	request := buyProductRequest{}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, getErrBody(err), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	user := ctx.Value(userCtxID).(User)
	market.Buy(&user, request.ProductID)
	response := buyProductResponse{
		ProductID: request.ProductID,
		Credits:   user.Credits,
	}
	Gorm.Save(&user)
	JSON(w, response)
}
