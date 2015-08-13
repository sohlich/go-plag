package main

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"
	"github.com/sohlich/go-plag/parser"
)

func processSubmission(submission *Submission) error {

	sChan, err := unzipFile(submission.Content)

	if err != nil {
		return err
	}

	//Generate new submission
	//Id for group of files
	//and save the files
	submissionID := uuid.NewV1().String()
	for submissionFile := range sChan {
		submissionFile.Submission = submissionID
		tokContent, err := parser.TokenizeContent(submissionFile.Content, submission.Lang)
		if err != nil {
			return err
		}
		submissionFile.Tokens = tokContent.NGrams
		mongo.Save(submissionFile)
	}

	return nil
}

//Unzipt the files in parallel
//returns buffered channel
func unzipFile(content []byte) (<-chan *SubmissionFile, error) {
	r, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))

	if err != nil {
		log.Error(err)
		return nil, err
	}

	files := r.File
	filesLen := len(files)
	unzippedStream := make(chan *SubmissionFile, filesLen)

	var wg sync.WaitGroup
	wg.Add(filesLen)
	for _, f := range files {
		go func(file *zip.File) {
			processed, err := processFile(file)
			wg.Done()
			if err != nil {
				return
			}
			unzippedStream <- processed
		}(f)
	}

	go func() {
		wg.Wait()
		close(unzippedStream)
	}()

	return unzippedStream, nil
}

//Unzip the file and return filled
//structure for further processing
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
