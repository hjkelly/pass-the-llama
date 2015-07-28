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

/*
import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
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
	data := url.Values{}
	data.Set("From", fromNumber)
	data.Set("To", toNumber)
	data.Set("Body", body)
	req, _ := http.NewRequest("POST", "https://api.twilio.com/2010-04-01/Accounts/"+accountSid+"/Messages", bytes.NewBufferString(data.Encode()))
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New("Couldn't output body of " + strconv.Itoa(resp.StatusCode) + " response.")
		} else {
			return errors.New("Got " + strconv.Itoa(resp.StatusCode) + " response from Twilio API: " + string(respBodyBytes))
		}
	} else {
		return nil
	}
}
*/
