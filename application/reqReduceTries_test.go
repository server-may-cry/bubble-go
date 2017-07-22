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

func TestReduceTries(t *testing.T) {
	server := httptest.NewServer(GetRouter(true))
	defer server.Close()

	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open("sqlite3", file.Name())
	db.AutoMigrate(&User{})
	user := User{
		SysID:          1,
		ExtID:          123,
		RemainingTries: 8,
	}
	db.Create(&user)
	Gorm = db

	data := []byte("_123_")
	authKey := fmt.Sprintf("%x", md5.Sum(data))
	jsonBytes, _ := json.Marshal(AuthRequestPart{
		ExtID:   123,
		SysID:   "VK",
		AuthKey: authKey,
	})

	reader := bytes.NewReader(jsonBytes)
	resp, err := http.Post(fmt.Sprint(server.URL, "/ReqReduceTries"), "application/json", reader)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	var response []int8
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if response[0] != 7 {
		t.Fatalf("expected 7 remaining tries, got %+v", response)
	}

	db.First(&user, user.ID)
	if user.RemainingTries != 7 {
		t.Fatalf("expected 7 remaining tries, got %+v", user)
	}
}
