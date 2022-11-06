package gdynia

import (
	"fmt"
	"log"
	"time"
	wheretopark "wheretopark/go"
)

type Provider struct {
	configuration Configuration
	mapping       map[int]wheretopark.ID
}

func (p Provider) GetMetadata() (map[wheretopark.ID]wheretopark.Metadata, error) {
	vendorMetadata, err := GetMetadata()
	if err != nil {
		return nil, err
	}
	metadatas := make(map[wheretopark.ID]wheretopark.Metadata)
	for _, vendor := range vendorMetadata.Parkings {
		configuration, exists := p.configuration.ParkingLots[vendor.ID]
		if !exists {
			log.Printf("missing configuration for %d\n", vendor.ID)
			continue
		}
		id := wheretopark.GeometryToID(vendor.Location)
		metadata := wheretopark.Metadata{
			Name:           vendor.Name,
			Address:        vendor.Address,
			Geometry:       vendor.Location,
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
	for _, vendor := range *vendorState {
		id, exists := p.mapping[vendor.ParkingID]
		if !exists {
			log.Printf("no mapping for %d\n", vendor.ID)
			continue
		}

		fmt.Printf("vendorstate: %+v\n", vendor)

		lastUpdate, err := time.Parse("2006-01-02 15:04:05", vendor.InsertTime)
		if err != nil {
			return nil, err
		}
		state := wheretopark.State{
			LastUpdated: lastUpdate.Format(time.RFC3339),
			AvailableSpots: map[string]uint{
				"CAR": uint(vendor.FreePlaces),
			},
		}
		states[id] = state
	}
	return states, nil
}

func NewProvider(configurationPath *string) (*Provider, error) {
	var configuration Configuration
	if configurationPath == nil {
		configuration = DefaultConfiguration
	} else {
		newConfiguration, err := LoadConfiguration(*configurationPath)
		if err != nil {
			return nil, err
		}
		configuration = *newConfiguration
	}
	return &Provider{
		configuration: configuration,
		mapping:       make(map[int]wheretopark.ID),
	}, nil

}
