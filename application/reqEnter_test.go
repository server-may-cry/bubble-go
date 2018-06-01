package application

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestFirstGameField(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open("sqlite3", file.Name())
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&User{})

	server := httptest.NewServer(GetRouter(true, db, nil, nil))
	defer server.Close()

	data := []byte("_123_")
	jsonBytes, _ := json.Marshal(AuthRequestPart{
		AuthKey: fmt.Sprintf("%x", md5.Sum(data)),
		ExtID:   123,
		SysID:   "VK",
	})
	reader := bytes.NewReader(jsonBytes)

	resp, err := http.Post(fmt.Sprint(server.URL, "/ReqEnter"), "application/json", reader)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	var response enterResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if response.FirstGame != 1 {
		t.Fatalf("first game expected, got %+v", response)
	}

	reader.Reset(jsonBytes)
	resp, err = http.Post(fmt.Sprint(server.URL, "/ReqEnter"), "application/json", reader)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("%s raw: %s", err, string(body))
	}
	if response.FirstGame != 0 {
		t.Fatalf("not first game expected, got %+v", response)
	}
}
