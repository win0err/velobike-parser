package parkings

import (
	"encoding/json"
	"fmt"
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

type Request struct {
	RawResponse []byte
	ParsedResponse *VelobikeResponse
}

func NewRequest() Request {
	return Request{
		ParsedResponse: &VelobikeResponse{},
	}
}

func (r *Request) Get() error {
	res, err := http.Get(VelobikeAPIEndpointURI)
	if err != nil {
		return err
	}

	r.RawResponse, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return nil
}

func (r *Request) Parse() error {
	if err := json.Unmarshal(r.RawResponse, r.ParsedResponse); err != nil {
		return err
	}

	if len(r.ParsedResponse.Items) == 0 {
		return fmt.Errorf("parkings data is empty: %+v", string(r.RawResponse))
	}

	r.ParsedResponse.Time = helpers.GetCurrentTime()
	return nil
}