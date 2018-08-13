package notification

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	batchLevelsGroupCount = 200
	eventTypeLevel        = 1
	eventTypeGeneral      = 2
)

// HTTPClient network client interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// VkWorker worker for app 2 user notifications
type VkWorker struct {
	ch     chan VkEvent
	config VkConfig
	token  string
	client HTTPClient

	// events to processing
	batchLevels []VkEvent
	listEvents  []VkEvent
}

// NewVkWorker create vk worker for app to user notification
func NewVkWorker(config VkConfig, client HTTPClient) *VkWorker {
	return &VkWorker{
		ch:          make(chan VkEvent),
		config:      config,
		client:      client,
		batchLevels: make([]VkEvent, 0),
		listEvents:  make([]VkEvent, 0),
	}
}

// VkConfig config for vk worker
type VkConfig struct {
	AppID           string
	Secret          string
	RequestInterval time.Duration
}

// VkEvent struct for app to user event notification in Vk
type VkEvent struct {
	ExtID int64
	Type  int
	Value int
}

type authorizationResponse struct {
	AccessToken string `json:"access_token"`
}

// SendEvent send event to worker
func (w *VkWorker) SendEvent(e VkEvent) {
	w.ch <- e
}

// LenEvents return amount of events to process
func (w *VkWorker) LenEvents() int {
	return len(w.listEvents)
}

// LenLevels return amount of levels to set
func (w *VkWorker) LenLevels() int {
	return len(w.listEvents)
}

// Initialize load access token
func (w *VkWorker) Initialize() error {
	// https://vk.com/dev/access_token
	req, err := http.NewRequest("GET", "https://oauth.vk.com/access_token", nil)
	if err != nil {
		return errors.Wrap(err, "failed to create token request")
	}

	q := req.URL.Query()
	q.Add("client_id", w.config.AppID)
	q.Add("client_secret", w.config.Secret)
	q.Add("grant_type", "client_credentials")
	req.URL.RawQuery = q.Encode()

	resp, err := w.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return errors.Wrap(err, "failed to get token")
	}
	var authResponse authorizationResponse
	if err = json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return errors.Wrap(err, "failed to decode token response")
	}
	w.token = authResponse.AccessToken
	return nil
}

// Work main work to send events
func (w *VkWorker) Work() {
	tickerLevels := time.NewTicker(w.config.RequestInterval * 2)
	tickerGeneral := time.NewTicker(w.config.RequestInterval * 2)
	for {
		select {
		case <-tickerLevels.C:
			err := w.processLevelEvents()
			if err != nil {
				log.Println("failure during processing level events", err)
			}
		case <-tickerGeneral.C:
			err := w.processGeneralEvents()
			if err != nil {
				log.Println("failure during processing general events", err)
			}
		case event, ok := <-w.ch:
			if ok {
				switch event.Type {
				case eventTypeLevel:
					w.batchLevels = append(w.batchLevels, event)
				case eventTypeGeneral:
					w.listEvents = append(w.listEvents, event)
				}
				continue
			}
			<-time.NewTicker(w.config.RequestInterval * 2).C
			if (len(w.listEvents) + len(w.batchLevels)) == 0 {
				log.Println("vk worker graceful shutdowned")
				break
			}
			log.Printf(
				"Prepare shutdown. VK worker must make:%d requests",
				len(w.listEvents)+len(w.batchLevels)/batchLevelsGroupCount,
			)
		}
	}
}

func (w *VkWorker) processLevelEvents() error {
	if len(w.batchLevels) == 0 {
		return nil
	}
	todoCount := len(w.batchLevels)
	if todoCount > batchLevelsGroupCount {
		todoCount = batchLevelsGroupCount
	}
	batchPart := w.batchLevels[:todoCount]
	userLevels := make([]string, todoCount)
	for i, e := range batchPart {
		userLevels[i] = fmt.Sprintf("%d:%d", e.ExtID, e.Value)
	}
	parameters := map[string]string{
		"levels": strings.Join(userLevels, ","),
	}
	// https://vk.com/dev/secure.setUserLevel deprecated
	err := w.sendRequest("secure.setUserLevel", parameters)
	if err != nil {
		return err
	}
	w.batchLevels = w.batchLevels[todoCount:]
	return nil
}

func (w *VkWorker) processGeneralEvents() error {
	if len(w.listEvents) == 0 {
		return nil
	}
	event := w.listEvents[0]
	parameters := map[string]string{
		"user_id":     strconv.FormatInt(event.ExtID, 10),
		"activity_id": strconv.Itoa(event.Type),
		"value":       strconv.Itoa(event.Value),
	}
	// https://vk.com/dev/secure.sendNotification
	err := w.sendRequest("secure.sendNotification", parameters)
	if err != nil {
		w.listEvents = append(w.listEvents[1:], w.listEvents[0]) // move to end of slice
		return err
	}
	w.listEvents = w.listEvents[1:]
	return nil
}

func (w *VkWorker) sendRequest(method string, parameters map[string]string) error {
	req, err := http.NewRequest("GET", "https://api.vk.com/method/"+method, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create vk request")
	}

	q := req.URL.Query()
	for k, v := range parameters {
		q.Add(k, v)
	}
	q.Add("access_token", w.token)
	q.Add("client_secret", w.config.Secret)
	// https://vk.com/dev/versions
	q.Add("v", "5.37")
	req.URL.RawQuery = q.Encode()

	resp, err := w.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return errors.Wrap(err, "failed to send vk request")
	}

	var rawResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&rawResponse)
	if err != nil {
		return errors.Wrap(err, "failed to decode vk response")
	}
	// https://vk.com/dev/errors
	val, exist := rawResponse["error"]
	if exist {
		return fmt.Errorf("vk error response %v", val)
	}
	return nil
}
