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
	models.RouteIncomingSmsPage(listPage)

	// So long as there's a next page, repeat that.
	for listPage.HasNextPage() {
		listPage, err = listPage.NextPage()
		models.RouteIncomingSmsPage(listPage)
	}
}
