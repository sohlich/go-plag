package main

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/validator.v2"
)

//DataStorage provides the
//interface for data storing
//in application
type DataStorage interface {
	OpenSession() error                    //OpenSession opens the session to mongodb server
	CloseSession()                         //ClloseSession closes the session to mongodb server
	Save(MongoObject) (interface{}, error) //Save save objects by its interface
	FindSubmissionFileById(id string) (*SubmissionFile, error)
	FindOneAssignmentByID(id string) (*Assignment, error)
	FindAllSubmissionsByAssignment(assignmentID string) ([]SubmissionFile, error)
	FindAllComparableSubmissionFiles(submissionfile *SubmissionFile) ([]SubmissionFile, error)
	FindMaxSimilarityBySubmission(assignmentID string) ([]ApacPlagiarismSync, error)
}

//Mongo implements the storage
//at mongodb.
type Mongo struct {
	ConnectionString     string
	Database             string `validate:"nonzero"`
	AssignmentCollection string `validate:"nonzero"`
	SubmissionCollection string `validate:"nonzero"`
	ResultCollection     string `validate:"nonzero"`
	mongoSession         *mgo.Session
	assignments          *mgo.Collection
	submissions          *mgo.Collection
	results              *mgo.Collection
	db                   *mgo.Database
}

//OpenSession opens connection/session for mongo session
//for given url and
//sets the session to global property
func (m *Mongo) OpenSession() error {
	Log.Infof(`Initializing MongoDB
		Connection string: %s
		Database: %s
		Assignment collection: %s
		Submission collection: %s`,
		m.ConnectionString,
		m.Database,
		m.AssignmentCollection,
		m.SubmissionCollection)

	validErr := validator.Validate(m)

	if validErr != nil {
		return &InitError{"Collection names Empty", validErr}
	}

	mongo, connError := mgo.Dial(m.ConnectionString)
	m.mongoSession = mongo
	if connError != nil {
		return &InitError{"Connection failed", connError}
	}

	m.db = mongo.DB(m.Database)
	m.assignments = m.db.C(m.AssignmentCollection)
	m.submissions = m.db.C(m.SubmissionCollection)
	m.results = m.db.C(m.ResultCollection)

	//Indexes

	m.submissions.EnsureIndexKey("assignment")
	m.submissions.EnsureIndexKey("submission")
	m.results.EnsureIndexKey("fileId")
	m.results.EnsureIndexKey("compareTo")

	return nil
}

//CloseSession closes
//the connection to mongodb
func (m *Mongo) CloseSession() {
	m.mongoSession.Close()
}

//Save saves given MongoObject
//to mongodb database
func (m *Mongo) Save(object MongoObject) (interface{}, error) {
	var err error

	object.NewId()

	if _, ok := object.(*Assignment); ok {
		err = m.assignments.Insert(object)
	} else if _, ok := object.(*SubmissionFile); ok {
		err = m.submissions.Insert(object)
	} else if comparison, ok := object.(*OutputComparisonResult); ok {
		err = m.updateSimilarityResults(comparison)
	} else {
		err = errors.New("Object not assignable")
	}

	return object, err
}

//updateSimilarityResults Upserts the similarity results.
func (m *Mongo) updateSimilarityResults(comparison *OutputComparisonResult) error {

	fileQuery := bson.M{"fileId": comparison.Files[0],
		"comparedTo.fileId":     comparison.Files[1],
		"comparedTo.submission": comparison.Submissions[1],
		"submission":            comparison.Submissions[0],
		"assignment":            comparison.Assignment}
	fileUpdate := bson.M{"$set": bson.M{"similarity": comparison.SimilarityIndex}}
	_, err := m.results.Upsert(fileQuery, fileUpdate)

	ifErr(err)

	fileQuery["comparedTo.fileId"] = comparison.Files[0]
	fileQuery["comparedTo.submission"] = comparison.Submissions[0]
	fileQuery["fileId"] = comparison.Files[1]
	fileQuery["submission"] = comparison.Submissions[1]
	_, err = m.results.Upsert(fileQuery, fileUpdate)

	ifErr(err)

	return err
}

//FindOneAssignmentByID finds the
//assignment by given id
func (m *Mongo) FindOneAssignmentByID(id string) (*Assignment, error) {
	assignment := &Assignment{}
	err := m.assignments.FindId(bson.ObjectIdHex(id)).One(assignment)
	if err != nil {
		return nil, err
	}
	return assignment, err
}

//FindSubmissionFileById finds submission file
//by its ID.
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

//Extracts the similarity information for APAC
//sync.
func (m *Mongo) FindMaxSimilarityBySubmission(assignmentId string) ([]ApacPlagiarismSync, error) {
	query := []bson.M{
		{"$match": bson.M{"assignment": assignmentId}},
		{"$project": bson.M{"_id": 0, "uuid": "$submission", "uuid2": "$comparedTo.submission", "similarity": 1}},
		{"$group": bson.M{"_id": "$uuid", "similarity": bson.M{"$max": "$similarity"}, "submissions": bson.M{"$addToSet": bson.M{"uuid": "$uuid2", "similarity": "$similarity"}}}},
		{"$project": bson.M{"_id": 0, "baseuuid": "$_id", "similarity": 1, "submissions": 1}},
	}
	qryRes := make([]ApacPlagiarismSync, 0)
	err := m.results.Pipe(query).All(&qryRes)
	return qryRes, err
}

func ifErr(err error) {
	if err != nil {
		Log.Errorln(err)
	}
}
