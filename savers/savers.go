package savers

import (
	"fmt"
	"time"

	"github.com/win0err/velobike-parser/database"
	"github.com/win0err/velobike-parser/helpers"
	"github.com/win0err/velobike-parser/parkings"
)

func ToDb(states []parkings.State) error {
	currentTime := states[0].Time.Truncate(time.Second)

	dbStateRepository := parkings.ProvideDbStateRepository(database.Connection)
	lastState, _ := dbStateRepository.GetLast()

	if !lastState.Time.Truncate(time.Second).Equal(currentTime) {
		if err := dbStateRepository.SaveAll(states); err != nil {
			return fmt.Errorf("error while saving data: %w", err)
		}
	} else {
		helpers.Log.Info("data already exists for %s, skipping...\n", currentTime)
	}

	return nil
}

func ToJson(states []parkings.State) error {
	currentTime := states[0].Time
	jsonStateRepository := parkings.ProvideJsonStateRepository(helpers.Config.BackupDir)

	return jsonStateRepository.SaveAll(states, currentTime)
}
