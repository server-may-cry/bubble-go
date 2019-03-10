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

func TestFirstGameField(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open("sqlite3", file.Name())
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&models.User{})

	data := []byte("_123_")
	jsonBytes, _ := json.Marshal(AuthRequestPart{
		AuthKey: fmt.Sprintf("%x", md5.Sum(data)),
		ExtID:   123,
		SysID:   "VK",
	})

	reader := bytes.NewReader(jsonBytes)
	handlerContainer := ReqEnter(db)
	resp, err := testAppHandler(handlerContainer.HTTPHandler, nil, "/ReqEnter", reader)
	if err != nil {
		t.Fatal(err)
	}
	var response enterResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if response.FirstGame != 1 {
		t.Fatalf("first game expected, got %+v", response)
	}

	jsonBytesResponse, err := json.Marshal(response.SubStagesRecordStats01)
	if err != nil {
		t.Fatal(err)
	}
	user := models.User{
		SysID:            1,
		ExtID:            123,
		Credits:          900,
		ProgressStandart: string(jsonBytesResponse),
	}
	db.Create(&user)

	reader.Reset(jsonBytes)
	resp, err = testAppHandler(handlerContainer.HTTPHandler, &user, "/ReqEnter", reader)
	if err != nil {
		t.Fatal(err)
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
