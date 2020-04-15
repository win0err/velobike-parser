package main

import (
	"log"
	"time"

	"github.com/win0err/velobike-parser/helpers"
	"github.com/win0err/velobike-parser/parkings"
)

func main() {
	for {
		response, err := parkings.Get()
		if err != nil {
			log.Println("Unable to get parkings data:", err)

			log.Println("Retry in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		go processResponse(response)
		helpers.SleepUntilNextMinute()
	}
}
