package parkings

import "time"

func ToState(item VelobikeResponseItem, time time.Time) State {
	return State{
		Time:                time,
		IsLocked:            item.IsLocked,
		IsFavourite:         item.IsFavourite,
		FreeElectricPlaces:  item.FreeElectricPlaces,
		FreeOrdinaryPlaces:  item.FreeOrdinaryPlaces,
		FreePlaces:          item.FreePlaces,
		TotalElectricPlaces: item.TotalElectricPlaces,
		TotalOrdinaryPlaces: item.TotalOrdinaryPlaces,
		TotalPlaces:         item.TotalPlaces,
		Station: Station{
			ID:          item.ID,
			Name:        item.Name,
			Address:     item.Address,
			PositionLat: item.Position.Lat,
			PositionLon: item.Position.Lon,
			HasTerminal: item.HasTerminal,
		},
	}
}

func ToStates(response VelobikeResponse) []State {
	states := make([]State, len(response.Items))

	for i, item := range response.Items {
		states[i] = ToState(item, response.Time)
	}

	return states
}
