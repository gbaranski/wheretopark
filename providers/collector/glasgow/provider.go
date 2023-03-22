package glasgow

import (
	"fmt"
	"time"
	wheretopark "wheretopark/go"
	"wheretopark/go/provider/simple"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/currency"

	geojson "github.com/paulmach/go.geojson"
)

type Provider struct{}

func (p Provider) Config() simple.Config {
	return simple.NewConfig(time.Hour * 1)
}

func (p Provider) Name() string {
	return "glasgow"
}

func (p Provider) GetParkingLots() (map[wheretopark.ID]wheretopark.ParkingLot, error) {
	data, err := GetData()
	if err != nil {
		return nil, err
	}
	parkingLots := make(map[wheretopark.ID]wheretopark.ParkingLot)
	for _, vendor := range data.LogicalModel.PayloadPublication.SituationItems {

		latitude, longitude :=
			vendor.Record.GroupOfLocations.LocationContainedInGroup.PointByCoordinates.PointByCoordinates.Latitude,
			vendor.Record.GroupOfLocations.LocationContainedInGroup.PointByCoordinates.PointByCoordinates.Longitude
		id := wheretopark.CoordinateToID(latitude, longitude)
		configuration, exists := configuration.ParkingLots[vendor.Record.ID]
		if !exists {
			log.Warn().Str("name", vendor.Record.ID).Msg("missing configuration")
			log.Warn().
				Str("id", vendor.Record.ID).
				Str("name", vendor.Record.Identity).
				Msg("missing configuration")
			continue
		}
		metadata := wheretopark.Metadata{
			LastUpdated:    configuration.LastUpdated,
			Name:           configuration.Name,
			Address:        configuration.Address,
			Geometry:       geojson.NewPointGeometry([]float64{longitude, latitude}),
			Resources:      configuration.Resources,
			TotalSpots:     map[string]uint{wheretopark.SpotTypeCar: vendor.Record.TotalCapacity},
			MaxDimensions:  nil,
			Features:       configuration.Features,
			PaymentMethods: configuration.PaymentMethods,
			Comment:        configuration.Comment,
			Currency:       currency.GBP,
			Timezone:       defaultTimezone,
			Rules:          configuration.Rules,
		}

		lastUpdate, err := time.ParseInLocation("2006-01-02T15:04:05", vendor.Record.DateTime, defaultTimezone)
		if err != nil {
			return nil, fmt.Errorf("failed to parse last update time: %w", err)
		}

		if vendor.Record.OccupiedSpaces < 0 || vendor.Record.TotalCapacity <= 0 {
			continue
		}
		state := wheretopark.State{
			LastUpdated: lastUpdate.In(defaultTimezone).Format(time.RFC3339),
			AvailableSpots: map[string]uint{
				wheretopark.SpotTypeCar: vendor.Record.TotalCapacity - uint(vendor.Record.OccupiedSpaces),
			},
		}
		parkingLots[id] = wheretopark.ParkingLot{
			Metadata: metadata,
			State:    state,
		}
	}
	return parkingLots, nil
}

func NewProvider() (simple.Provider, error) {
	return Provider{}, nil
}
