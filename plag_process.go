package main

import (
	"archive/zip"
	"bytes"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
)

func unzipFile(content []byte) {
	r, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))

	if err != nil {
		log.Error(err)
		return
	}

	zipFileStream := make(chan *SubmissionFile, 10)

	files := r.File

	for _, f := range files {
		go func(file *zip.File) {
			processed, err := processFile(file)
			if err != nil {
				return
			}
			zipFileStream <- processed
		}(f)
	}

	for i := 0; i < len(files); i++ {
		val := <-zipFileStream
		log.Info(val.Name)
	}
	close(zipFileStream)
}

func processFile(file *zip.File) (*SubmissionFile, error) {
	submissionFile := &SubmissionFile{Name: file.Name}
	rc, rcError := file.Open()
	if rcError != nil {
		log.Error(rcError)
		return nil, rcError
	}
	content, rError := ioutil.ReadAll(rc)
	if rError != nil {
		log.Error(rError)
		return nil, rError
	}
	submissionFile.Content = string(content)
	return submissionFile, nil
}
