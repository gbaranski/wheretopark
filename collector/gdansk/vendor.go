package gdansk

import (
	"github.com/rs/zerolog/log"
)

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Metadata struct {
	LastUpdate  string `json:"lastUpdate"`
	ParkingLots []struct {
		ID             string     `json:"id"`
		Name           string     `json:"name"`
		ShortName      string     `json:"shortName"`
		Address        string     `json:"address"`
		StreetEntrance string     `json:"streetEntrance"`
		Location       Coordinate `json:"location"`
	} `json:"parkingLots"`
}

type State struct {
	LastUpdate  string `json:"lastUpdate"`
	ParkingLots []struct {
		ID             string `json:"parkingId"`
		AvailableSpots uint   `json:"availableSpots"`
		LastUpdate     string `json:"lastUpdate"`
	} `json:"parkingLots"`
}

func (s *State) AvailableSpotsByID(id string) uint {
	for _, parkingLot := range s.ParkingLots {
		if parkingLot.ID == id {
			return parkingLot.AvailableSpots
		}
	}
	log.Panic().Str("id", id).Msg("parking lot not found")
	return 0
}
