package main

import (
	"gopkg.in/mgo.v2/bson"
)

//Fake database connection
//and implementation
//of storage interface
type FakeDataStorage struct {
}

func (f *FakeDataStorage) OpenSession() error {
	return nil
}

func (f *FakeDataStorage) CloseSession() {
}

func (m *FakeDataStorage) Save(object MongoObject) (interface{}, error) {
	return nil, nil
}

func (m *FakeDataStorage) FindOneAssignmentByID(id string) (*Assignment, error) {
	return &Assignment{
		Lang: Language("java"),
	}, nil
}

func (m *FakeDataStorage) FindSubmissionFileByID(id string) (*SubmissionFile, error) {
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

func (f *FakeDataStorage) FindMaxSimilarityBySubmission(assignmentId string) ([]ApacPlagiarismSync, error) {
	result := []ApacPlagiarismSync{
		ApacPlagiarismSync{
			Baseuuid:   "10",
			Similarity: 0.23,
			Submissions: []ApacSubmissionSimilarity{
				ApacSubmissionSimilarity{
					UUID:       "12",
					Similarity: 0.23,
				},
			},
		},
	}
	return result, nil
}
