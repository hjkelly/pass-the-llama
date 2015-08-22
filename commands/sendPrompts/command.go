package main

import (
	"fmt"
	"github.com/hjkelly/pass-the-llama/models"
	"os"
	"time"
)

func SendHourlyPrompts() (int, error) {
	// Get the accounts that need this.
	accounts, err := models.ListAccountsNeedingPrompt(time.Now().Hour())
	if err != nil {
		return 0, err
	}

	// Send the prompt for each.
	for idx, a := range *accounts {
		err = a.SendPrompt()
		if err != nil {
			return idx, err
		}
	}

	// Report that all went well.
	return len(*accounts), nil
}

func main() {
	numSent, err := SendHourlyPrompts()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		fmt.Fprintf(os.Stdout, "Sent %d checkin prompts.", numSent)
		os.Exit(1)
	} else {
		fmt.Fprintf(os.Stdout, "Sent %d checkin prompts.", numSent)
		os.Exit(0)
	}
}
