package main

import (
	"gopkg.in/mgo.v2/bson"
)

type Language string

const (
	JAVA Language = "java"
)

//
// const (
// 	NEW    string = "new"
// 	CLOSED string = "closed"
// 	SYNCED string = "synced"
// )

type MongoObject interface {
	NewId()
}

type Assignment struct {
	ID   bson.ObjectId `bson:"_id"`
	Name string
	Lang Language
	// State string
}

func (object *Assignment) NewId() {
	object.ID = bson.NewObjectId()
}

//Files from uploaded zipped file
//for submission
type SubmissionFile struct {
	ID         bson.ObjectId `bson:"_id"`
	Name       string
	Assignment string
	Submission string
	Similarity float32
	Content    string
	Tokens     []uint32
	TokenMap   map[string]int
}

func (object *SubmissionFile) NewId() {
	object.ID = bson.NewObjectId()
}

//DTO for http
type Submission struct {
	ID           string
	Owner        string
	AssignmentID string
	Lang         string
	Content      []byte //base64 file zip content
}

//DTO for fileComparison
type OutputComparisonResult struct {
	ID              bson.ObjectId `bson:"_id"`
	Files           []string
	SimilarityIndex float32
}

func (object *OutputComparisonResult) NewId() {
	object.ID = bson.NewObjectId()
}
