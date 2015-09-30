package main

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
)

//TODO fix test
func TestCompareFiles(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	log.Debugln("TestCompareFiles")
	oldMongo := mongo
	mongo = &FakeDataStorage{}
	defer func(s DataStorage) { mongo = s }(oldMongo)

	inputChan := make(chan OutputComparisonResult)
	ctx, _ := context.WithCancel(context.TODO())
	outputChannel := compareFiles(ctx, inputChan)
	go func() {
		f1, _ := mongo.FindSubmissionFileById("1")
		f2, _ := mongo.FindSubmissionFileById("2")
		inputChan <- OutputComparisonResult{
			Files:  []string{"1", "2"},
			Tokens: []map[string]int{f1.TokenMap, f2.TokenMap},
		}
	}()

	output := <-outputChannel

	assert.True(t, output.SimilarityIndex > 0)
	log.Debugf("Compared files wth index {}", output.SimilarityIndex)
}

func TestGenerateTuples(t *testing.T) {
	log.Debugln("TestGenerateTuples")
	oldMongo := mongo
	mongo = &FakeDataStorage{}
	defer func(s DataStorage) { mongo = s }(oldMongo)

	files, _ := mongo.FindAllSubmissionsByAssignment("")

	ctx, cancel := context.WithCancel(context.Background())
	outChan := generateTuples(ctx, files)

	cancel()

	count := 0
	for out := range outChan {
		log.Infoln(out)
		count++
	}
	// assert.Equal(t, 10, count, "Did not comapred all files")
}

func TestCheckAssignmentPipeline(t *testing.T) {
	oldMongo := mongo
	mongo = &FakeDataStorage{}
	defer func(s DataStorage) { mongo = s }(oldMongo)
	assignment := &Assignment{ID: bson.NewObjectId()}
	compCount := checkAssignment(assignment)
	assert.True(t, compCount == nil, "Error in comparison")
}
