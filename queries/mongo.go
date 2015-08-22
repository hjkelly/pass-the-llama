package queries

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

var accountCollection = "accounts"

func getDb() *mgo.Database {
	url := "localhost"
	database := "passTheLlama"

	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error when calling `mgo.Dial('"+url+"')`")
		panic(err)
	}

	// TODO: add a deferred Close()

	return session.DB(database)
}

func GetAccountByPhoneNumber(phoneNumber string, result interface{}) error {
	return getDb().C(accountCollection).Find(bson.M{
		"phoneNumber": phoneNumber,
	}).One(result)
}

func ListAccountsNeedingPrompt(hour int, result interface{}) error {
	return getDb().C(accountCollection).Find(bson.M{
		"promptHour": hour,
	}).All(result)
}

func ListAccountsWithNewCheckins(result interface{}) error {
	return getDb().C(accountCollection).Find(bson.M{
		"newCheckins": bson.M{"$size": bson.M{"$gt": 0}},
	}).All(result)
}

func UpdateAccount(query, changes bson.M) error {
	return getDb().C(accountCollection).Update(query, changes)
}

func UpdateManyAccounts(query, changes bson.M) error {
	_, err := getDb().C(accountCollection).UpdateAll(query, changes)
	return err
}

/*



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

func ListAccountsWithNewCheckins() (*[]Account, error) {
	db := getDb()

	accounts := make([]Account, 0, 10)
	err := db.C(accountCollection).Find(bson.M{
		"newCheckins": bson.M{"$size": bson.M{"gt": 0}},
	}).All(&accounts)
	return &accounts, err
}

func AddCheckinForAccountPhone(accountPhoneNumber string, c Checkin) error {
	db := getDb()
	_, err := db.C(accountCollection).UpdateAll(bson.M{
		"phoneNumber":                accountPhoneNumber,
		"newCheckins.twilioSid":      bson.M{"$ne": c.TwilioSid},
		"archivedCheckins.twilioSid": bson.M{"$ne": c.TwilioSid},
	}, bson.M{
		"$push": bson.M{"newCheckins": c},
	})
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
*/
