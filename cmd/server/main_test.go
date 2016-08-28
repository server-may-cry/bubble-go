package main_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	main "github.com/server-may-cry/bubble-go/cmd/server"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestIndexRoute(c *C) {
	testRouter := main.GetEngine()
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Println(err)
	}
	response := httptest.NewRecorder()
	testRouter.ServeHTTP(response, request)
	c.Assert(response.Code, Equals, http.StatusOK)
}
