package application

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	newrelic "github.com/newrelic/go-agent"
	"github.com/server-may-cry/bubble-go/models"
	"github.com/server-may-cry/bubble-go/mynewrelic"
)

func testAppHandler(
	hc HTTPHandler,
	user *models.User,
	urlExpected string,
	body io.Reader,
) (*httptest.ResponseRecorder, error) {
	if urlExpected != hc.URL {
		return nil, fmt.Errorf("Wrong url. expected '%s', got '%s'", urlExpected, hc.URL)
	}

	req, err := http.NewRequest("POST", urlExpected, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	config := newrelic.NewConfig("bubble-go", "1234567890123456789012345678901234567890")
	app, err := newrelic.NewApplication(config)
	if err != nil {
		return nil, err
	}
	ctx := context.WithValue(req.Context(), mynewrelic.Ctx, app.StartTransaction("test", nil, nil))
	req = req.WithContext(ctx)

	if user != nil {
		ctx := req.Context()
		ctx = context.WithValue(ctx, userCtxID, *user)
		req = req.WithContext(ctx)
	}

	rr := httptest.NewRecorder()

	hc.HTTPHandler.ServeHTTP(rr, req)

	if rr.Code != 200 {
		return nil, fmt.Errorf("Bad status code returned '%d'", rr.Code)
	}

	return rr, nil
}
