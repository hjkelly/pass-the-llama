package main

import (
	"github.com/hjkelly/pass-the-llama/models"
	"log"
	"time"
)

func main() {
	accounts := models.GetAccountsNeedingPrompt(time.Now().Hour())

	successes := 0
	failures := 0
	for _, a := range *accounts {
		err := a.SendPrompt()
		if err != nil {
			failures += 1
			log.Printf("Couldn't send a message because: %+v", err)
		} else {
			successes += 1
		}
	}
	log.Printf("Successfully sent prompts for %d accounts, but another %d failed.", successes, failures)
}
