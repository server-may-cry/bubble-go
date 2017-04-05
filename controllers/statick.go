package controllers

import (
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/gin-gonic/gin.v1"
)

const (
	cdnroot = "http://119226.selcdn.ru/bubble"
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
func ServeStatick(c *gin.Context) {
	filePath := c.Param("filePath")
	fullFilePath := filepath.ToSlash(tmpDirName + "/bubble" + filePath)
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
		resp, err := http.Get(cdnroot + filePath)
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
	c.Writer.Header().Set("Content-Type", mime.TypeByExtension(ext))
	_, err = c.Writer.WriteString(string(dat))
	if err != nil {
		panic(err)
	}
}

// ClearStatickCache remove statick files
func ClearStatickCache(c *gin.Context) {
	err := os.RemoveAll(tmpDirName)
	if err != nil {
		panic(err)
	}
	c.String(http.StatusOK, "done")
}
