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

type reduceCreditsCompleteRequest struct {
	AuthRequestPart
	reduceCreditsRequest
}

func TestReduceCredits(t *testing.T) {
	server := httptest.NewServer(GetRouter(true))
	defer server.Close()

	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open("sqlite3", file.Name())
	db.AutoMigrate(&User{})
	user := User{
		SysID:   1,
		ExtID:   123,
		Credits: 900,
	}
	db.Create(&user)
	Gorm = db

	data := []byte("_123_")
	authKey := fmt.Sprintf("%x", md5.Sum(data))
	jsonBytes, _ := json.Marshal(reduceCreditsCompleteRequest{
		AuthRequestPart: AuthRequestPart{
			ExtID:   123,
			SysID:   "VK",
			AuthKey: authKey,
		},
		reduceCreditsRequest: reduceCreditsRequest{
			Amount: 150,
		},
	})

	reader := bytes.NewReader(jsonBytes)
	resp, err := http.Post(fmt.Sprint(server.URL, "/ReqReduceCredits"), "application/json", reader)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	decoder := json.NewDecoder(resp.Body)
	var response reduceCreditsResponse
	err = decoder.Decode(&response)
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
