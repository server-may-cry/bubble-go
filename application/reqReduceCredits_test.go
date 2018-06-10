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
)

type reduceCreditsCompleteRequest struct {
	AuthRequestPart
	reduceCreditsRequest
}

func TestReduceCredits(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open("sqlite3", file.Name())
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&User{})
	user := User{
		SysID:   1,
		ExtID:   123,
		Credits: 900,
	}
	db.Create(&user)

	data := []byte("_123_")
	jsonBytes, _ := json.Marshal(reduceCreditsCompleteRequest{
		AuthRequestPart: AuthRequestPart{
			ExtID:   123,
			SysID:   "VK",
			AuthKey: fmt.Sprintf("%x", md5.Sum(data)),
		},
		reduceCreditsRequest: reduceCreditsRequest{
			Amount: 150,
		},
	})

	reader := bytes.NewReader(jsonBytes)
	handlerContainer := ReqReduceCredits(db)
	resp, err := testAppHandler(handlerContainer.HTTPHandler, &user, "/ReqReduceCredits", reader)
	if err != nil {
		t.Fatal(err)
	}
	var response reduceCreditsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if response.Credits != 750 {
		t.Fatalf("expected 750 credits, got %+v", response)
	}

	db.First(&user, user.ID)
	if user.Credits != 750 {
		t.Fatalf("expected 750 credits in db, got %+v", user)
	}
}
