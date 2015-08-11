package main

import (
	//"encoding/json"
	"io/ioutil"
	//"strings"
	"testing"

	//log "github.com/Sirupsen/logrus"
)

func TestFileUzip(t *testing.T) {
	output, err := ioutil.ReadFile("test/test.zip")
	if err != nil {
		t.Error("Cannot read file")
	}
	unzipFile(output)
}
