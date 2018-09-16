package application

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
	newrelic "github.com/newrelic/go-agent"
	"github.com/server-may-cry/bubble-go/models"
	"github.com/server-may-cry/bubble-go/mynewrelic"
)

type usersProgressRequest struct {
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
func ReqUsersProgress(db *gorm.DB) HTTPHandlerContainer {
	handler := HTTPHandler{
		URL: "/ReqUsersProgress",
	}
	handler.HTTPHandler = func(w http.ResponseWriter, r *http.Request) {
		request := usersProgressRequest{}
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user := r.Context().Value(userCtxID).(models.User)
		usersLen := len(request.SocIDs)
		users := make([]models.User, usersLen)

		s := newrelic.DatastoreSegment{
			StartTime:  newrelic.StartSegmentNow(r.Context().Value(mynewrelic.Ctx).(newrelic.Transaction)),
			Product:    newrelic.DatastorePostgres,
			Collection: "user",
			Operation:  "SELECT",
		}
		db.Where("sys_id = ? and ext_id in (?)", user.SysID, request.SocIDs).Find(&users)
		_ = s.End()

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

	return HTTPHandlerContainer{
		HTTPHandler: handler,
	}
}
