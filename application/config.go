package application

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var userConfig struct {
}

func init() {
	configFile := "./config/user.json"
	file, err := os.Open(filepath.ToSlash(configFile))
	if err != nil {
		log.Fatal(err)
	}
	err = json.NewDecoder(file).Decode(&userConfig)
	if err != nil {
		log.Fatal(err)
	}
}
