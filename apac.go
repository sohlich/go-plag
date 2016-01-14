package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func syncWithApac(assignmentID string) error {
	result, err := mongo.FindMaxSimilarityBySubmission(assignmentID)

	if err != nil {
		return err
	}
	byteBody, err := json.Marshal(result)
	if err != nil {
		return err
	}

	Log.Infof(string(byteBody))

	reader := bytes.NewReader(byteBody)
	res, apacErr := http.Post(Apac.URL, "application/json", reader)
	if res != nil {
		Log.Infof("APAC response code %v", res.StatusCode)
		resBytes, _ := ioutil.ReadAll(res.Body)
		Log.Infof("Apac output is %s", string(resBytes))
	}
	return apacErr
}
