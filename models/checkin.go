package models

import (
	"errors"
	"github.com/hjkelly/pass-the-llama/queries"
	"github.com/hjkelly/twiliogo"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type Checkin struct {
	TwilioSid  string    `bson:"twilioSid"`
	Status     string    `bson:"status"`
	ReceivedAt time.Time `bson:"receivedAt"`
}

func ParseCheckinFromTwilioMessage(m twiliogo.Message) (*Checkin, error) {
	// Is this even a status?
	status, err := getCheckinStatus(m.Body)
	if err != nil {
		return nil, err
	}

	// Get the timestamp for when it was received.
	timestamp, err := parseTwilioTimestamp(m.DateCreated)
	if err != nil {
		return nil, err
	}

	return &Checkin{
		TwilioSid:  m.Sid,
		Status:     status,
		ReceivedAt: timestamp,
	}, nil
}

func (c *Checkin) PushToAccount(phoneNumber string) error {
	return queries.UpdateAccount(bson.M{
		"phoneNumber":                phoneNumber,
		"archivedCheckins.twilioSid": bson.M{"$ne": c.TwilioSid},
		"newCheckins.twilioSid":      bson.M{"$ne": c.TwilioSid},
	}, bson.M{
		"$push": bson.M{"newCheckins": c},
	})
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
		return "", errors.New("Message body didn't appear to begin with a status: " + messageBody)
	}
}

func parseTwilioTimestamp(timestampStr string) (time.Time, error) {
	timestamp, err := time.Parse(time.RFC1123Z, timestampStr)
	if err != nil {
		err = errors.New("Couldn't parse message's DateCreated as RFC 1123Z: '" + timestampStr + "'")
	}
	return timestamp, err
}
