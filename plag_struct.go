package main

import (
	"gopkg.in/mgo.v2/bson"
)

type Language string

const (
	JAVA Language = "java"
)

type MongoObject interface {
	NewId()
}

type Assignment struct {
	ID   bson.ObjectId `bson:"_id"`
	Name string
	Lang Language
}

func (object *Assignment) NewId() {
	object.ID = bson.NewObjectId()
}
