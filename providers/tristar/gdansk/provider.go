package gdansk

import (
	"errors"
	wheretopark "wheretopark/go"

	geojson "github.com/paulmach/go.geojson"
)

type Provider struct {
	configuration Configuration
	mapping       map[string]wheretopark.ID
}

func (p Provider) GetMetadata() (map[wheretopark.ID]wheretopark.Metadata, error) {
	vendorMetadata, err := GetMetadata()
	if err != nil {
		return nil, err
	}
	metadatas := make(map[wheretopark.ID]wheretopark.Metadata)
	for _, vendor := range vendorMetadata.ParkingLots {
		configuration := p.configuration.ParkingLots[vendor.ID]
		id := wheretopark.CoordinateToID(vendor.Location.Latitude, vendor.Location.Longitude)
		metadata := wheretopark.Metadata{
			Name:           vendor.Name,
			Address:        vendor.Address,
			Location:       *geojson.NewPointFeature([]float64{vendor.Location.Longitude, vendor.Location.Latitude}),
			Resources:      configuration.Resources,
			TotalSpots:     configuration.TotalSpots,
			MaxWidth:       configuration.MaxWidth,
			MaxHeight:      configuration.MaxHeight,
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
	for _, vendor := range vendorState.ParkingLots {
		id, exists := p.mapping[vendor.ID]
		if !exists {
			return nil, errors.New("no mapping for " + vendor.ID)
		}

		state := wheretopark.State{
			LastUpdated: vendor.LastUpdate,
			AvailableSpots: map[string]uint{
				"CAR": vendor.AvailableSpots,
			},
		}
		states[id] = state
	}
	return states, nil
}

func NewProvider(configurationPath string) (*Provider, error) {
	configuration, err := LoadConfiguration(configurationPath)
	if err != nil {
		return nil, err
	}
	return &Provider{
		configuration: *configuration,
		mapping:       make(map[string]wheretopark.ID),
	}, nil

}
