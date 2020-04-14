package main

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/win0err/velobike-parser/parkings"
)

func initDB() *gorm.DB {
	db, err := gorm.Open("postgres", os.Getenv("DB_URI"))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&parkings.State{}, parkings.Station{})
	db.Model(&parkings.State{}).AddForeignKey("station_id", "stations(id)", "RESTRICT", "RESTRICT")

	return db
}

func main() {
	db := initDB()
	defer db.Close()

	stateRepository := parkings.ProvideStateRepostiory(db)

	response, _ := parkings.Get()
	states := parkings.ToStates(*response)

	stateRepository.SaveAll(states)
}
