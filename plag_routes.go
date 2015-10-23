package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sohlich/go-plag/parser"
)

func putAssignment(ctx *gin.Context) {
	assignment := &Assignment{}
	ctx.BindJSON(assignment)
	if !assignment.Valid() {
		ctx.JSON(405, "Invalid assignment definition")
	}
	assignment.Lang = Language(strings.ToLower(string(assignment.Lang)))
	Log.Debugf("assignment object %s", assignment)

	_, err := mongo.Save(assignment)
	if err != nil {
		Log.Error(err)
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

	//Start assembling
	submission := &Submission{}
	decoder.Decode(submission)

	Log.Info(submission.AssignmentID)

	if !submission.Valid() {
		Log.Errorf("Invalid submission data %v", submission)
		ctx.JSON(405, "Invalid submission data")
		return
	}

	assignment, mgoErr := mongo.FindOneAssignmentById(submission.AssignmentID)
	if notifyError(mgoErr, ctx) {
		return
	}
	submission.Lang = string(assignment.Lang)
	Log.Debug("Lang of submission is %v", submission.Lang)

	//Read content and append to entity
	fileContent, fError := ioutil.ReadAll(file)
	if notifyError(fError, ctx) {
		return
	}
	submission.Content = fileContent

	//Launch Goroutine to process submission
	go func(sub *Submission, assGnmnt *Assignment) {
		err := processSubmission(sub)
		if err == nil {
			err = checkAssignment(assGnmnt)
		}

		if err != nil {
			Log.Errorf("Error in processSubmission %s the error: %s", sub.ID, err.Error())
		}
	}(submission, assignment)
}

func getSupportedLangs(ctx *gin.Context) {
	Log.Infoln("Getting supported languages")
	ctx.JSON(200, parser.GetSupportedLangs())
}

func notifyError(err error, ctx *gin.Context) bool {
	if err != nil {
		Log.Error(err)
		ctx.JSON(405, err)
		return true
	}
	return false
}
