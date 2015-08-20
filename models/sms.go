package models

import (
	"errors"
	"fmt"
	"github.com/hjkelly/twiliogo"
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
	yesterday := time.Now().Add(-24 * time.Hour)
	yesterdayStr := fmt.Sprintf("%4d-%02d-%02d", yesterday.Year(), yesterday.Month(), yesterday.Day())
	listPage, err := twiliogo.GetMessageList(client, twiliogo.To(fromNumber), twiliogo.DateSentAfter(yesterdayStr))
	if err != nil {
		return nil, err
	} else {
		return listPage, nil
	}
}

func RouteIncomingSmsPage(listPage *twiliogo.MessageList) (numCheckins int, numMisses int, errs []error) {
	for _, message := range listPage.GetMessages() {
		isCheckin, isMiss, err := RouteIncomingSms(message)
		if err != nil {
			errs = append(errs, err)
		} else if isMiss {
			numMisses += 1
		} else if isCheckin {
			numCheckins += 1
		} else {
			errs = append(errs, errors.New("RouteIncomingSms didn't report the outcome, nor did it report an error."))
		}
	}
	return
}

func RouteIncomingSms(message twiliogo.Message) (isCheckin bool, isMiss bool, err error) {
	// Right now, something is either a checkin or it isn't. Consider it a
	// checkin if it starts with good, bad, or ugly.
	status, statusErr := getCheckinStatus(message.Body)

	// If we couldn't figure out what to do with the message body, give up.
	if statusErr != nil {
		isMiss = true
		return
	} else {
		isCheckin = true
	}

	// Parse the DateCreated and use that as the timestamp.
	timestamp, err := time.Parse(time.RFC1123Z, message.DateCreated)
	if err != nil {
		err = errors.New("Couldn't parse message's DateCreated as RFC 1123Z: '" + message.DateCreated + "'")
		return
	}

	// We aren't tallying up errors just yet...
	err = AddCheckinForAccountPhone(message.From, Checkin{
		TwilioSid:        message.Sid,
		Status:           status,
		PartnersNotified: false,
		ReceivedAt:       timestamp,
		RoutedAt:         time.Now(),
	})
	return
}
