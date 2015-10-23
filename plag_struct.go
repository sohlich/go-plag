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

func (a *Assignment) Valid() bool {
	return len(a.Name) > 0 && len(a.Lang) > 0
}

//Files from uploaded zipped file
//for submission
type SubmissionFile struct {
	ID         bson.ObjectId `bson:"_id"`
	Owner      string        `omitempty`
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

func (s *Submission) Valid() bool {
	return len(s.AssignmentID) > 0
}

//DTO for fileComparison
type OutputComparisonResult struct {
	ID              bson.ObjectId `bson:"_id"`
	Assignment      string
	Files           []string
	Submissions     []string
	Tokens          []map[string]int `omitempty`
	SimilarityIndex float32
}

func (object *OutputComparisonResult) NewId() {
	object.ID = bson.NewObjectId()
}

//Embedded Structure for sending similarities
//to APAC system
type ApacSubmissionSimilarity struct {
	Uuid       string  `json:"uuid"`
	Similarity float64 `json:"similarity"`
}

//Main strucutre to send similarities
//to APAC system. This wraps one submission
//and extracts the max similarity from all
//comparisons of this submission
type ApacPlagiarismSync struct {
	Baseuuid    string                     `json:"baseuuid"`    //submission id
	Similarity  float64                    `json:"similarity"`  //max similarity
	Submissions []ApacSubmissionSimilarity `json:"submissions"` //other submissions
}
