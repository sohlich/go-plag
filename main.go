package main

import (
	"expvar"
	"flag"
	"fmt"
	"os"
	"runtime/pprof"

	"code.google.com/p/gcfg"
	log "github.com/Sirupsen/logrus"
	expvarGin "github.com/gin-gonic/contrib/expvar"
	"github.com/gin-gonic/gin"
)

const (
	config = "application.conf"
)

var (
	engine             = gin.Default()
	mongo  DataStorage = &Mongo{
		Database:             "plag",
		AssignmentCollection: "assignments",
		SubmissionCollection: "submissions",
		ResultCollection:     "results",
	}

	//expvar
	comparison_count = expvar.NewInt("comparison_count")
	cpuprofile       = flag.String("cpuprofile", "", "write cpu profile to file")
)

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	//Load config
	cfg := loadProperties(config)

	//Setup and init storage
	mgoConf := cfg.Mongo
	mgoConnString := fmt.Sprintf("%s:%s",
		mgoConf.Host,
		mgoConf.Port)
	mgo, _ := mongo.(*Mongo)
	mgo.AssignmentCollection = mgoConf.Assignments
	mgo.ResultCollection = mgoConf.Results
	mgo.SubmissionCollection = mgoConf.Submissions
	mgo.Database = mgoConf.Database
	mgo.ConnectionString = mgoConnString
	initStorage()
	defer mgo.CloseSession()

	//Setup and start webserver
	webConf := cfg.Server
	initGin(engine)
	address := fmt.Sprintf("%s:%s",
		webConf.Host,
		webConf.Port)
	engine.Run(address)
}

//Setup gin Engine server
func initGin(ginEngine *gin.Engine) {
	ginEngine.PUT("/assignment", putAssignment)
	ginEngine.PUT("/submission", putSubmission)
	ginEngine.Use(logrusLogger()).GET("/debug/vars", expvarGin.Handler())
}

func logrusLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		comparison_count.Add(1)
	}
}

func initStorage() {
	err := mongo.OpenSession()
	if err != nil {
		log.Fatal(err)
	}
}

func loadProperties(cfgFile string) configFile {
	var err error
	var cfg configFile
	if cfgFile != "" {
		err = gcfg.ReadFileInto(&cfg, cfgFile)
	}
	if err != nil {
		log.Panic(err)
	}
	return cfg
}
