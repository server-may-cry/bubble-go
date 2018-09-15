package sentry

import (
	"github.com/getsentry/raven-go"
)

// HandleError send error into sentry
func HandleError(err error) {
	raven.CaptureError(err, nil)
}
