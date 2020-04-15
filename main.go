package main

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/win0err/velobike-parser/helpers"
	"github.com/win0err/velobike-parser/parkings"
)

func initDB() *gorm.DB {
	db, err := gorm.Open("postgres", os.Getenv("DB_URI"))
	if err != nil {
		panic(err)
	}

	return db
}

func init() {
	os.Setenv("TZ", "Europe/Moscow")

	db := initDB()
	defer db.Close()

	if !db.HasTable(&parkings.State{}) {
		db.
			AutoMigrate(&parkings.State{}, parkings.Station{}).
			Model(&parkings.State{}).
			AddForeignKey("station_id", "stations(id)", "RESTRICT", "RESTRICT")
	}
}

func main() {
	for {
		if response, err := parkings.Get(); err == nil {
			go processResponse(response)

			helpers.SleepUntilNextMinute()
		}
	}
}

func processResponse(response *parkings.VelobikeResponse) {
	db := initDB()
	defer db.Close()

	stateRepository := parkings.ProvideStateRepository(db)
	lastState := stateRepository.GetLast()

	alreadyParsed := lastState.Time.Equal(
		helpers.GetCurrentTime(),
	)

	if !alreadyParsed {
		states := parkings.ToStates(*response)
		stateRepository.SaveAll(states)
	}
}