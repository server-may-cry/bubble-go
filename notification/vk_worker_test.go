package notification

import (
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestErrorOnInvalidResponse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://api.vk.com/method/mockMethod?access_token=&client_secret=secret&foo=bar&v=5.37",
		httpmock.NewStringResponder(200, `{"error": "mock error"}`),
	)
	worker := NewVkWorker(VkConfig{
		AppID:           "123",
		Secret:          "secret",
		RequestInterval: time.Second,
	}, &http.Client{})
	err := worker.sendRequest("mockMethod", map[string]string{"foo": "bar"})
	if err == nil || err.Error() != "vk error response mock error" {
		t.Error(err)
	}
}

func TestInitialization(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://oauth.vk.com/access_token?client_id=123&client_secret=secret&grant_type=client_credentials",
		httpmock.NewStringResponder(200, `{"access_token": "mockTocken"}`),
	)
	worker := NewVkWorker(VkConfig{
		AppID:           "123",
		Secret:          "secret",
		RequestInterval: time.Second,
	}, &http.Client{})
	err := worker.Initialize()
	if err != nil {
		t.Error(err)
	}
	if worker.token != "mockTocken" {
		t.Errorf("Expected mock tocken in worker, got %s", worker.token)
	}
}

func TestSendLevelEvent(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://api.vk.com/method/secure.setUserLevel?access_token=&client_secret=secret&levels=123%3A1%2C1234%3A4&v=5.37",
		httpmock.NewStringResponder(200, `{}`),
	)
	worker := NewVkWorker(VkConfig{
		AppID:           "123",
		Secret:          "secret",
		RequestInterval: time.Second,
	}, &http.Client{})
	worker.batchLevels = append(worker.batchLevels, VkEvent{
		ExtID: 123,
		Type:  eventTypeLevel,
		Value: 1,
	})
	worker.batchLevels = append(worker.batchLevels, VkEvent{
		ExtID: 1234,
		Type:  eventTypeLevel,
		Value: 4,
	})
	err := worker.processLevelEvents()
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkSendLevelEvent(b *testing.B) {
	worker := NewVkWorker(VkConfig{
		AppID:           "123",
		Secret:          "secret",
		RequestInterval: time.Second,
	}, &http.Client{})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		httpmock.Activate()
		httpmock.RegisterResponder(
			"GET",
			"https://api.vk.com/method/secure.setUserLevel?access_token=&client_secret=secret&levels=123%3A1%2C1234%3A4&v=5.37",
			httpmock.NewStringResponder(200, `{}`),
		)
		worker.batchLevels = append(worker.batchLevels, VkEvent{
			ExtID: 123,
			Type:  eventTypeLevel,
			Value: 1,
		})
		worker.batchLevels = append(worker.batchLevels, VkEvent{
			ExtID: 1234,
			Type:  eventTypeLevel,
			Value: 4,
		})
		err := worker.processLevelEvents()
		if err != nil {
			b.Fatal(err)
		}
		httpmock.DeactivateAndReset()
	}
}

func TestSendGeneralEvent(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://api.vk.com/method/"+
			"secure.sendNotification?access_token=&activity_id=2&client_secret=secret&user_id=123&v=5.37&value=1",
		httpmock.NewStringResponder(200, `{}`),
	)
	worker := NewVkWorker(VkConfig{
		AppID:           "123",
		Secret:          "secret",
		RequestInterval: time.Second,
	}, &http.Client{})
	worker.listEvents = append(worker.listEvents, VkEvent{
		ExtID: 123,
		Type:  eventTypeGeneral,
		Value: 1,
	})
	err := worker.processGeneralEvents()
	if err != nil {
		t.Error(err)
	}
}
