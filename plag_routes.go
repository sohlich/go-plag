package main

import (
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
