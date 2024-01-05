package krakow

import (
	"context"
	"time"
	wheretopark "wheretopark/go"
)

var (
	METADATA_URL = wheretopark.MustParseURL("http://zdmk.krakow.pl/wp-content/themes/justidea_theme/assets/xml/parkomaty.xml")
)

type Source struct{}

func (s Source) ParkingLots(ctx context.Context) (<-chan map[wheretopark.ID]wheretopark.ParkingLot, error) {
	vMetadata, err := wheretopark.Get[Metadata](METADATA_URL, nil)
	if err != nil {
		return nil, err
	}

	parkingLots := make(map[wheretopark.ID]wheretopark.ParkingLot)
	for _, placemark := range vMetadata.Placemarks {
		metadata := placemark.Metadata(10)
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

func New() Source {
	return Source{}
}
