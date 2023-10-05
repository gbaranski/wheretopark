package gdynia

import (
	geojson "github.com/paulmach/go.geojson"
	"github.com/rs/zerolog/log"
)

type Metadata struct {
	ParkingLots []struct {
		ID             int               `json:"id"`
		Code           string            `json:"code"`
		Name           string            `json:"name"`
		Address        string            `json:"address"`
		StreetEntrance string            `json:"streetEntrance"`
		Location       *geojson.Geometry `json:"location"`
		LastUpdate     string            `json:"lastUpdate"`
	} `json:"parkings"`
}

type State = []struct {
	ID         int    `json:"id"`
	ParkingID  int    `json:"parkingId"`
	Capacity   int    `json:"capacity"`
	FreePlaces int    `json:"freePlaces"`
	InsertTime string `json:"insertTime"`
}

func StatePositionByID(state State, id int) int {
	for i, s := range state {
		if s.ID == id {
			return i
		}
	}
	log.Panic().Int("id", id).Msg("matching state not found")
	return 0
}
