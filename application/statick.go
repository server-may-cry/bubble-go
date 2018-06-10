package application

import (
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// StatickHandler handler to resolve work with static content from CDN
type StatickHandler struct {
	cdnroot    string
	tmpDirName string
}

// http://119226.selcdn.ru/bubble/ShootTheBubbleDevVK.html
// http://bubble-srv-dev.herokuapp.com/bubble/ShootTheBubbleDevVK.html

// NewStatickHandler create static handler
func NewStatickHandler(cdnroot string) (*StatickHandler, error) {
	tmpDirName, err := ioutil.TempDir("", "bubble_cache_")
	if err != nil {
		return nil, errors.Wrap(err, "can't create tmp dir for static")
	}
	return &StatickHandler{
		cdnroot:    cdnroot,
		tmpDirName: tmpDirName,
	}, nil
}

// Serve resolve content from CDN
func (sh StatickHandler) Serve(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}
	filePath := r.URL.Path
	fullFilePath := filepath.ToSlash(sh.tmpDirName + filePath)
	if _, err := os.Stat(fullFilePath); os.IsNotExist(err) {
		dirToStoreFile := filepath.Dir(fullFilePath)
		if _, err = os.Stat(dirToStoreFile); os.IsNotExist(err) {
			err = os.MkdirAll(dirToStoreFile, 0777)
			if err != nil {
				panic(err)
			}
		}
		out, err := os.Create(fullFilePath)
		if err != nil {
			panic(err)
		}
		defer out.Close()
		resp, err := http.Get(sh.cdnroot + filePath)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			panic(err)
		}
	}

	dat, err := ioutil.ReadFile(fullFilePath)
	if err != nil {
		panic(err)
	}
	ext := filepath.Ext(fullFilePath)
	w.Header().Set("Content-Type", mime.TypeByExtension(ext))
	_, err = w.Write(dat)
	if err != nil {
		panic(err)
	}
}

// Clear remove statick files
func (sh StatickHandler) Clear(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}
	err := os.RemoveAll(sh.tmpDirName)
	if err != nil {
		panic(err)
	}
	JSON(w, "done")
}
