package application

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStaticHandler(t *testing.T) {
	statickHandler, err := NewStatickHandler("http://119226.selcdn.ru")
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/bubble/ShootTheBubbleDevVK.html?some=query/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(statickHandler.Serve)

	handler.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Received non-200 response: %d\n", rr.Code)
	}

	actual, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(actual), ".swf") {
		t.Errorf("String don`t contain `.swf` '%s'\n", actual)
	}

	files, _ := ioutil.ReadDir(statickHandler.tmpDirName)
	if len(files) < 1 {
		t.Errorf("No files in cache folder found. %d", len(files))
	}
}

func TestClearStaticHandler(t *testing.T) {
	statickHandler, err := NewStatickHandler("http://119226.selcdn.ru")
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/cache-clear", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(statickHandler.Clear)

	handler.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Received non-200 response: %d\n", rr.Code)
	}

	actual, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(actual) != "\"done\"\n" {
		t.Errorf("Invalid response (%s), expected (\"done\"\n)\n", actual)
	}
	files, _ := ioutil.ReadDir(statickHandler.tmpDirName)
	if len(files) > 0 {
		t.Errorf("no files expected in tmp directory. found: %d", len(files))
	}
}
