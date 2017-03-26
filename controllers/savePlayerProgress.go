package controllers

import (
	"net/http"

	"github.com/server-may-cry/bubble-go/models"
	"github.com/server-may-cry/bubble-go/storage"

	"gopkg.in/gin-gonic/gin.v1"
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
	ReachedSubStage            uint8  `json:"reachedSubStage,string" binding:"required"`
	CurrentStage               uint8  `json:"currentStage,string" binding:"required"`
	ReachedStage               uint8  `json:"reachedStage,string" binding:"required"`
	CompleteSubStage           uint8  `json:"completeSubStage,string" binding:"required"`
	CompleteSubStageRecordStat uint8  `json:"completeSubStageRecordStat,string" binding:"required"`
	LevelMode                  string `json:"levelMode,string" binding:"required"`
}

// ReqSavePlayerProgress save player progress
func ReqSavePlayerProgress(c *gin.Context) {
	request := savePlayerProgressRequest{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, getErrBody(err))
		return
	}
	user := c.MustGet("user").(models.User)
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
		panic("not implemented level mode")
	}
	// TODO
	// logic progress
	if needUpdate {
		storage.Gorm.Save(&user)
	}
	// social logic
	response := "ok"
	c.JSON(http.StatusOK, response)
}
