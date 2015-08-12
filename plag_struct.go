package main

import (
	"gopkg.in/mgo.v2/bson"
)

type Language string

const (
	JAVA Language = "java"
)

type MongoObject interface {
	NewId()
}

type Assignment struct {
	ID   bson.ObjectId `bson:"_id"`
	Name string
	Lang Language
}

func (object *Assignment) NewId() {
	object.ID = bson.NewObjectId()
}

//Files from uploaded zipped file
//for submission
type SubmissionFile struct {
	ID         bson.ObjectId `bson:"_id"`
	Name       string
	Submission string
	Similarity float32
	Content    string
}

func (object *SubmissionFile) NewId() {
	object.ID = bson.NewObjectId()
}

//DTO for http
type Submission struct {
	ID           string
	Owner        string
	AssignmentID string
	Content      []byte //base64 file zip content
}
