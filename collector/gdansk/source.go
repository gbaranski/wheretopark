package gdansk

import (
	"context"
	"fmt"
	"time"
	"wheretopark/collector/client"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"

	geojson "github.com/paulmach/go.geojson"
)

var (
	METADATA_URL = wheretopark.MustParseURL("https://ckan.multimediagdansk.pl/dataset/cb1e2708-aec1-4b21-9c8c-db2626ae31a6/resource/d361dff3-202b-402d-92a5-445d8ba6fd7f/download/parking-lots.jso")
	STATE_URL    = wheretopark.MustParseURL("https://ckan2.multimediagdansk.pl/parkingLots")
)

type Source struct{}

func (s Source) ParkingLots(ctx context.Context) (<-chan map[wheretopark.ID]wheretopark.ParkingLot, error) {
	vMetadata, err := client.Get[Metadata](METADATA_URL, nil)
	if err != nil {
		return nil, err
	}
	vState, err := client.Get[State](STATE_URL, nil)
	if err != nil {
		return nil, err
	}

	lastUpdated, err := time.Parse(time.RFC3339, vState.LastUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse last update time: %w", err)
	}
	lastUpdated = lastUpdated.In(defaultTimezone)

	parkingLots := make(map[wheretopark.ID]wheretopark.ParkingLot)
	for _, vMetadata := range vMetadata.ParkingLots {
		configuration, exists := configuration.ParkingLots[vMetadata.ID]
		if !exists {
			log.Ctx(ctx).
				Warn().
				Str("id", vMetadata.ID).
				Str("name", vMetadata.Name).
				Msg("missing configuration")
			continue
		}
		id := wheretopark.CoordinateToID(vMetadata.Location.Latitude, vMetadata.Location.Longitude)
		metadata := wheretopark.Metadata{
			LastUpdated:    configuration.LastUpdated,
			Name:           vMetadata.Name,
			Address:        vMetadata.Address,
			Geometry:       geojson.NewPointGeometry([]float64{vMetadata.Location.Longitude, vMetadata.Location.Latitude}),
			Resources:      configuration.Resources,
			TotalSpots:     configuration.TotalSpots,
			MaxDimensions:  configuration.MaxDimensions,
			Features:       configuration.Features,
			PaymentMethods: configuration.PaymentMethods,
			Comment:        configuration.Comment,
			Currency:       configuration.Currency,
			Timezone:       configuration.Timezone,
			Rules:          configuration.Rules,
		}

		state := wheretopark.State{
			LastUpdated: lastUpdated,
			AvailableSpots: map[wheretopark.SpotType]uint{
				wheretopark.SpotTypeCar: vState.AvailableSpotsByID(vMetadata.ID),
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
