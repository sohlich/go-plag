package main

import (
	"testing"
)

func TestApacSync(t *testing.T) {
	oldMongo := mongo
	mongo = &FakeDataStorage{}
	defer func(s DataStorage) { mongo = s }(oldMongo)
	syncWithApac("assignmentId")
}
