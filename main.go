package main

import (
	"log"
	"time"

	"github.com/win0err/velobike-parser/database"
	"github.com/win0err/velobike-parser/helpers"
	"github.com/win0err/velobike-parser/parkings"
)

func init() {
	db, err := database.GetConnection()
	if err != nil {
		log.Fatalln("[FATAL] Unable get initial DB connection:", err)
	}
	defer db.Close()

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalln("[FATAL] Unable to migrate DB:", err)
	}
}

func main() {
	for {
		response, err := parkings.Get()
		if err != nil {
			log.Println("[ERROR] Unable to get parkings data:", err)

			log.Println("[INFO] Retry in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		go processResponse(response)
		helpers.SleepUntilNextMinute()
	}
}
