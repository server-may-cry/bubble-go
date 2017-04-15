package notification

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type authorizationResponse struct {
	AccessToken string `json:"access_token"`
}

var authResponse authorizationResponse
var client *http.Client

const requestInterval = time.Millisecond * 300
const batchLevelsGroupCount = 200

// VkWorkerInit run vk app2user notification
func VkWorkerInit(ch <-chan VkEvent) {
	client = &http.Client{}

	req, err := http.NewRequest("GET", "https://oauth.vk.com/access_token", nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	q := req.URL.Query()
	q.Add("client_id", os.Getenv("VK_APP_ID"))
	q.Add("client_secret", os.Getenv("VK_SECRET"))
	q.Add("v", "5.37")
	q.Add("grant_type", "client_credentials")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&authResponse)
	if err != nil {
		log.Println(err.Error())
		return
	}

	ticker := time.NewTicker(requestInterval)
	batchLevels := make([]VkEvent, 0)
	listEvents := make([]VkEvent, 0)
	for {
		select {
		case <-ticker.C:
			var parameters map[string]string
			if len(listEvents) > 0 {
				event := listEvents[0]
				parameters["user_id"] = event.ExtID
				parameters["activity_id"] = strconv.Itoa(event.Type)
				parameters["value"] = strconv.Itoa(event.Value)
				sendRequest("secure.sendNotification", parameters)
				continue
			}
			if len(batchLevels) > 0 {
				// TODO make batch request
				// sendRequest("secure.setUserLevel", parameters)
				continue
			}
		case event, ok := <-ch:
			if ok {
				switch event.Type {
				case 1:
					batchLevels = append(batchLevels, event)
				case 2:
					listEvents = append(listEvents, event)
				}
			} else {
				ch = nil
			}
		}
		if ch == nil { // closed channel on shutdown. TODO make shutdown
			log.Printf(
				"Prepare shutdown. VK worker must make:%d requests",
				len(listEvents)+len(batchLevels)/batchLevelsGroupCount,
			)
			if (len(listEvents) + len(batchLevels)/batchLevelsGroupCount) == 0 {
				break
			}
		}
	}
}

func sendRequest(method string, parameters map[string]string) {
	req, err := http.NewRequest("GET", fmt.Sprint("https://api.vk.com/method/", method), nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	q := req.URL.Query()
	for k, v := range parameters {
		q.Add(k, v)
	}
	q.Add("access_token", authResponse.AccessToken)
	q.Add("client_secret", os.Getenv("VK_SECRET"))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	var rawResponse map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&rawResponse)
	if err != nil {
		log.Println(err.Error())
		return
	}
	val, exist := rawResponse["error"]
	if exist {
		log.Println(val)
		return
	}
}
