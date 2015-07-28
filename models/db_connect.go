package models

import (
	"gopkg.in/mgo.v2"
	"log"
)

func getDb() *mgo.Database {
	url := "localhost"
	database := "passTheLlama"

	session, err := mgo.Dial(url)
	if err != nil {
		log.Println("Error when calling `mgo.Dial('"+url+"')`")
		panic(err)
	}

	return session.DB(database)
}

