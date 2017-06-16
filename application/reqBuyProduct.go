package application

import (
	"encoding/json"
	"net/http"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	user := ctx.Value(userCtxID).(User)
	err = Market.Buy(&user, request.ProductID)
	if err != nil {
		panic(err)
	}
	response := buyProductResponse{
		ProductID: request.ProductID,
		Credits:   user.Credits,
	}
	Gorm.Save(&user)
	JSON(w, response)
}
