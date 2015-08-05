package models

import (
	"errors"
	"github.com/carlosdp/twiliogo"
	"os"
)

func SendSms(toNumber, body string) error {
	// Pull info from the environment vars.
	fromNumber := os.Getenv("TWILIO_FROM_NUMBER")
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	if len(accountSid) == 0 || len(authToken) == 0 || len(fromNumber) == 0 {
		return errors.New("These three environment variables are required: TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, and TWILIO_FROM_NUMBER")
	}

	// Otherwise, go ahead and send the message!
	client := twiliogo.NewClient(accountSid, authToken)
	_, err := twiliogo.NewMessage(client, fromNumber, toNumber, twiliogo.Body(body))
	return err
}
