package application

import (
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	newrelic "github.com/newrelic/go-agent"
	"github.com/server-may-cry/bubble-go/mynewrelic"
)

type fileCache struct {
	content  []byte
	mimeType string
}

// StaticHandler handler to resolve work with static content from CDN
type StaticHandler struct {
	cdnroot   string
	httpClien *http.Client
	mutex     sync.RWMutex
	storage   map[string]fileCache // For now total size of all static files ~65 MiB
}

// https://119226.selcdn.ru/bubble/ShootTheBubbleDevVK.html
// https://bubble-srv-dev.herokuapp.com/bubble/ShootTheBubbleDevVK.html

// NewStaticHandler create static handler
func NewStaticHandler(cdnroot string) (*StaticHandler, error) {
	return &StaticHandler{
		cdnroot: cdnroot,
		httpClien: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns: 5,
			},
			Timeout: 5 * time.Second,
		},
		storage: make(map[string]fileCache),
	}, nil
}

// Serve resolve content from CDN
func (sh *StaticHandler) Serve(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}
	err := r.Context().Value(mynewrelic.Ctx).(newrelic.Transaction).SetName("/static_serve")
	if err != nil {
		panic(err)
	}
	filePath := r.URL.Path
	fileForResponse, ok := sh.storage[filePath]
	if !ok {
		resp, err := sh.httpClien.Get(sh.cdnroot + filePath)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			panic(filePath)
		}
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		ext := filepath.Ext(filePath)
		fileForResponse = fileCache{
			content:  content,
			mimeType: mime.TypeByExtension(ext),
		}
		sh.mutex.Lock()
		sh.storage[filePath] = fileForResponse
		sh.mutex.Unlock()
	}

	w.Header().Set("Content-Type", fileForResponse.mimeType)
	_, err = w.Write(fileForResponse.content)
	if err != nil {
		panic(err)
	}
}

// Clear remove statick files
func (sh *StaticHandler) Clear(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}
	sh.mutex.Lock()
	defer sh.mutex.Unlock()
	sh.storage = make(map[string]fileCache)
	JSON(w, "done")
}
