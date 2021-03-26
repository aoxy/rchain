package mongodb

import (
	"fmt"
	"github.com/globalsign/mgo"
	"os"
)

var session *mgo.Session

const MongoHost = "127.0.0.1"
const MongoPort = 27017
const MongoDb = "wallet"

func init() {
	var err error
	session, err = mgo.Dial(fmt.Sprintf("%v:%v", MongoHost, MongoPort))
	if err != nil {
		fmt.Println("mongodb connection failed ", err)
		os.Exit(2)
	}
}

func DB() *mgo.Database {
	return WithNameDB(MongoDb)
}

func WithNameDB(dbname string) *mgo.Database {
	return session.DB(dbname)
}

func Close() {
	if session != nil {
		session.Close()
	}
}
