package main

import (
	// 	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
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

//The submission is sent as json with assignment identifier and
//zipped content in base64
func putSubmission(ctx *gin.Context) {
	submission := &Submission{}
	ctx.BindJSON(submission)

	log.Debugf("new submission for assignment %s", submission.AssignmentID)

}
