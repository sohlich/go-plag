package main

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"regexp"
	"sync"

	"github.com/satori/go.uuid"
	"github.com/sohlich/go-plag/parser"
)

var (
	extensionRegex = regexp.MustCompile("^.*\\.")
)

//Process the submission that comes from api,
//unzip, parse and save to database.
func processSubmission(submission *Submission) error {

	filter, err := parser.GetLangFileFilter(submission.Lang)
	if err != nil {
		return err
	}
	sChan, err := unzipFile(submission.Content, filter)
	if err != nil {
		return err
	}

	//Generate new submission
	//Id for group of files
	//and save the files
	var submissionID string
	if len(submission.ID) > 0 {
		submissionID = submission.ID
	} else {
		submissionID = uuid.NewV1().String()
	}

	for submissionFile := range sChan {
		submissionFile.Submission = submissionID
		tokContent, tokMap, err := parser.TokenizeContent(submissionFile.Content, submission.Lang)
		if err != nil {
			Log.Errorf("Cannot process %s error: %s", submissionFile.Name, err)
			metrics.ErrorInc()
			continue
		}
		submissionFile.TokenMap = tokMap
		submissionFile.Tokens = tokContent
		submissionFile.Owner = submission.Owner
		submissionFile.Assignment = submission.AssignmentID
		mongo.Save(submissionFile)
	}
	return nil
}

//Unzipt the files in parallel
//returns buffered channel
func unzipFile(content []byte, filter map[string]bool) (<-chan *SubmissionFile, error) {
	r, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))

	if err != nil {
		Log.Error(err)
		return nil, err
	}

	files := r.File
	filesLen := len(files)
	unzippedStream := make(chan *SubmissionFile, filesLen)

	var wg sync.WaitGroup

	for _, f := range files {
		extension := parseExtension(f.Name)
		if filter[extension] {
			wg.Add(1)
			go func(file *zip.File) {
				processed, err := processFile(file)
				defer wg.Done()
				if err != nil {
					return
				}
				unzippedStream <- processed
			}(f)
		} else {
			Log.Infof("Not processing %s", f.Name)
		}
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
		Log.Error(rcError)
		return nil, rcError
	}
	content, rError := ioutil.ReadAll(rc)
	if rError != nil {
		Log.Error(rError)
		return nil, rError
	}
	submissionFile.Content = string(content)
	return submissionFile, nil
}

func parseExtension(fullPath string) string {
	return extensionRegex.ReplaceAllString(fullPath, "")
}
