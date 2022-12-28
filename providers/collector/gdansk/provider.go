package gdansk

import (
	"fmt"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"

	geojson "github.com/paulmach/go.geojson"
)

var defaultLocation *time.Location

func init() {
	location, err := time.LoadLocation("Europe/Warsaw")
	if err != nil {
		panic(err)
	}
	defaultLocation = location
}

type Provider struct {
	mapping map[string]wheretopark.ID
}

func (p Provider) GetMetadata() (map[wheretopark.ID]wheretopark.Metadata, error) {
	vendorMetadata, err := GetMetadata()
	if err != nil {
		return nil, err
	}
	metadatas := make(map[wheretopark.ID]wheretopark.Metadata)
	for _, vendor := range vendorMetadata.ParkingLots {
		configuration, exists := configuration.ParkingLots[vendor.ID]
		if !exists {
			log.Warn().Str("id", vendor.ID).Msg("missing configuration")
			continue
		}
		id := wheretopark.CoordinateToID(vendor.Location.Latitude, vendor.Location.Longitude)
		metadata := wheretopark.Metadata{
			Name:           vendor.Name,
			Address:        vendor.Address,
			Geometry:       *geojson.NewPointGeometry([]float64{vendor.Location.Longitude, vendor.Location.Latitude}),
			Resources:      configuration.Resources,
			TotalSpots:     configuration.TotalSpots,
			MaxDimensions:  configuration.MaxDimensions,
			Features:       configuration.Features,
			PaymentMethods: configuration.PaymentMethods,
			Comment:        configuration.Comment,
			Currency:       "PLN",
			Timezone:       "Europe/Warsaw",
			Rules:          configuration.Rules,
		}
		metadatas[id] = metadata
		p.mapping[vendor.ID] = id
	}
	return metadatas, nil
}

func (p Provider) GetState() (map[wheretopark.ID]wheretopark.State, error) {
	vendorState, err := GetState()
	if err != nil {
		return nil, err
	}
	lastUpdate, err := time.Parse(time.RFC3339, vendorState.LastUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse last update time: %w", err)
	}
	states := make(map[wheretopark.ID]wheretopark.State)

	for _, vendor := range vendorState.ParkingLots {
		id, exists := p.mapping[vendor.ID]
		if !exists {
			log.Warn().Str("id", vendor.ID).Msg("no mapping")
			continue
		}

		state := wheretopark.State{
			LastUpdated: lastUpdate.In(defaultLocation).Format(time.RFC3339),
			AvailableSpots: map[wheretopark.SpotType]uint{
				wheretopark.SpotTypeCar: vendor.AvailableSpots,
			},
		}
		states[id] = state
	}
	return states, nil
}

func NewProvider() (wheretopark.Provider, error) {
	return Provider{
		mapping: make(map[string]wheretopark.ID),
	}, nil
}
