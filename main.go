package main

import (
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

//Todo mongo string to config
const (
	mgoConnString string = "localhost:27017"
)

var engine *gin.Engine = gin.Default()
var mongo DataStorage = &Mongo{
	Database:             "plag",
	AssignmentCollection: "assignments",
	SubmissionCollection: "submissions",
	ResultCollection:     "results",
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// log.SetLevel(log.DebugLevel)

	initGin(engine)
	initStorage()

	//Init db connection
	log.Info("Executing server")
	engine.Run("0.0.0.0:8080")
}

//Setup gin Engine server
func initGin(ginEngine *gin.Engine) {

	ginEngine.PUT("/assignment", putAssignment)
	ginEngine.PUT("/submission", putSubmission)
	ginEngine.Use(gin.Logger())
}

func initStorage() {
	err := mongo.OpenSession(mgoConnString)
	if err != nil {
		log.Fatal(err)
	}
}
