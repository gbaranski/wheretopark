package krakow

import (
	"context"
	"time"
	wheretopark "wheretopark/go"
)

type Source struct {
	placemarks []Placemark
}

func (s Source) ParkingLots(ctx context.Context) (<-chan map[wheretopark.ID]wheretopark.ParkingLot, error) {
	parkingLots := make(map[wheretopark.ID]wheretopark.ParkingLot)
	for _, placemark := range s.placemarks {
		metadata := placemark.Metadata(0)
		id := wheretopark.GeometryToID(metadata.Geometry)
		state := wheretopark.State{
			LastUpdated: time.Now(),
			AvailableSpots: map[string]uint{
				wheretopark.SpotTypeCar: 0,
			},
		}
		parkingLots[id] = wheretopark.ParkingLot{
			Metadata: metadata,
			State:    state,
		}
	}
	ch := make(chan map[wheretopark.ID]wheretopark.ParkingLot, 1)
	ch <- parkingLots
	close(ch)
	return ch, nil
}

func New(placemarks []Placemark) Source {
	return Source{
		placemarks: placemarks,
	}
}
