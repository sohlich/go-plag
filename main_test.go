package main

import (
	"testing"
)

func TestLoadProperties(t *testing.T) {
	result := loadProperties("test/plag.conf")

	assertProperties(t, result.Server.Host, "0.0.0.0")
	assertProperties(t, result.Server.Port, "8080")
	assertProperties(t, result.Mongo.ConnectionString(), "localhost:27017")
	assertProperties(t, result.Mongo.Assignments, "assignments")
	assertProperties(t, result.Mongo.Results, "results")
	assertProperties(t, result.Mongo.Submissions, "submissions")
	assertProperties(t, result.Log.Path, "/home/radek/tmp/log/plag.log")

	if result.Server.Port != "8080" {
		t.Error("Error in configuration reading")
	}
}

func assertProperties(t *testing.T, result, expected string) {
	if result != expected {
		t.Errorf("Error in configuration reading %s expected got %s", expected, result)
	}
}
