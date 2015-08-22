package queries

import (
	"errors"
	"fmt"
	"github.com/hjkelly/twiliogo"
	"os"
	"time"
)

// Use the environment variables to pass back a Twilio client and the phone
// number (string) we're sending from.
func getTwilioClientAndFromNumber() (twiliogo.Client, string, error) {
	// Pull info from the environment vars.
	fromNumber := os.Getenv("TWILIO_FROM_NUMBER")
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	if len(accountSid) == 0 || len(authToken) == 0 || len(fromNumber) == 0 {
		return nil, "", errors.New("These three environment variables are required: TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, and TWILIO_FROM_NUMBER")
	}
	return twiliogo.NewClient(accountSid, authToken), fromNumber, nil
}

// Send a message to someone with a certain body. This will always be from our
// 'From' number.
func SendSms(toNumber, body string) error {
	// Get credentials and make sure we didn't have any problems initializing
	// the client.
	client, fromNumber, err := getTwilioClientAndFromNumber()
	if err != nil {
		return err
	}

	// Otherwise, go ahead and send the message!
	_, err = twiliogo.NewMessage(client, fromNumber, toNumber, twiliogo.Body(body))
	return err
}

func FetchIncomingSmsPage() (*twiliogo.MessageList, error) {
	// Get credentials and make sure we didn't have any problems initializing
	// the client.
	client, fromNumber, err := getTwilioClientAndFromNumber()
	if err != nil {
		return nil, err
	}

	// Query for new messages, but only get ones from today just for sanity's sake.
	yesterday := time.Now().Add(-24 * time.Hour)
	yesterdayStr := fmt.Sprintf("%4d-%02d-%02d", yesterday.Year(), yesterday.Month(), yesterday.Day())
	listPage, err := twiliogo.GetMessageList(client, twiliogo.To(fromNumber), twiliogo.DateSentAfter(yesterdayStr))
	if err != nil {
		return nil, err
	} else {
		return listPage, nil
	}
}
