package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	// "gopkg.in/mgo.v2/bson"
)

func putAssignment(ctx *gin.Context) {
	assignment := &Assignment{}
	ctx.BindJSON(assignment)

	log.Debugf("assignment object %s", assignment)

	_, err := mongo.Save(assignment)
	if err != nil {
		log.Error(err)
		ctx.JSON(405, "Object not stored")
		return
	}
	ctx.JSON(200, assignment)
}

//The submission is sent as mutipart form
func putSubmission(ctx *gin.Context) {
	meta := ctx.Request.FormValue("submission-meta")
	file, _, err := ctx.Request.FormFile("submission-data")
	if notifyError(err, ctx) {
		return
	}
	//Decode metadata
	decoder := json.NewDecoder(strings.NewReader(meta))
	submission := &Submission{}
	decoder.Decode(submission)

	log.Info(submission.AssignmentID)
	_, err = mongo.FindOneAssignment(submission.AssignmentID)
	if notifyError(err, ctx) {
		return
	}

	//Read content and append to entity
	fileContent, fError := ioutil.ReadAll(file)
	if notifyError(fError, ctx) {
		return
	}
	submission.Content = fileContent
	processSubmission(submission)
}

func notifyError(err error, ctx *gin.Context) bool {
	if err != nil {
		log.Error(err)
		ctx.JSON(405, err)
		return true
	}

	return false
}
