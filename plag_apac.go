package main

import (
	"encoding/json"
	"fmt"
)

func syncWithApac(assignmentId string) {
	result, err := mongo.FindMaxSimilarityBySubmission(assignmentId)
	if err != nil {
		return
	}
	bytes, err := json.Marshal(result)
	fmt.Println(string(bytes))
	// http.Post(url, bodyType, body)
}
