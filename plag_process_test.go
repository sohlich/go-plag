package main

import (
	"io/ioutil"
	"testing"

	log "github.com/Sirupsen/logrus"
)

//Test parallel unzip function
func TestFileUzip(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	output, err := ioutil.ReadFile("test/test.zip")
	if err != nil {
		t.Error("Cannot read file")
	}
	unzipChannel, err := unzipFile(output, map[string]bool{"c": true, "cpp": true})
	log.Error(err)

	count := 0
	for item := range unzipChannel {
		log.Debug(item.Name)
		count++
	}

	if count != 9 {
		t.Error("Function do not unzip all files")
	}
}

func TestExtensionRegex(t *testing.T) {
	filePath := "home/radek/obr.png"
	expected := "png"
	result := extensionRegex.ReplaceAllString(filePath, "")
	if expected != result {
		t.Errorf("Regex do not removed expected part: %s", result)
	}
}
