package main

import (
	"errors"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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

	return nil, errors.New("No such file")
}

func (f *FakeDataStorage) FindAllSubmissionsByAssignment(assignmentId string) ([]SubmissionFile, error) {

	return []SubmissionFile{
		SubmissionFile{ID: bson.NewObjectId()},
		SubmissionFile{ID: bson.NewObjectId()},
		SubmissionFile{ID: bson.NewObjectId()},
		SubmissionFile{ID: bson.NewObjectId()},
		SubmissionFile{ID: bson.NewObjectId()},
	}, nil
}

func (f *FakeDataStorage) FindAllComparableSubmissionFiles(submissionfile *SubmissionFile) ([]SubmissionFile, error) {
	return []SubmissionFile{}, nil
}

func TestCompareFiles(t *testing.T) {
	log.Debugln("TestCompareFiles")
	oldMongo := mongo
	mongo = &FakeDataStorage{}
	defer func(s DataStorage) { mongo = s }(oldMongo)

	inputChan := make(chan OutputComparisonResult)
	outputChannel := compareFiles(inputChan)
	go func() { inputChan <- OutputComparisonResult{Files: []string{"1", "2"}} }()

	output := <-outputChannel

	assert.True(t, output.SimilarityIndex > 0)
	log.Debugf("Compared files wth index {}", output.SimilarityIndex)
}

func TestGenerateTuples(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	log.Debugln("TestCompareFiles")
	oldMongo := mongo
	mongo = &FakeDataStorage{}
	defer func(s DataStorage) { mongo = s }(oldMongo)

	checkAssignment(&Assignment{})
}
