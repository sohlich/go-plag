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
	"github.com/rifflock/lfshook"
	"github.com/sohlich/go-plag/parser"
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
	comparison_count  = expvar.NewInt("comparison_count")
	submission_errors = expvar.NewInt("submission_error")
	cpuprofile        = flag.String("cpuprofile", "", "write cpu profile to file")
	Log               *log.Logger
)

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		//runtime.SetCPUProfileRate(100)
		err = pprof.StartCPUProfile(f)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	//Load config
	cfg := loadProperties(config)

	Log = NewLogger(cfg.Log.Path)

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

	//Load plugins
	loadPlugins()

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
	ginEngine.Use(logrusLogger()).PUT("/assignment", putAssignment)
	ginEngine.PUT("/submission", putSubmission)
	ginEngine.GET("/plugin/langs", getSupportedLangs)
	ginEngine.GET("/debug/vars", expvarGin.Handler())
}

func loadPlugins() {
	//Load plugins
	parser.LoadPlugins("plugin")
}

func logrusLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		Log.Infof("%s:%s from %s", c.Request.Method, c.Request.URL.String(), c.Request.RemoteAddr)
	}
}

func initStorage() {
	Log.Infoln("Initializing storage")
	err := mongo.OpenSession()
	if err != nil {
		Log.Fatal(err)
	}
}

func loadProperties(cfgFile string) configFile {
	var err error
	var cfg configFile
	if cfgFile != "" {
		err = gcfg.ReadFileInto(&cfg, cfgFile)
	}
	if err != nil {
		Log.Panic(err)
	}
	return cfg
}

//File logger for logrus
func NewLogger(path string) *log.Logger {
	if Log != nil {
		return Log
	}
	//This creates new logger
	Log = log.New()
	Log.Formatter = new(log.JSONFormatter)
	Log.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
		log.InfoLevel:  path,
		log.ErrorLevel: path,
		log.DebugLevel: path,
	}))
	return Log
}
