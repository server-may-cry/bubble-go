package application

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	newrelic "github.com/newrelic/go-agent"
	"github.com/server-may-cry/bubble-go/mynewrelic"
)

// ReqReduceTries reduce user tries by one
func ReqReduceTries(db *gorm.DB) HTTPHandlerContainer {
	handler := HTTPHandler{
		URL: "/ReqReduceTries",
	}
	handler.HTTPHandler = func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		user := r.Context().Value(userCtxID).(User)
		ts := time.Now().Unix()
		if user.RestoreTriesAt != 0 && user.RestoreTriesAt <= ts {
			if user.RemainingTries < defaultConfig.DefaultRemainingTries {
				user.RemainingTries = defaultConfig.DefaultRemainingTries
			}
		}
		user.RemainingTries--
		if user.RestoreTriesAt < 0 {
			user.RemainingTries = 0
		}
		if user.RestoreTriesAt == 0 {
			user.RestoreTriesAt = ts + int64(defaultConfig.IntervalTriesRestoration)
		}

		s := newrelic.DatastoreSegment{
			StartTime:  newrelic.StartSegmentNow(r.Context().Value(mynewrelic.Ctx).(newrelic.Transaction)),
			Product:    newrelic.DatastorePostgres,
			Collection: "user",
			Operation:  "UPDATE",
		}
		db.Save(&user)
		_ = s.End()

		JSON(w, []int8{user.RemainingTries})
	}

	return HTTPHandlerContainer{
		HTTPHandler: handler,
	}
}
