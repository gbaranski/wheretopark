package gdynia

import (
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
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
	mapping map[int]wheretopark.ID
}

func (p Provider) Name() string {
	return "gdynia"
}

func (p Provider) GetMetadata() (map[wheretopark.ID]wheretopark.Metadata, error) {
	vendorMetadata, err := GetMetadata()
	if err != nil {
		return nil, err
	}
	metadatas := make(map[wheretopark.ID]wheretopark.Metadata)
	for _, vendor := range vendorMetadata.Parkings {
		configuration, exists := configuration.ParkingLots[vendor.ID]
		if !exists {
			log.Warn().
				Int("id", vendor.ID).
				Str("name", vendor.Name).
				Msg("missing configuration")
			continue
		}
		id := wheretopark.GeometryToID(vendor.Location)
		metadata := wheretopark.Metadata{
			Name:           vendor.Name,
			Address:        vendor.Address,
			Geometry:       vendor.Location,
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
	states := make(map[wheretopark.ID]wheretopark.State)
	for _, vendor := range *vendorState {
		id, exists := p.mapping[vendor.ParkingID]
		if !exists {
			log.Warn().Int("id", vendor.ID).Msg("no mapping")
			continue
		}
		lastUpdate, err := time.ParseInLocation("2006-01-02 15:04:05", vendor.InsertTime, defaultLocation)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse time")
			continue
		}
		state := wheretopark.State{
			LastUpdated: lastUpdate.UTC().Format(time.RFC3339),
			AvailableSpots: map[wheretopark.SpotType]uint{
				wheretopark.SpotTypeCar: uint(vendor.FreePlaces),
			},
		}
		states[id] = state
	}
	return states, nil
}

func NewProvider() (wheretopark.Provider, error) {
	return Provider{
		mapping: make(map[int]wheretopark.ID),
	}, nil

}
