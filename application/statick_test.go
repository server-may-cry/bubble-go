package application

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStaticHandler(t *testing.T) {
	staticHandler, err := NewStaticHandler("http://119226.selcdn.ru")
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/bubble/ShootTheBubbleDevVK.html?some=query/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	http.HandlerFunc(staticHandler.Serve).ServeHTTP(rr, req)

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

	files := staticHandler.storage
	if len(files) != 1 {
		t.Errorf("Expected one file in static handler, found: %+v", files)
	}
}

func TestClearStaticHandler(t *testing.T) {
	staticHandler, err := NewStaticHandler("http://119226.selcdn.ru")
	staticHandler.storage["test"] = fileCache{}
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/cache-clear", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	http.HandlerFunc(staticHandler.Clear).ServeHTTP(rr, req)

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
	files := staticHandler.storage
	if len(files) != 0 {
		t.Errorf("No files expected in static handler. found: %d", len(files))
	}
}
