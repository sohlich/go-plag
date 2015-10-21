package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func syncWithApac(assignmentId string) error {
	result, err := mongo.FindMaxSimilarityBySubmission(assignmentId)
	if err != nil {
		return err
	}
	byteBody, err := json.Marshal(result)
	if err != nil {
		return err
	}
	Log.Infof("Syncing to APAC %v", string(byteBody))
	reader := bytes.NewReader(byteBody)
	res, apacErr := http.Post(Apac.Url, "application/json", reader)
	if res != nil {
		Log.Infof("APAC response code %v", res.StatusCode)
		resBytes, _ := ioutil.ReadAll(res.Body)
		Log.Infof("Apac output is %s", string(resBytes))
	}
	return apacErr
}
