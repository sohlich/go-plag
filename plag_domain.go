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
	Save(interface{}) (interface{}, error)
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
func (instance *Mongo) OpenSession(url string) error {
	log.Info("")
	log.Infof(`Initializing MongoDB
		Connection string %s
		Database: %s
		Assignment collection: %s
		Submission collection: %s`,
		url,
		instance.Database,
		instance.AssignmentCollection,
		instance.SubmissionCollection)

	validErr := validator.Validate(instance)

	if validErr != nil {
		return &InitError{"Collection names Empty", validErr}
	}

	mongo, connError := mgo.Dial(url)
	instance.mongoSession = mongo
	if connError != nil {
		return &InitError{"Connection failed", connError}
	}

	db := mongo.DB(instance.Database)
	instance.assignments = db.C(instance.AssignmentCollection)
	instance.submissions = db.C(instance.SubmissionCollection)
	return nil
}

//Silently close mongo session
func (instance *Mongo) CloseSession() {
	instance.mongoSession.Close()
}

func (instance *Mongo) Save(object MongoObject) (interface{}, error) {
	var err error

	object.NewId()

	if _, ok := object.(*Assignment); ok {
		err = instance.assignments.Insert(object)
	} else if _, ok := object.(*SubmissionFile); ok {
		err = instance.submissions.Insert(object)
	} else {
		err = errors.New("Object not assignable")
	}

	return object, err
}

func (instance *Mongo) FindOneAssignment(id string) (*Assignment, error) {
	assignment := &Assignment{}
	err := instance.assignments.FindId(bson.ObjectIdHex(id)).One(assignment)
	if err != nil {
		return nil, err
	}
	return assignment, err
}
