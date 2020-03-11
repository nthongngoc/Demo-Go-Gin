package shared

import (
	"log"
	"time"

	"github.com/khoa5773/go-server/src/configs"
	mgo "gopkg.in/mgo.v2"
)

var MongoSession *mgo.Database

func init() {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{configs.ConfigsService.DbHost},
		Timeout:  60 * time.Second,
		Database: configs.ConfigsService.DbName,
		Username: configs.ConfigsService.DbUser,
		Password: configs.ConfigsService.DbPass,
	}

	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}
	mongoSession.SetMode(mgo.Monotonic, true)

	MongoSession = mongoSession.DB("")
}
