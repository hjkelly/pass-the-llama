package models

import (
	"errors"
	"github.com/hjkelly/pass-the-llama/queries"
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	Id               bson.ObjectId `bson:"_id"`
	Name             string        `bson:"name"`
	PhoneNumber      string        `bson:"phoneNumber"`
	Partners         []string      `bson:"partners"`
	PromptHour       int           `bson:"promptHour"`
	NewCheckins      []Checkin     `bson:"newCheckins"`
	ArchivedCheckins []Checkin     `bson:"archivedCheckins"`
}

func GetAccountByPhoneNumber(phoneNumber string) (*Account, error) {
	account := new(Account)
	err := queries.GetAccountByPhoneNumber(phoneNumber, account)
	return account, err
}

func ListAccountsNeedingPrompt(hour int) (*[]Account, error) {
	accounts := make([]Account, 0, 10)
	err := queries.ListAccountsNeedingPrompt(hour, &accounts)
	return &accounts, err
}

func ListAccountsWithNewCheckins() (*[]Account, error) {
	accounts := make([]Account, 0, 10)
	err := queries.ListAccountsWithNewCheckins(&accounts)
	return &accounts, err
}

func (a *Account) SendPrompt() error {
	// Pull info from the environment vars.
	err := queries.SendSms(a.PhoneNumber, "How have things been since your last checkin? Reply 'good', 'bad', or 'ugly'.")
	// If there's an error, give them some context.
	if err != nil {
		return errors.New("Sending text to " + a.PhoneNumber + "resulted in error: " + err.Error())
	} else {
		return nil
	}
}
