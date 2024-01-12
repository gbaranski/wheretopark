package krakow

import (
	"context"
	"time"
	wheretopark "wheretopark/go"
	"wheretopark/go/timeseries"
)

type Provider struct {
	placemarks []Placemark
	timeseries timeseries.TimeSeries
}

func (p Provider) ParkingLots(ctx context.Context) (<-chan map[wheretopark.ID]wheretopark.ParkingLot, error) {
	parkingLots := make(map[wheretopark.ID]wheretopark.ParkingLot)
	for _, placemark := range p.placemarks {
		id := placemark.ID()
		totalSpots := p.timeseries.MaxOccupancyOf(id)
		metadata := placemark.Metadata(totalSpots)
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

func New(placemarks []Placemark, timeseries timeseries.TimeSeries) Provider {
	return Provider{
		placemarks,
		timeseries,
	}
}
