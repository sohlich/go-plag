package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMongoInitValidationFailed(t *testing.T) {
	mongo := &Mongo{}
	err := mongo.OpenSession("url")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Collection names Empty")
	_, ok := err.(*InitError)
	if !ok {
		t.Error("Obtained error not type InitError")
	}
}

func TestMongoConnectionFailed(t *testing.T) {
	mongo := &Mongo{}
	err := mongo.OpenSession("url")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Collection names Empty")
	_, ok := err.(*InitError)
	if !ok {
		t.Error("Obtained error not type InitError")
	}
}

// func TestFindAllByAssignment(t *testing.T) {
// 	initStorage()
// 	result, err := mongo.FindAllSubmissionsByAssignment("55c7b6ebe13823356f000001")
// 	assert.Empty(t, err)
// 	assert.NotEmpty(t, result)
// }

// func TestFindAllBut(t *testing.T) {
// 	initStorage()

// 	submission := &SubmissionFile{
// 		Assignment: "55c7b6ebe13823356f000001",
// 		Submission: "3a66ab0d-4559-11e5-a728-f0def193326b",
// 	}

// 	result, err := mongo.FindAllComparableSubmissionFiles(submission)
// 	assert.Empty(t, err)
// 	assert.Equal(t, 9, len(result))
// }
