package application

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/server-may-cry/bubble-go/models"
)

func TestReduceTries(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open("sqlite3", file.Name())
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&models.User{})
	user := models.User{
		SysID:          1,
		ExtID:          123,
		RemainingTries: 8,
	}
	db.Create(&user)

	data := []byte("_123_")
	jsonBytes, _ := json.Marshal(AuthRequestPart{
		ExtID:   123,
		SysID:   "VK",
		AuthKey: fmt.Sprintf("%x", md5.Sum(data)),
	})

	reader := bytes.NewReader(jsonBytes)
	handlerContainer := ReqReduceTries(db)
	resp, err := testAppHandler(handlerContainer.HTTPHandler, &user, "/ReqReduceTries", reader)
	if err != nil {
		t.Fatal(err)
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
