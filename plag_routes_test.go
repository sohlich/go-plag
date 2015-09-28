package main

import (
	"bytes"
	// "fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestPutAssignment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	oldMongo := mongo
	mongo = &FakeDataStorage{}
	defer func(s DataStorage) { mongo = s }(oldMongo)
	router := gin.New()

	router.POST("/test", putAssignment)
	w := httptest.NewRecorder()
	json := "{\"name\": \"Test1\",\"lang\": \"java\"}"
	r, _ := http.NewRequest("POST", "/test", strings.NewReader(json))
	router.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Error("Bad response")
	}
}

func TestPutAssignmenBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	oldMongo := mongo
	mongo = &FakeDataStorage{}
	defer func(s DataStorage) { mongo = s }(oldMongo)
	router := gin.New()

	router.POST("/test", putAssignment)
	w := httptest.NewRecorder()
	json := "{\"nm\": \"Test1\",\"llng\": \"java\""
	r, _ := http.NewRequest("POST", "/test", strings.NewReader(json))
	router.ServeHTTP(w, r)
	if w.Code != 405 {
		t.Error("Bad response")
	}
}

func TestGetSupportedLangs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	oldMongo := mongo
	mongo = &FakeDataStorage{}
	defer func(s DataStorage) { mongo = s }(oldMongo)
	router := gin.New()

	router.GET("/test", getSupportedLangs)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Error("Bad response")
	}
}

func TestPutSubmission(t *testing.T) {
	gin.SetMode(gin.TestMode)
	oldMongo := mongo
	mongo = &FakeDataStorage{}
	defer func(s DataStorage) { mongo = s }(oldMongo)
	router := gin.New()

	params := map[string]string{
		"submission-meta": "{ \"owner\": \"'$p'\",   \"assignmentId\": \"55c7a86e8543eb08edca6b51\",   \"id\":\"'$p'\" }",
	}
	w := httptest.NewRecorder()
	r, _ := newfileUploadRequest("/test", params, "submission-data", "test/test.zip")

	router.POST("/test", putSubmission)
	router.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Errorf("Bad response %v", w)
	}
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
