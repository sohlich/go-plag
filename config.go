package main

type MongoConfig struct {
	Port        string
	Host        string
	Database    string
	Assignments string
	Submissions string
	Results     string
}

type ServerConfig struct {
	Port string
	Host string
}

type configFile struct {
	Mongo  MongoConfig
	Server ServerConfig
}
