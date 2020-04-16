package main

import (
	"log"

	"github.com/win0err/velobike-parser/helpers"
	"github.com/win0err/velobike-parser/parkings"
	"github.com/win0err/velobike-parser/savers"
)

func processResponse(response *parkings.VelobikeResponse) {
	states := parkings.ToStates(*response)
	currentTime := helpers.GetCurrentTime()

	if err := savers.ToDb(states, currentTime); err == nil {
		log.Println("[INFO] Data successfully saved for", currentTime)
	} else {
		log.Println("[WARN] Unable to save to database:", err)
		if err := savers.ToJson(states, currentTime); err == nil {
			log.Printf("[INFO] Data backuped for %s\n", currentTime)
		} else {
			log.Println("[ERROR] Error while saving to DB:", err)
		}
	}

	wg.Done()
}
