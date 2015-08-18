package models

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
	"time"
)

// MODEL / LOGIC LEVEL: **********

type Account struct {
	Id          bson.ObjectId `bson:"_id"`
	PhoneNumber string        `bson:"phoneNumber"`
	Partners    []string      `bson:"partners"`
	PromptHour  int           `bson:"promptHour"`
	Checkins    []Checkin     `bson:"checkins"`
}

type Checkin struct {
	TwilioSid        string    `bson:"twilioSid"`
	Status           string    `bson:"status"`
	PartnersNotified bool      `bson:"partnersNotified"`
	RoutedAt         time.Time `bson:"routedAt"`
}

func (a *Account) SendPrompt() error {
	// Pull info from the environment vars.
	return SendSms(a.PhoneNumber, "How have things been since your last checkin? Reply 'good', 'bad', or 'ugly'.")
}

func getCheckinStatus(messageBody string) (string, error) {
	lowercaseBody := strings.ToLower(messageBody)
	if strings.HasPrefix(lowercaseBody, "good") {
		return "good", nil
	} else if strings.HasPrefix(lowercaseBody, "bad") {
		return "bad", nil
	} else if strings.HasPrefix(lowercaseBody, "ugly") {
		return "ugly", nil
	} else {
		return "", errors.New("no valid status found")
	}
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

func AddCheckinForAccountPhone(accountPhoneNumber string, c Checkin) error {
	db := getDb()
	err := db.C(accountCollection).Update(bson.M{"phoneNumber": accountPhoneNumber}, bson.M{"$push": bson.M{"checkins": c}})
	if err != nil {
		log.Printf("Couldn't update account with phone number:\n%+v\nReason:\n%s", accountPhoneNumber, err.Error())
	}
	return err
}

func (a *Account) Save() {
	db := getDb()
	err := db.C(accountCollection).Update(bson.M{"_id": a.Id}, a)
	if err != nil {
		log.Printf("Couldn't update account:\n%+v\nReason:\n%s", a, err.Error())
	}
}
