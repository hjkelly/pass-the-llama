package main

import (
	"fmt"
	"github.com/hjkelly/pass-the-llama/models"
	"github.com/hjkelly/pass-the-llama/queries"
	"github.com/hjkelly/twiliogo"
	"gopkg.in/mgo.v2/bson"
	"os"
)

func SaveNewCheckins() (numSaved int, numNotSaved int, err error) {
	page := new(twiliogo.MessageList)

	// Load the first page.
	page, err = queries.FetchIncomingSmsPage()

	// Handle the first and the last pages the same way.
	for {
		// If we hit an error, give up.
		if err != nil {
			break
		}

		// Process each incoming message.
		for _, message := range page.GetMessages() {
			c, err := models.ParseCheckinFromTwilioMessage(message)
			// If there was NOT an error, save the new checkin and tally it up.
			if err == nil {
				nonFatalErr := c.PushToAccount(message.From)
				if nonFatalErr != nil {
					numNotSaved++
				} else {
					numSaved++
				}
			}
		}

		// If there's no next page, give up.
		if page.HasNextPage() == false {
			break
		} else {
			// Otherwise, download the new page.
			page, err = page.NextPage()
		}
	}
	return
}

func NotifyPartnersOfNewCheckins() (numNotificationsSent int, err error) {
	var accountsWithNewCheckins *[]models.Account
	// Which accounts need to have notifications sent?
	accountsWithNewCheckins, err = models.ListAccountsWithNewCheckins()
	if err != nil {
		return
	}

	// Loop through each and send notifications.
	for _, a := range *accountsWithNewCheckins {
		// Of all the new checkins, find the most recent.
		mostRecentCheckin := new(models.Checkin)
		for _, newCheckin := range a.NewCheckins {
			if mostRecentCheckin == nil || newCheckin.ReceivedAt.After(mostRecentCheckin.ReceivedAt) {
				*mostRecentCheckin = newCheckin
			}
		}

		// Form the notification body.
		body := fmt.Sprintf("%s says things have been %s.", a.Name, mostRecentCheckin.Status)
		// Send it to each partner.
		for _, partnerPhoneNumber := range a.Partners {
			err = queries.SendSms(partnerPhoneNumber, body)
			if err != nil {
				return
			} else {
				numNotificationsSent++
			}
		}

		// Move the relevant checkins to the archived list.
		err = queries.UpdateAccount(
			bson.M{"_id": a.Id},
			bson.M{
				"$push":  bson.M{"archivedCheckins": bson.M{"$each": a.NewCheckins}},
				"$unset": bson.M{"newCheckins": true},
			},
		)
		if err != nil {
			return
		}
	}

	return
}

func main() {
	numSaved, numNotSaved, err := SaveNewCheckins()
	if err != nil {
		fmt.Fprintf(os.Stdout, "SaveNewCheckins: %d new checkins, %d not saved, err: %s\n", numSaved, numNotSaved, err.Error())
		os.Exit(1)
	} else {
		fmt.Fprintf(os.Stdout, "SaveNewCheckins: %d new checkins, %d not saved.\n", numSaved, numNotSaved)
	}

	// Notify partners of new checkins.
	numNotificationsSent, err := NotifyPartnersOfNewCheckins()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		fmt.Fprintf(os.Stdout, "Sent %d notifications.\n", numNotificationsSent)
		os.Exit(1)
	} else {
		fmt.Fprintf(os.Stdout, "Sent %d notifications.\n", numNotificationsSent)
	}
}
