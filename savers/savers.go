package savers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/win0err/velobike-parser/database"
	"github.com/win0err/velobike-parser/parkings"
)

var BackupDir = os.Getenv("BACKUP_DIR")

func ToDb(states []parkings.State, currentTime time.Time) error {
	db, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("unable to connect database: %w", err)
	}
	defer db.Close()

	dbStateRepository := parkings.ProvideDbStateRepository(db)
	lastState, _ := dbStateRepository.GetLast()

	if !lastState.Time.Equal(currentTime) {
		if err := dbStateRepository.SaveAll(states); err != nil {
			return fmt.Errorf("error while saving data: %w", err)
		}
	} else {
		log.Printf("[INFO] Data already exists for %s, skipping...\n", currentTime)
	}

	return nil
}

func ToJson(states []parkings.State, currentTime time.Time) error {
	jsonStateRepository := parkings.ProvideJsonStateRepository(BackupDir)

	return jsonStateRepository.SaveAll(states, currentTime)
}
