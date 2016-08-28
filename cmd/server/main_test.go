package main_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	main "github.com/server-may-cry/bubble-go/cmd/server"
)

func TestIndexRoute(t *testing.T) {
	testRouter := main.GetEngine()
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		fmt.Println(err)
	}
	response := httptest.NewRecorder()
	testRouter.ServeHTTP(response, request)
	//assert.Equal(t, response.Code, 200)

	log.Print(request, response)
}
