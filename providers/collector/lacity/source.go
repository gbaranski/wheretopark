package lacity

import (
	"context"
	"time"
	wheretopark "wheretopark/go"
	"wheretopark/providers/collector/client"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/currency"

	geojson "github.com/paulmach/go.geojson"
)

type Source struct {
	// mapping space.ID -> wheretopark.ID
	mapping map[string]wheretopark.ID
}

var (
	METADATA_URL = wheretopark.MustParseURL("https://data.lacity.org/resource/s49e-q6j2.json")
	STATE_URL    = wheretopark.MustParseURL("https://data.lacity.org/resource/e7h6-4a3e.json")
)

var timezone *time.Location = wheretopark.MustLoadLocation("America/Los_Angeles")

func (s Source) Metadata(ctx context.Context) (map[wheretopark.ID]wheretopark.Metadata, error) {
	vendor, err := client.Get[Metadata](METADATA_URL, nil)
	if err != nil {
		return nil, err
	}

	metadatas := make(map[wheretopark.ID]wheretopark.Metadata, len(*vendor))
	for _, space := range *vendor {
		rules, err := rulesFor(space.MeteredTimeLimit, space.RateType, space.RateRange)
		if err != nil {
			log.Ctx(ctx).
				Warn().
				Err(err).
				Str("spaceID", space.SpaceID).
				Msg("failed to parse as rules")
			continue
		}
		metadata := wheretopark.Metadata{
			LastUpdated: nil,
			Name:        space.BlockFace,
			Address:     space.BlockFace,
			Geometry:    geojson.NewPointGeometry([]float64{space.Coordinate.Longitude, space.Coordinate.Latitude}),
			Resources: []string{
				"https://ladotparking.org",
			},
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: 1,
			},
			MaxDimensions: nil,
			Features: []wheretopark.Feature{
				wheretopark.FeatureUncovered,
			},
			PaymentMethods: []wheretopark.PaymentMethod{
				wheretopark.PaymentMethodCash, wheretopark.PaymentMethodCard,
			},
			Comment: map[string]string{
				"en": "Source of data: http://www.laexpresspark.org/",
			},
			Currency: currency.USD,
			Timezone: timezone,
			Rules:    rules,
		}
		id := wheretopark.GeometryToID(metadata.Geometry)
		if _, exists := metadatas[id]; exists {
			log.Ctx(ctx).
				Error().
				Str("blockFace", space.BlockFace).
				Str("id", id).
				Msg("metadata duplicate")
		}
		metadatas[id] = metadata
		s.mapping[space.SpaceID] = id
	}
	return metadatas, nil
}

func (s Source) State(ctx context.Context) (map[wheretopark.ID]wheretopark.State, error) {
	vendor, err := client.Get[State](STATE_URL, nil)
	if err != nil {
		return nil, err
	}

	states := make(map[wheretopark.ID]wheretopark.State)
	for _, vendor := range *vendor {
		id, exists := s.mapping[vendor.SpaceID]
		if !exists {
			log.Ctx(ctx).
				Debug().
				Str("spaceID", vendor.SpaceID).
				Msg("no mapping")
			continue
		}

		lastUpdated, err := time.ParseInLocation("2006-01-02T15:04:05.000", vendor.EventTime, time.UTC)
		if err != nil {
			log.Ctx(ctx).
				Warn().
				Err(err).
				Str("spaceID", vendor.SpaceID).
				Msg("failed to parse last update time")
			continue
		}

		var availableSpots uint
		if vendor.OccupancyState == OccupancyStateVacant {
			availableSpots = 1
		}
		states[id] = wheretopark.State{
			LastUpdated: lastUpdated,
			AvailableSpots: map[string]uint{
				wheretopark.SpotTypeCar: availableSpots,
			},
		}
	}

	return states, nil
}

func New() Source {
	return Source{
		mapping: make(map[string]wheretopark.ID),
	}
}
