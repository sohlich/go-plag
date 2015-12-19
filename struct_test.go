package main

import (
	"encoding/json"
	"log"
	"strings"
	"testing"
)

func TestJsonApac(t *testing.T) {
	apacDto := ApacPlagiarismSync{
		Baseuuid:   "123",
		Similarity: 0.24,
		Submissions: []ApacSubmissionSimilarity{
			ApacSubmissionSimilarity{
				UUID:       "1234",
				Similarity: 1,
			}},
	}
	output, err := json.Marshal(apacDto)
	if err != nil {
		t.Error("Json serialization failed")
	}
	stringOutput := string(output)
	log.Printf("APAC Serialization outout\n%s", stringOutput)
	if !strings.Contains(stringOutput, "baseUuid") {
		t.Error("Bad serialization output")
	}
}
