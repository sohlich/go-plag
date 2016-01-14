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

func TestJsonApacInput(t *testing.T) {
	input := "{\"identificator\":\"123455\"}"

	output := &Submission{}

	json.Unmarshal([]byte(input), output)

	if output.ID != "123455" {
		t.Error("Bad deserialization")
	}

}
