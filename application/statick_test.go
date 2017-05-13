package application

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStaticHandler(t *testing.T) {
	server := httptest.NewServer(GetRouter(true))
	defer server.Close()

	resp, err := http.Get(fmt.Sprint(server.URL, "/bubble/ShootTheBubbleDevVK.html?some=query/"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	actual, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(actual), ".swf") {
		t.Errorf("String don`t contain `.swf` '%s'\n", actual)
	}
}
