package main

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"sync"

	log "github.com/Sirupsen/logrus"
)

func processSubmission(submission *Submission) {

	sChan, err := unzipFile(submission.Content)

	if err != nil {
		return
	}

	for submissionFile := range sChan {
		submissionFile.Submission = submission.AssignmentID
		mongo.Save(submissionFile)
	}
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
