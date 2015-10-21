package main

import (
	"fmt"
)

type MongoConfig struct {
	Port        string
	Host        string
	Database    string
	Assignments string
	Submissions string
	Results     string
}

func (config MongoConfig) ConnectionString() string {
	return fmt.Sprintf("%s:%s",
		config.Host,
		config.Port)
}

type ServerConfig struct {
	Port string
	Host string
}

type LogConfig struct {
	Path string
}

type ApacConfig struct {
	Url string
}

type configFile struct {
	Mongo  MongoConfig
	Server ServerConfig
	Log    LogConfig
	Apac   ApacConfig
}
