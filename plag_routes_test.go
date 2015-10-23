package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sohlich/go-plag/parser"
)

//Sets gin to testing mode and presets the
//data storage to FakeDataStorage
//to emulate the database.
func preTest() (*gin.Engine, DataStorage) {
	gin.SetMode(gin.TestMode)
	oldMongo := mongo
	mongo = &FakeDataStorage{}
	router := gin.New()
	return router, oldMongo
}

func TestPutAssignment(t *testing.T) {
	router, oldMongo := preTest()
	defer func(s DataStorage) { mongo = s }(oldMongo)

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
	router, oldMongo := preTest()
	defer func(s DataStorage) { mongo = s }(oldMongo)

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
	router, oldMongo := preTest()
	defer func(s DataStorage) { mongo = s }(oldMongo)

	router.GET("/test", getSupportedLangs)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Error("Bad response")
	}
}

func TestPutSubmission(t *testing.T) {
	parser.LoadPlugins("plugin")
	router, oldMongo := preTest()
	defer func(s DataStorage) { mongo = s }(oldMongo)

	params := map[string]string{
		"submission-meta": "{ \"owner\": \"radek\",   \"assignmentId\": \"55c7a86e8543eb08edca6b51\",   \"id\":\"25\" }",
	}
	w := httptest.NewRecorder()
	r, _ := newfileUploadRequest("/test", params, "submission-data", "test/test.zip")

	router.POST("/test", putSubmission)
	router.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Errorf("Bad response %v", w)
	}
	time.Sleep(100 * time.Millisecond)
}

func TestPutSubmissionInvalidAssignment(t *testing.T) {
	parser.LoadPlugins("plugin")
	router, oldMongo := preTest()
	defer func(s DataStorage) { mongo = s }(oldMongo)

	params := map[string]string{
		"submission-meta": "{ \"owner\": \"'$p'\",   \"assignmentId\": \"\",   \"id\":\"'$p'\" }",
	}
	w := httptest.NewRecorder()
	r, _ := newfileUploadRequest("/test", params, "submission-data", "test/test.zip")

	router.POST("/test", putSubmission)
	router.ServeHTTP(w, r)
	if w.Code != 405 {
		t.Errorf("Bad response expected %s got %v", 405, w)
	}

	time.Sleep(100 * time.Millisecond)
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
