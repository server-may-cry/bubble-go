package application

import (
	"encoding/json"
	"net/http"

	"github.com/server-may-cry/bubble-go/market"
	"github.com/server-may-cry/bubble-go/storage"
)

type buyProductRequest struct {
	baseRequest
	ProductID string `json:"productId" binding:"required"`
}

type buyProductResponse struct {
	ProductID string `json:"productId"`
	Credits   int16  `json:"credits,uint16"`
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
	user := ctx.Value(UserCtxID).(User)
	market.Buy(&user, request.ProductID)
	response := buyProductResponse{
		ProductID: request.ProductID,
		Credits:   user.Credits,
	}
	storage.Gorm.Save(&user)
	JSON(w, response)
}
