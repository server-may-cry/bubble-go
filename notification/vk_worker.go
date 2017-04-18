package notification

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type authorizationResponse struct {
	AccessToken string `json:"access_token"`
}

// VkConfig config for vk worker
type VkConfig struct {
	AppID           string
	Secret          string
	RequestInterval time.Duration
}

// VkWorker worker for app 2 user notifications
type VkWorker struct {
	ch     chan VkEvent
	config VkConfig
	token  string
	client *http.Client
}

// VkEvent struct for app2user event notification in Vk
type VkEvent struct {
	ExtID string
	Type  int
	Value int
}

const batchLevelsGroupCount = 200

// NewVkWorker create vk worker for app2user notification
func NewVkWorker(config VkConfig) *VkWorker {
	worker := &VkWorker{
		ch:     make(chan VkEvent),
		config: config,
	}
	go worker.work()
	return worker
}

func (w *VkWorker) sendEvent(e VkEvent) {
	w.ch <- e
}

func (w *VkWorker) work() {
	w.client = &http.Client{}

	req, err := http.NewRequest("GET", "https://oauth.vk.com/access_token", nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	q := req.URL.Query()
	q.Add("client_id", w.config.AppID)
	q.Add("client_secret", w.config.Secret)
	q.Add("grant_type", "client_credentials")
	req.URL.RawQuery = q.Encode()

	resp, err := w.client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var authResponse authorizationResponse
	err = decoder.Decode(&authResponse)
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.token = authResponse.AccessToken

	ticker := time.NewTicker(w.config.RequestInterval)
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
				w.sendRequest("secure.sendNotification", parameters)
				continue
			}
			if len(batchLevels) > 0 {
				// TODO make batch request
				// sendRequest("secure.setUserLevel", parameters)
				continue
			}
		case event, ok := <-w.ch:
			if ok {
				switch event.Type {
				case 1:
					batchLevels = append(batchLevels, event)
				case 2:
					listEvents = append(listEvents, event)
				}
			} else {
				w.ch = nil
			}
		}
		if w.ch == nil { // closed channel on shutdown. TODO make shutdown
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

func (w *VkWorker) sendRequest(method string, parameters map[string]string) {
	req, err := http.NewRequest("GET", fmt.Sprint("https://api.vk.com/method/", method), nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	q := req.URL.Query()
	for k, v := range parameters {
		q.Add(k, v)
	}
	q.Add("access_token", w.token)
	q.Add("client_secret", w.config.Secret)
	req.URL.RawQuery = q.Encode()

	resp, err := w.client.Do(req)
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
