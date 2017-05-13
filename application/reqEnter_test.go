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

func TestRegistration(t *testing.T) {
	server := httptest.NewServer(GetRouter(true))
	defer server.Close()

	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open("sqlite3", file.Name())
	db.AutoMigrate(&User{})
	Gorm = db

	data := []byte("_123_")
	authKey := fmt.Sprintf("%x", md5.Sum(data))
	jsonBytes, _ := json.Marshal(AuthRequestPart{
		AuthKey: authKey,
		ExtID:   123,
		SysID:   "VK",
	})
	reader := bytes.NewReader(jsonBytes)

	resp, err := http.Post(fmt.Sprint(server.URL, "/ReqEnter"), "application/json", reader)
	if err != nil {
		t.Fatal(err)
	}
	decoder := json.NewDecoder(resp.Body)
	var response enterResponse
	err = decoder.Decode(response)
	if err != nil {
		t.Fatal(err)
	}
	t.Fatal(response)
}
