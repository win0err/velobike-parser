package database

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/win0err/velobike-parser/parkings"
)

var DbDialect = os.Getenv("DB_DIALECT")
var DbUri = os.Getenv("DB_URI")

func GetConnection() (*gorm.DB, error) {
	db, err := gorm.Open(DbDialect, DbUri)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	if !db.HasTable(&parkings.State{}) {
		if err := db.AutoMigrate(&parkings.State{}, parkings.Station{}).Error; err != nil {
			return err
		}

		db.Model(&parkings.State{}).
			AddUniqueIndex("idx_station_id_time", "station_id", "time")

		if DbDialect != "sqlite3" {
			return db.Model(&parkings.State{}).
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
