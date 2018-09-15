package sentry

import (
	"github.com/getsentry/raven-go"
)

func HandleError(err error) {
	raven.CaptureError(err, nil)
}
