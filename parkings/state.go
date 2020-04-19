package parkings

import "time"

// State contains information about the station at a specific point in time
type State struct {
	ID                  uint64    `gorm:"primary_key;auto_increment" json:"-"`
	Time                time.Time `gorm:"index:time_idx"`
	IsLocked            bool
	IsFavourite         bool
	FreeElectricPlaces  uint8
	FreeOrdinaryPlaces  uint8
	FreePlaces          uint8
	TotalElectricPlaces uint8
	TotalOrdinaryPlaces uint8
	TotalPlaces         uint8
	Station             Station `gorm:"foreignkey:StationID"`
	StationID           string  `gorm:"type:char(5)"`
}

// Station contains information about Velobike station
type Station struct {
	ID          string `gorm:"primary_key;index:station_idx;type:char(5)"`
	Name        string
	Address     string
	PositionLat float32
	PositionLon float32
	HasTerminal bool
}
