package test

import (
	"github.com/globalsign/mgo"
	"rchain/db/mongodb"
	"testing"
)

func TestCreate(t *testing.T) {
	c := mongodb.DB().C("test")
	c.Create(&mgo.CollectionInfo{})
	c.Insert(map[string]string{"name": "Tom", "age": "30"})
}
