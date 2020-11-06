package database

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/win0err/velobike-parser/helpers"
	"github.com/win0err/velobike-parser/parkings"
)

var Connection *gorm.DB

func init() {
	var err error

	Connection, err = GetConnection()
	if err != nil {
		helpers.Log.Fatal("database connection failed:", err)
	}
}

func GetConnection() (*gorm.DB, error) {
	db, err := gorm.Open(helpers.Config.Database.Dialect, helpers.Config.Database.Uri)
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Hour)

	return db, nil
}

func AutoMigrate() error {
	if !Connection.HasTable(&parkings.State{}) {
		if err := Connection.AutoMigrate(&parkings.State{}, parkings.Station{}).Error; err != nil {
			return err
		}

		Connection.Model(&parkings.State{}).
			AddUniqueIndex("idx_station_id_time", "station_id", "time")

		if helpers.Config.Database.Dialect != "sqlite3" {
			return Connection.Model(&parkings.State{}).
				AddForeignKey(
					"station_id",
					"stations(id)",
					"RESTRICT",
					"RESTRICT",
				).Error
		}
	}

	return nil
}
