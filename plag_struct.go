package main

import (
	"gopkg.in/mgo.v2/bson"
)

//Language defines the
//assignment tyoe
type Language string

//MongoObject is common interface
//for handling Ids in structs
type MongoObject interface {
	NewID()
}

//Assignment is a strucutre
//that holds one assignmen
//in DataStorage. The name is
//human readable identifier
//Lang defines the programming
//language that the assignment is
//for
type Assignment struct {
	ID   bson.ObjectId `bson:"_id"`
	Name string
	Lang Language
}

//NewID generates new bson.ObjectId
//for assignment.
func (a *Assignment) NewID() {
	a.ID = bson.NewObjectId()
}

//Valid returns id the assignemtn
// has valid attributes.
func (a *Assignment) Valid() bool {
	return len(a.Name) > 0 && len(a.Lang) > 0
}

//SubmissionFile holds the files
//from uploaded zipped file
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

//NewID generates new bson.ObjectId
//for SubmisisonFile.
func (sf *SubmissionFile) NewID() {
	sf.ID = bson.NewObjectId()
}

//Submission is dto for REST
type Submission struct {
	ID           string
	Owner        string
	AssignmentID string
	Lang         string
	Content      []byte
}

//Valid returns if the structure has
//valid properties.
func (s *Submission) Valid() bool {
	return len(s.AssignmentID) > 0
}

//OutputComparisonResult is output dto for
//file comparison process.
type OutputComparisonResult struct {
	ID              bson.ObjectId `bson:"_id"`
	Assignment      string
	Files           []string
	Submissions     []string
	Tokens          []map[string]int `omitempty`
	SimilarityIndex float32
}

//NewID generates new bson.ObjectId
//for struct.
func (ocr *OutputComparisonResult) NewID() {
	ocr.ID = bson.NewObjectId()
}

//ApacSubmissionSimilarity is embedded struct
//for sending similarities
//to APAC system
type ApacSubmissionSimilarity struct {
	UUID       string  `json:"uuid"`
	Similarity float64 `json:"similarity"`
}

//ApacPlagiarismSync is strucutre to send similarities
//to APAC system. This wraps one submission
//and extracts the max similarity from all
//comparisons of this submission
type ApacPlagiarismSync struct {
	Baseuuid    string                     `json:"baseUuid"`           //submission id
	Similarity  float64                    `json:"similarity"`         //max similarity
	Submissions []ApacSubmissionSimilarity `json:"similarSubmissions"` //other submissions
}
