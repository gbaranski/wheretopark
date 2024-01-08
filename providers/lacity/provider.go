package lacity

import (
	"context"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/currency"

	geojson "github.com/paulmach/go.geojson"
)

type Provider struct{}

var (
	METADATA_URL = wheretopark.MustParseURL("https://data.lacity.org/resource/s49e-q6j2.json")
	STATE_URL    = wheretopark.MustParseURL("https://data.lacity.org/resource/e7h6-4a3e.json")
)

var timezone *time.Location = wheretopark.MustLoadLocation("America/Los_Angeles")

func (p Provider) ParkingLots(ctx context.Context) (<-chan map[wheretopark.ID]wheretopark.ParkingLot, error) {
	vMetadata, err := wheretopark.Get[Metadata](METADATA_URL, nil)
	if err != nil {
		return nil, err
	}

	vState, err := wheretopark.Get[State](STATE_URL, nil)
	if err != nil {
		return nil, err
	}

	parkingLots := make(map[wheretopark.ID]wheretopark.ParkingLot, len(*vMetadata))
	for _, vSpaceMetadata := range *vMetadata {

		// metadata
		rules, err := rulesFor(vSpaceMetadata.MeteredTimeLimit, vSpaceMetadata.RateType, vSpaceMetadata.RateRange)
		if err != nil {
			log.Ctx(ctx).
				Warn().
				Err(err).
				Str("spaceID", vSpaceMetadata.SpaceID).
				Msg("failed to parse as rules")
			continue
		}
		metadata := wheretopark.Metadata{
			LastUpdated: nil,
			Name:        vSpaceMetadata.BlockFace,
			Address:     vSpaceMetadata.BlockFace,
			Geometry:    geojson.NewPointGeometry([]float64{vSpaceMetadata.Coordinate.Longitude, vSpaceMetadata.Coordinate.Latitude}),
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

		// state
		vSpaceStateIndex := statePositionBySpaceID(*vState, vSpaceMetadata.SpaceID)
		if vSpaceStateIndex == nil {
			continue
		}
		vSpaceState := (*vState)[*vSpaceStateIndex]
		lastUpdated, err := time.ParseInLocation("2006-01-02T15:04:05.000", vSpaceState.EventTime, time.UTC)
		if err != nil {
			log.Ctx(ctx).
				Warn().
				Err(err).
				Str("spaceID", vSpaceState.SpaceID).
				Msg("failed to parse last update time")
			continue
		}
		var availableSpots uint
		if vSpaceState.OccupancyState == OccupancyStateVacant {
			availableSpots = 1
		}
		state := wheretopark.State{
			LastUpdated: lastUpdated,
			AvailableSpots: map[string]uint{
				wheretopark.SpotTypeCar: availableSpots,
			},
		}

		// finish
		id := wheretopark.GeometryToID(metadata.Geometry)
		if _, exists := parkingLots[id]; exists {
			log.Ctx(ctx).
				Error().
				Str("blockFace", vSpaceMetadata.BlockFace).
				Str("id", id).
				Msg("parking lot duplicate")
		}
		parkingLot := wheretopark.ParkingLot{
			Metadata: metadata,
			State:    state,
		}
		parkingLots[id] = parkingLot
	}
	ch := make(chan map[wheretopark.ID]wheretopark.ParkingLot, 1)
	ch <- parkingLots
	close(ch)
	return ch, nil
}

func New() Provider {
	return Provider{}
}
