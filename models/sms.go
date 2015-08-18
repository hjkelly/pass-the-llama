package models

import (
	"errors"
	"fmt"
	"github.com/carlosdp/twiliogo"
	"os"
	"time"
)

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
	now := time.Now()
	dateSent := fmt.Sprintf("%4d-%2d-%2d", now.Year(), now.Month(), now.Day())
	listPage, err := twiliogo.GetMessageList(client, twiliogo.To(fromNumber), twiliogo.DateSent(dateSent))
	if err != nil {
		return nil, err
	} else {
		return listPage, nil
	}
}

func RouteIncomingSmsPage(listPage *twiliogo.MessageList) (int, int) {
	var checkins, misses int

	for _, message := range listPage.GetMessages() {
		error := RouteIncomingSms(message)
		if error != nil {
			misses += 1
		} else {
			checkins += 1
		}
	}

	return checkins, misses
}

func RouteIncomingSms(message twiliogo.Message) error {
	// Right now, something is either a checkin or it isn't. Consider it a
	// checkin if it starts with good, bad, or ugly.
	status, err := getCheckinStatus(message.Body)
	if err == nil {
		// We aren't tallying up errors just yet...
		err = AddCheckinForAccountPhone(message.From, Checkin{
			TwilioSid:        message.Sid,
			Status:           status,
			PartnersNotified: false,
			RoutedAt:         time.Now(),
		})
	}
	return err
}
