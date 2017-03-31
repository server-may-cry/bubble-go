package controllers

import (
	"io"
	"io/ioutil"
	"log"
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

// ServeStatick load (if not exist) static from file server (crutch for spend less money and not store static files in repo)
func ServeStatick(c *gin.Context) {
	filePath := c.Param("filePath")
	if tmpDirName == "" {
		dir, err := ioutil.TempDir("", "bubble_cache_")
		if err != nil {
			log.Fatal(err)
		}
		tmpDirName = dir
	}
	fullFilePath := filepath.ToSlash(tmpDirName + "/bubble" + filePath)
	if _, err := os.Stat(fullFilePath); os.IsNotExist(err) {
		dirToStoreFile := filepath.Dir(fullFilePath)
		if _, err = os.Stat(dirToStoreFile); os.IsNotExist(err) {
			err = os.MkdirAll(dirToStoreFile, 0777)
			if err != nil {
				log.Fatal(err)
			}
		}
		out, err := os.Create(fullFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		resp, err := http.Get(cdnroot + filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Fatal(err)
		}
	}

	dat, err := ioutil.ReadFile(fullFilePath)
	if err != nil {
		log.Fatal(err)
	}
	ext := filepath.Ext(fullFilePath)
	c.Writer.Header().Set("Content-Type", mime.TypeByExtension(ext))
	_, err = c.Writer.WriteString(string(dat))
	if err != nil {
		log.Fatal(err)
	}
}

// ClearStatickCache remove statick files
func ClearStatickCache(c *gin.Context) {
	err := os.RemoveAll(tmpDirName)
	if err != nil {
		log.Fatal(err)
	}
	c.String(http.StatusOK, "done")
}
