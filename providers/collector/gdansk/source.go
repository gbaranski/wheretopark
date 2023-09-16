package gdansk

import (
	"context"
	"fmt"
	"time"
	wheretopark "wheretopark/go"
	"wheretopark/providers/collector/client"

	"github.com/rs/zerolog/log"

	geojson "github.com/paulmach/go.geojson"
)

var (
	METADATA_URL = wheretopark.MustParseURL("https://ckan.multimediagdansk.pl/dataset/cb1e2708-aec1-4b21-9c8c-db2626ae31a6/resource/d361dff3-202b-402d-92a5-445d8ba6fd7f/download/parking-lots.jso")
	STATE_URL    = wheretopark.MustParseURL("https://ckan2.multimediagdansk.pl/parkingLots")
)

type Source struct {
	mapping map[string]wheretopark.ID
}

func (s Source) Metadata(ctx context.Context) (map[wheretopark.ID]wheretopark.Metadata, error) {
	vendor, err := client.Get[Metadata](METADATA_URL, nil)
	if err != nil {
		return nil, err
	}

	metadatas := make(map[wheretopark.ID]wheretopark.Metadata)
	for _, vMetadata := range vendor.ParkingLots {
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
		metadatas[id] = metadata
		s.mapping[vMetadata.ID] = id
	}
	return metadatas, nil
}

func (s Source) State(ctx context.Context) (map[wheretopark.ID]wheretopark.State, error) {
	vendor, err := client.Get[State](STATE_URL, nil)
	if err != nil {
		return nil, err
	}

	lastUpdate, err := time.Parse(time.RFC3339, vendor.LastUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse last update time: %w", err)
	}
	states := make(map[wheretopark.ID]wheretopark.State)
	for _, vState := range vendor.ParkingLots {
		id, exists := s.mapping[vState.ID]
		if !exists {
			log.Ctx(ctx).
				Debug().
				Str("id", vState.ID).
				Msg("no mapping")
			continue
		}

		state := wheretopark.State{
			LastUpdated: lastUpdate.In(defaultTimezone),
			AvailableSpots: map[wheretopark.SpotType]uint{
				wheretopark.SpotTypeCar: vState.AvailableSpots,
			},
		}
		states[id] = state
	}
	return states, nil
}

func New() Source {
	return Source{
		mapping: make(map[string]wheretopark.ID),
	}
}
