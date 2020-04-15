package parkings

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/win0err/velobike-parser/helpers"
)

// VelobikeResponse structure is properly parsed JSON with request time added.
type VelobikeResponse struct {
	Items []VelobikeResponseItem
	Time  time.Time `json:"-"`
}

// VelobikeResponseItem structure is properly parsed JSON item.
type VelobikeResponseItem struct {
	Address            string
	FreeElectricPlaces uint8
	FreeOrdinaryPlaces uint8
	FreePlaces         uint8
	HasTerminal        bool
	ID                 string
	IsFavourite        bool
	IsLocked           bool
	Name               string
	Position           struct {
		Lat float32
		Lon float32
	}
	StationTypes        []string
	TotalElectricPlaces uint8
	TotalOrdinaryPlaces uint8
	TotalPlaces         uint8
}

const VelobikeAPIEndpointURI = "https://velobike.ru/ajax/parkings/"

// Get information about parkings from Velobike's API and parse response.
func Get() (*VelobikeResponse, error) {
	res, err := http.Get(VelobikeAPIEndpointURI)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	responseData := &VelobikeResponse{}
	json.Unmarshal(body, &responseData)

	responseData.Time = helpers.GetCurrentTime()

	return responseData, nil
}
