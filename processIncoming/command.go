package main

import (
	"github.com/hjkelly/pass-the-llama/models"
	"log"
)

func main() {
	listPage, err := models.FetchIncomingSmsPage()
	if err != nil {
		log.Printf("Failed to fetch incoming messages: " + err.Error())
	}

	// Work on the first page's items.
	numCheckins, numMisses, errs := models.RouteIncomingSmsPage(listPage)

	// So long as there's a next page, repeat that.
	for listPage.HasNextPage() {
		listPage, err = listPage.NextPage()
		addlCheckins, addlMisses, addlErrs := models.RouteIncomingSmsPage(listPage)
		numCheckins += addlCheckins
		numMisses += addlMisses
		for _, err := range addlErrs {
			errs = append(errs, err)
		}
	}

	log.Printf("Processed %d checkins, skipped %d incoming SMSes, and found %d errors.", numCheckins, numMisses, len(errs))
	for _, err := range errs {
		log.Println("    " + err.Error())
	}
}
