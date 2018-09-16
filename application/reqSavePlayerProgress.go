package application

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	newrelic "github.com/newrelic/go-agent"
	"github.com/server-may-cry/bubble-go/models"
	"github.com/server-may-cry/bubble-go/mynewrelic"
	"github.com/server-may-cry/bubble-go/notification"
	"github.com/server-may-cry/bubble-go/platforms"
)

var completeLastLevelOnIslandEventMap []int

func init() {
	completeLastLevelOnIslandEventMap = []int{
		4,
		5,
		6,
		7,
		8,
		9,
	}
}

/*
{
    ...
    "reachedSubStage":"1", // номер точки на острове до которой игрок дошел за все время игры
    "currentStage":"0", // номер острова на котором игрок прошел точку в текущей сессии
    "reachedStage":"0", // номер острова до которого игрок дошел за все время игры
    "completeSubStage":"1", // номер точки на острове которую игрок прошел в текущей сессии
    "completeSubStageRecordStat":"2", // количество звезд набранных на пройденной точке
    "levelMode":"standart", // режим игры, может принимать значения "standart" (01) и "arcade" (02)
    Если режим игры "standart", то перезаписываются значения reachedStage01 reachedSubStage01,
    которые приходят в ReqEnter`e, если же "arcade" то reachedStage02 и reachedSubStage02
}
*/
type savePlayerProgressRequest struct {
	ReachedSubStage            int8   `json:"reachedSubStage,string"`
	CurrentStage               int8   `json:"currentStage,string"` // island number
	ReachedStage               int8   `json:"reachedStage,string"`
	CompleteSubStage           int8   `json:"completeSubStage,string"`           // level number on island
	CompleteSubStageRecordStat int8   `json:"completeSubStageRecordStat,string"` // starCount
	LevelMode                  string `json:"levelMode"`                         // standart
}

// ReqSavePlayerProgress save player progress
func ReqSavePlayerProgress(
	vkWorker *notification.VkWorker,
	db *gorm.DB,
) HTTPHandlerContainer {
	handler := HTTPHandler{
		URL: "/ReqSavePlayerProgress",
	}
	handler.HTTPHandler = func(w http.ResponseWriter, r *http.Request) {
		request := savePlayerProgressRequest{}
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user := r.Context().Value(userCtxID).(models.User)
		var needUpdate bool
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
			progress := user.GetProgresStandart()
			if int(request.CurrentStage) > len(progress) {
				log.Panicf("cheating? can`t save progress %+v", request)
			}
			if request.CompleteSubStageRecordStat > progress[request.CurrentStage][request.CompleteSubStage] {
				needUpdate = true
				progress[request.CurrentStage][request.CompleteSubStage] = request.CompleteSubStageRecordStat
				user.SetProgresStandart(progress)
			}
		default:
			log.Panicf("not implemented level mode %s", request.LevelMode)
		}
		if needUpdate {
			s := newrelic.DatastoreSegment{
				StartTime:  newrelic.StartSegmentNow(r.Context().Value(mynewrelic.Ctx).(newrelic.Transaction)),
				Product:    newrelic.DatastorePostgres,
				Collection: "user",
				Operation:  "UPDATE",
			}
			db.Save(&user)
			_ = s.End()
		}
		if user.SysID == platforms.VK {
			vkSocialLogic(vkWorker, request, user)
		}
		JSON(w, "ok")
	}

	return HTTPHandlerContainer{
		HTTPHandler: handler,
	}
}

func vkSocialLogic(vkWorker *notification.VkWorker, request savePlayerProgressRequest, user models.User) {
	if request.CompleteSubStageRecordStat > 0 {
		// not failed level
		levelOrder := 0
		if request.CurrentStage > 0 {
			levelOrder = int(request.CurrentStage)*14 - 6
		}
		levelOrder += int(request.CompleteSubStage) + 1
		prevReachedLevelOrder := 0
		if user.ReachedStage01 > 0 {
			prevReachedLevelOrder = int(user.ReachedStage01)*14 - 6
		}
		prevReachedLevelOrder += int(request.ReachedSubStage) + 1
		if levelOrder > prevReachedLevelOrder {
			vkWorker.SendEvent(notification.VkEvent{
				ExtID: user.ExtID,
				Type:  1,
				Value: levelOrder,
			})
		}
	}
	if request.CompleteSubStage == 14 || (request.CompleteSubStage == 8 && request.CurrentStage == 0) {
		// open new island event
		islandOrder := request.CurrentStage // complete last mission on island
		eventID := completeLastLevelOnIslandEventMap[islandOrder]
		if eventID != 0 {
			if request.CurrentStage > user.ReachedStage01 {
				vkWorker.SendEvent(notification.VkEvent{
					ExtID: user.ExtID,
					Type:  2,
					Value: eventID,
				})
			}
		}
	}
}
