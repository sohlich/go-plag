package main

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
)

//Fake database connection
//and implementation
//of storage interface
type FakeDataStorage struct {
}

func (f *FakeDataStorage) OpenSession(url string) error {
	return nil
}

func (f *FakeDataStorage) CloseSession() {

}

func (m *FakeDataStorage) Save(object MongoObject) (interface{}, error) {
	return nil, nil
}

func (m *FakeDataStorage) FindOneAssignmentById(id string) (*Assignment, error) {
	return nil, nil
}

func (m *FakeDataStorage) FindSubmissionFileById(id string) (*SubmissionFile, error) {
	fileOne := &SubmissionFile{

		TokenMap: map[string]int{"1": 2, "5": 3},
	}

	fileTwo := &SubmissionFile{
		TokenMap: map[string]int{"1": 4, "5": 2},
	}

	switch id {
	case "1":
		return fileOne, nil
	case "2":
		return fileTwo, nil
	}

	return fileOne, nil
}

func (f *FakeDataStorage) FindAllSubmissionsByAssignment(assignmentId string) ([]SubmissionFile, error) {

	return []SubmissionFile{
		SubmissionFile{Submission: "1", ID: bson.NewObjectId()},
		SubmissionFile{Submission: "2", ID: bson.NewObjectId()},
		SubmissionFile{Submission: "3", ID: bson.NewObjectId()},
		SubmissionFile{Submission: "4", ID: bson.NewObjectId()},
		SubmissionFile{Submission: "5", ID: bson.NewObjectId()},
	}, nil
}

func (f *FakeDataStorage) FindAllComparableSubmissionFiles(submissionfile *SubmissionFile) ([]SubmissionFile, error) {
	return []SubmissionFile{}, nil
}

func TestCompareFiles(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	log.Debugln("TestCompareFiles")
	oldMongo := mongo
	mongo = &FakeDataStorage{}
	defer func(s DataStorage) { mongo = s }(oldMongo)

	inputChan := make(chan OutputComparisonResult)
	ctx, _ := context.WithCancel(context.TODO())
	outputChannel := compareFiles(ctx, inputChan)
	go func() { inputChan <- OutputComparisonResult{Files: []string{"1", "2"}} }()

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
	assert.True(t, compCount > 0, "No files were compared")
}
