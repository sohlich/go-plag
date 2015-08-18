package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/validator.v2"
)

type DataStorage interface {
	OpenSession(url string) error
	CloseSession()
	Save(MongoObject) (interface{}, error)
	FindSubmissionFileById(id string) (*SubmissionFile, error)
	FindOneAssignmentById(id string) (*Assignment, error)
	FindAllSubmissionsByAssignment(assignmentId string) ([]SubmissionFile, error)
	FindAllComparableSubmissionFiles(submissionfile *SubmissionFile) ([]SubmissionFile, error)
}

type Mongo struct {
	Database             string `validate:"nonzero"`
	AssignmentCollection string `validate:"nonzero"`
	SubmissionCollection string `validate:"nonzero"`
	mongoSession         *mgo.Session
	assignments          *mgo.Collection
	submissions          *mgo.Collection
}

//Opens mongo session for given url and
//sets the session to global property
func (m *Mongo) OpenSession(url string) error {
	log.Info("")
	log.Infof(`Initializing MongoDB
		Connection string %s
		Database: %s
		Assignment collection: %s
		Submission collection: %s`,
		url,
		m.Database,
		m.AssignmentCollection,
		m.SubmissionCollection)

	validErr := validator.Validate(m)

	if validErr != nil {
		return &InitError{"Collection names Empty", validErr}
	}

	mongo, connError := mgo.Dial(url)
	m.mongoSession = mongo
	if connError != nil {
		return &InitError{"Connection failed", connError}
	}

	db := mongo.DB(m.Database)
	m.assignments = db.C(m.AssignmentCollection)
	m.submissions = db.C(m.SubmissionCollection)
	return nil
}

//Silently close mongo session
func (m *Mongo) CloseSession() {
	m.mongoSession.Close()
}

func (m *Mongo) Save(object MongoObject) (interface{}, error) {
	var err error

	object.NewId()

	if _, ok := object.(*Assignment); ok {
		err = m.assignments.Insert(object)
	} else if _, ok := object.(*SubmissionFile); ok {
		err = m.submissions.Insert(object)
	} else {
		err = errors.New("Object not assignable")
	}

	return object, err
}

func (m *Mongo) FindOneAssignmentById(id string) (*Assignment, error) {
	assignment := &Assignment{}
	err := m.assignments.FindId(bson.ObjectIdHex(id)).One(assignment)
	if err != nil {
		return nil, err
	}
	return assignment, err
}

func (m *Mongo) FindSubmissionFileById(id string) (*SubmissionFile, error) {
	sFile := &SubmissionFile{}
	err := m.submissions.FindId(bson.ObjectIdHex(id)).One(sFile)
	if err != nil {
		return nil, err
	}
	return sFile, err
}

func (m *Mongo) FindAllSubmissionsByAssignment(assignmentId string) ([]SubmissionFile, error) {
	var result []SubmissionFile
	err := m.submissions.Find(bson.M{"assignment": assignmentId}).All(&result)
	return result, err
}

func (m *Mongo) FindAllComparableSubmissionFiles(submissionF *SubmissionFile) ([]SubmissionFile, error) {
	var result []SubmissionFile
	err := m.submissions.Find(bson.M{"$and": []bson.M{bson.M{"assignment": submissionF.Assignment}, bson.M{"submission": bson.M{"$ne": submissionF.Submission}}}}).All(&result)
	return result, err
}
