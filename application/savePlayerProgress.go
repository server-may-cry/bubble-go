package application

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/server-may-cry/bubble-go/storage"
)

/*
{
    ...
    "reachedSubStage":"1", // номер точки на острове до которой игрок дошел за все время игры
    "currentStage":"0", // номер острова на котором игрок прошел точку в текущей сессии
    "reachedStage":"0", // номер острова до которого игрок дошел за все время игры
    "completeSubStage":"1", // номер точки на острове которую игрок прошел в текущей сессии
    "completeSubStageRecordStat":"2", // количество звезд набранных на пройденной точке(перезаписывается только в случае если новое значение больше предыдущего)
    "levelMode":"standart", // режим игры, может принимать значения "standart" (01) и "arcade" (02)
    Если режим игры "standart", то перезаписываются значения reachedStage01 reachedSubStage01,
    которые приходят в ReqEnter`e, если же "arcade" то reachedStage02 и reachedSubStage02
}
*/
type savePlayerProgressRequest struct {
	baseRequest
	ReachedSubStage            int8   `json:"reachedSubStage,string" binding:"required"`
	CurrentStage               int8   `json:"currentStage,string" binding:"required"`
	ReachedStage               int8   `json:"reachedStage,string" binding:"required"`
	CompleteSubStage           int8   `json:"completeSubStage,string" binding:"required"`
	CompleteSubStageRecordStat int8   `json:"completeSubStageRecordStat,string" binding:"required"`
	LevelMode                  string `json:"levelMode,string" binding:"required"`
}

// ReqSavePlayerProgress save player progress
func ReqSavePlayerProgress(w http.ResponseWriter, r *http.Request) {
	request := savePlayerProgressRequest{}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, getErrBody(err), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	user := ctx.Value(UserCtxID).(User)
	needUpdate := false
	switch request.LevelMode {
	case "standart":
		if request.ReachedStage > user.ReachedStage01 {
			needUpdate = true
			user.ReachedStage01 = request.ReachedStage
			user.ReachedSubStage01 = request.ReachedSubStage
		} else if request.ReachedStage == user.ReachedStage01 && request.ReachedSubStage > user.ReachedSubStage01 {
			needUpdate = true
			user.ReachedSubStage01 = request.ReachedSubStage
		}
	default:
		log.Panicf("not implemented level mode %s", request.LevelMode)
	}
	// TODO
	// logic progress
	if needUpdate {
		storage.Gorm.Save(&user)
	}
	// social logic
	response := "ok"
	JSON(w, response)
}
