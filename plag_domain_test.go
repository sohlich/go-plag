package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMongoInitValidationFailed(t *testing.T) {
	mongo := &Mongo{}
	err := mongo.OpenSession("url")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Collection names Empty")
	_, ok := err.(*InitError)
	if !ok {
		t.Error("Obtained error not type InitError")
	}
}

func TestMongoConnectionFailed(t *testing.T) {
	mongo := &Mongo{}
	err := mongo.OpenSession("url")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Collection names Empty")
	_, ok := err.(*InitError)
	if !ok {
		t.Error("Obtained error not type InitError")
	}
}
