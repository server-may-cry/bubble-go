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

type usersProgressCompleteRequest struct {
	AuthRequestPart
	usersProgressRequest
}

func TestUsersProgress(t *testing.T) {
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
	jsonBytes, _ := json.Marshal(usersProgressCompleteRequest{
		AuthRequestPart{
			ExtID:   123,
			SysID:   "VK",
			AuthKey: fmt.Sprintf("%x", md5.Sum(data)),
		},
		usersProgressRequest{
			SocIDs: []int64{123},
		},
	})

	reader := bytes.NewReader(jsonBytes)
	handlerContainer := ReqUsersProgress(db)
	resp, err := testAppHandler(handlerContainer.HTTPHandler, &user, "/ReqUsersProgress", reader)
	if err != nil {
		t.Fatal(err)
	}
	var response usersProgressResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if len(response.UsersProgress) != 1 || response.UsersProgress[0].SocID != 123 {
		t.Fatalf("expected users progress, got %+v", response)
	}
}
