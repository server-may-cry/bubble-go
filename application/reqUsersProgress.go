package application

import (
	"encoding/json"
	"net/http"
)

type usersProgressRequest struct {
	baseRequest
	SocIDs []uint64 `json:"socIds,[]uint64"`
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

// ReqUsersProgress return progres of recieved users
func ReqUsersProgress(w http.ResponseWriter, r *http.Request) {
	request := usersProgressRequest{}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, getErrBody(err), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	user := ctx.Value(userCtxID).(User)
	usersLen := len(request.SocIDs)
	users := make([]User, usersLen)
	Gorm.Where("sys_id = ? and ext_id in (?)", user.SysID, request.SocIDs).Find(&users)
	response := usersProgressResponse{
		UsersProgress: make([]userProgres, usersLen),
	}
	for i, friend := range users {
		response.UsersProgress[i] = userProgres{
			UserID:            friend.ID,
			SocID:             friend.ExtID,
			ReachedStage01:    friend.ReachedStage01,
			ReachedSubStage01: friend.ReachedSubStage01,
		}
	}
	JSON(w, response)
}
