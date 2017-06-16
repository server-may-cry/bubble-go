package application

import (
	"io/ioutil"
	"testing"
)

var testsDataset = []struct {
	content string
	err     string
}{
	{`
{
  "default_remaining_tries": 5,
  "default_credits": {
    "vk": 1000,
    "ok": 3000
  },
  "interval_tries_restoration": 1800,
  "friends_bonus_credits_multiplier": 40,
  "init_progress": [
      [-1,-1,-1,-1,-1,-1,-1,-1],
      [-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1],
      [-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1],
      [-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1],
      [-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1],
      [-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1],
      [-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1]
  ]
}`, ""},
	{``, "can`t decode user.json: EOF"},
}

func TestConfigInit(t *testing.T) {
	for _, tt := range testsDataset {
		file, err := ioutil.TempFile("", "")
		if err != nil {
			panic(err)
		}
		_, err = file.WriteString(tt.content)
		if err != nil {
			panic(err)
		}
		err = ConfigInit(file.Name())
		if tt.err != "" && err != nil && err.Error() != tt.err {
			t.Errorf("data %s config error expexted '%s', got '%s'", tt.content, tt.err, err)
		} else if tt.err == "" && err != nil {
			t.Errorf("data %s config error expexted '%s', got '%s'", tt.content, tt.err, err)
		}
	}
}
