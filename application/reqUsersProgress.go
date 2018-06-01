package application

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
)

type usersProgressRequest struct {
	baseRequest
	SocIDs []int64 `json:"socIds"`
}

type userProgres struct {
	UserID            uint64 `json:"userId"`
	SocID             int64  `json:"socId,string"`
	ReachedStage01    int8   `json:"reachedStage01"`
	ReachedStage02    int8   `json:"reachedStage02"`
	ReachedSubStage01 int8   `json:"reachedSubStage01"`
	ReachedSubStage02 int8   `json:"reachedSubStage02"`
}
type usersProgressResponse struct {
	UsersProgress []userProgres `json:"usersProgress"`
}

// ReqUsersProgress return progres of received users
func ReqUsersProgress(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request := usersProgressRequest{}
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user := r.Context().Value(userCtxID).(User)
		usersLen := len(request.SocIDs)
		users := make([]User, usersLen)
		db.Where("sys_id = ? and ext_id in (?)", user.SysID, request.SocIDs).Find(&users)
		response := usersProgressResponse{
			UsersProgress: make([]userProgres, usersLen),
		}
		for i, friend := range users {
			response.UsersProgress[i] = userProgres{
				UserID:            friend.ID,
				SocID:             friend.ExtID,
				ReachedStage01:    friend.ReachedStage01,
				ReachedSubStage01: friend.ReachedSubStage01,
				ReachedStage02:    friend.ReachedStage02,
				ReachedSubStage02: friend.ReachedSubStage02,
			}
		}
		JSON(w, response)
	}
}
