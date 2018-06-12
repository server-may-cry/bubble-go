package application

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

func testAppHandler(
	hc HTTPHandler,
	user *User,
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
