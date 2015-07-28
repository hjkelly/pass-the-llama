package models

import (
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

// MODEL / LOGIC LEVEL: **********

type Account struct {
	Id          bson.ObjectId `bson:"_id"`
	PhoneNumber string        `bson:"phoneNumber"`
	Partners    []string      `bson:"partners"`
	PromptHour  int           `bson:"promptHour"`
	LastCheckin *time.Time    `bson:"lastCheckin"`
}

func (a *Account) SendPrompt() error {
	// Pull info from the environment vars.
	err := SendSms(a.PhoneNumber, "How have things been since your last checkin? Reply 'good', 'bad', or 'ugly'.")
	if err == nil {
		now := time.Now()
		a.LastCheckin = &now
		a.Save()
	}
	return err
}

// DB / QUERY LEVEL: **********

var accountCollection = "accounts"

func GetAccountsNeedingPrompt(hour int) *[]Account {
	db := getDb()
	accounts := make([]Account, 0)
	query := bson.M{
		"promptHour": hour,
	}
	err := db.C(accountCollection).Find(query).All(&accounts)
	if err != nil {
		if err.Error() != "not found" {
			panic(err)
		} else {
			log.Printf("Found no accounts using query %+v", query)
		}
	}
	return &accounts
}

func (a *Account) Save() {
	db := getDb()
	err := db.C(accountCollection).Update(bson.M{"_id": a.Id}, a)
	if err != nil {
		log.Printf("Couldn't update account:\n%+v\nReason:\n%s", a, err.Error())
	}
}
