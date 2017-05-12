package application

import (
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

const (
	cdnroot = "http://119226.selcdn.ru"
)

var tmpDirName string

func init() {
	if tmpDirName == "" {
		dir, err := ioutil.TempDir("", "bubble_cache_")
		if err != nil {
			panic(err)
		}
		tmpDirName = dir
	}
}

// ServeStatick load (if not exist) static from file server (crutch for spend less money and not store static files in repo)
func ServeStatick(w http.ResponseWriter, r *http.Request) {
	fullFilePath := filepath.ToSlash(tmpDirName + r.RequestURI)
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
		resp, err := http.Get(cdnroot + r.RequestURI)
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

// ClearStatickCache remove statick files
func ClearStatickCache(w http.ResponseWriter, r *http.Request) {
	err := os.RemoveAll(tmpDirName)
	if err != nil {
		panic(err)
	}
	JSON(w, "done")
}
