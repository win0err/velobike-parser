package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/win0err/velobike-parser/database"
	"github.com/win0err/velobike-parser/helpers"
	"github.com/win0err/velobike-parser/parkings"
	"github.com/win0err/velobike-parser/savers"
)

func processResponse(states []parkings.State) {
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

func importData(reader io.Reader) error {
	var states []parkings.State

	if err := json.NewDecoder(reader).Decode(&states); err != nil {
		if err != io.EOF {
			return fmt.Errorf("unable to parse json: %w", err)
		}
	}

	wg.Add(1)

	db, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("unable to connect database: %w", err)
	}
	defer db.Close()

	savedCount := 0
	dbStateRepository := parkings.ProvideDbStateRepository(db)

	for _, state := range states {
		if err := dbStateRepository.Save(state); err != nil {
			if !strings.HasPrefix(err.Error(), "UNIQUE") {
				log.Printf(
					"[INFO] %s\n",
					fmt.Errorf("unable to save state: %w", err),
				)
			}
		} else {
			savedCount += 1
		}
	}

	log.Printf("[INFO] %d of %d states has been imported\n", savedCount, len(states))
	wg.Done()

	return nil
}

func exportData(all bool, to, from string) ([]byte, error) {
	var err error
	var states []parkings.State

	wg.Add(1)

	db, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("unable to connect database: %w", err)
	}
	defer db.Close()

	dbStateRepository := parkings.ProvideDbStateRepository(db)

	if all {
		if states, err = dbStateRepository.FindAll(); err != nil {
			return nil, err
		}
	} else {
		var fromTime, toTime time.Time

		fromTime, err = time.Parse("2006-01-02 15:04 MST", from)
		if err != nil {
			return nil, fmt.Errorf("unable to parse date: %w", err)
		}
		toTime, err = time.Parse("2006-01-02 15:04 MST", to)
		if err != nil {
			return nil, fmt.Errorf("unable to parse date: %w", err)
		}

		if states, err = dbStateRepository.FindByTimeRange(fromTime, toTime); err != nil {
			return nil, err
		}
	}

	data, err := json.Marshal(states)

	wg.Done()

	return data, err
}
