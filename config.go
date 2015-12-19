package main

import (
	"fmt"
)

//MongoConfig holds the configuration
//for connection to mongo database.
type MongoConfig struct {
	Port        string
	Host        string
	Database    string
	Assignments string
	Submissions string
	Results     string
}

//ConnectionString returns the connection
//string for mongodb database
func (config MongoConfig) ConnectionString() string {
	return fmt.Sprintf("%s:%s",
		config.Host,
		config.Port)
}

//ServerConfig holds
//the configuration
//for server binding
type ServerConfig struct {
	Port string
	Host string
}

//LogConfig holds
//the configuration
//for logger
type LogConfig struct {
	Path string
}

//ApacConfig holds the
//configuration for APAC
type ApacConfig struct {
	URL string
}

type configFile struct {
	Mongo  MongoConfig
	Server ServerConfig
	Log    LogConfig
	Apac   ApacConfig
}
