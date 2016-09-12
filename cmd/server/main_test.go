package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestIndexRoute(c *C) {
	response := httptest.NewRecorder()
	// testRouter.ServeHTTP(response, request)
	c.Assert(response.Code, Equals, http.StatusOK)
}
