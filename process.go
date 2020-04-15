package main

import (
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/win0err/velobike-parser/helpers"
	"github.com/win0err/velobike-parser/parkings"
)

var DbDialect = os.Getenv("DB_DIALECT")
var DbUri = os.Getenv("DB_URI")
var BackupDir = os.Getenv("BACKUP_DIR")

func init() {
	db, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Migrate
	if !db.HasTable(&parkings.State{}) {
		db.AutoMigrate(&parkings.State{}, parkings.Station{})

		if "sqlite3" == DbDialect {
			db.Exec("PRAGMA foreign_keys = ON;")
		} else {
			db.Model(&parkings.State{}).
				AddForeignKey("station_id", "stations(id)", "RESTRICT", "RESTRICT")
		}
	}
}

func getDBConnection() (*gorm.DB, error) {
	db, err := gorm.Open(DbDialect, DbUri)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func processResponse(response *parkings.VelobikeResponse) {
	states := parkings.ToStates(*response)
	currentTime := helpers.GetCurrentTime()

	if err := saveToDb(states, currentTime); err == nil {
		log.Printf("Data successfully saved for %s\n", currentTime)
	} else {
		if err := saveAsJson(states, currentTime); err == nil {
			log.Printf("Data backuped for %s to %s\n", currentTime, BackupDir)
		} else {
			log.Println("Error while saving to DB:", err)
		}
	}
}

func saveToDb(states []parkings.State, currentTime time.Time) error {
	db, err := getDBConnection()
	if err != nil {
		log.Println("Unable to connect database:", err)

		return err
	}
	defer db.Close()

	dbStateRepository := parkings.ProvideDbStateRepository(db)
	lastState, _ := dbStateRepository.GetLast()

	if !lastState.Time.Equal(currentTime) {
		if err := dbStateRepository.SaveAll(states); err != nil {
			log.Println("Error while saving data:", err)
			return err
		}
	} else {
		log.Printf("Data already exists for %s, skipping...\n", currentTime)
	}

	return nil
}

func saveAsJson(states []parkings.State, currentTime time.Time) error {
	jsonStateRepository := parkings.ProvideJsonStateRepository(BackupDir)

	return jsonStateRepository.SaveAll(states, currentTime)
}
