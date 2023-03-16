package glasgow

import (
	"fmt"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"

	geojson "github.com/paulmach/go.geojson"
)

var defaultLocation *time.Location

func init() {
	location, err := time.LoadLocation("Europe/London")
	if err != nil {
		panic(err)
	}
	defaultLocation = location
}

type Provider struct {
}

func (p Provider) Name() string {
	return "glasgow"
}

func (p Provider) GetMetadata() (map[wheretopark.ID]wheretopark.Metadata, error) {
	data, err := GetData()
	if err != nil {
		return nil, err
	}
	metadatas := make(map[wheretopark.ID]wheretopark.Metadata)
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
			Name:      configuration.Name,
			Address:   configuration.Address,
			Geometry:  *geojson.NewPointGeometry([]float64{longitude, latitude}),
			Resources: configuration.Resources,
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: vendor.Record.TotalCapacity,
			},
			MaxDimensions:  nil,
			Features:       configuration.Features,
			PaymentMethods: configuration.PaymentMethods,
			Comment:        configuration.Comment,
			Currency:       "GBP",
			Timezone:       defaultLocation.String(),
			Rules:          configuration.Rules,
		}
		metadatas[id] = metadata
	}
	return metadatas, nil
}

func (p Provider) GetState() (map[wheretopark.ID]wheretopark.State, error) {
	data, err := GetData()
	if err != nil {
		return nil, err
	}
	states := make(map[wheretopark.ID]wheretopark.State)

	for _, vendor := range data.LogicalModel.PayloadPublication.SituationItems {
		id := wheretopark.CoordinateToID(
			vendor.Record.GroupOfLocations.LocationContainedInGroup.PointByCoordinates.PointByCoordinates.Latitude,
			vendor.Record.GroupOfLocations.LocationContainedInGroup.PointByCoordinates.PointByCoordinates.Longitude,
		)

		lastUpdate, err := time.ParseInLocation("2006-01-02T15:04:05", vendor.Record.DateTime, defaultLocation)
		if err != nil {
			return nil, fmt.Errorf("failed to parse last update time: %w", err)
		}

		if vendor.Record.OccupiedSpaces < 0 || vendor.Record.TotalCapacity <= 0 {
			continue
		}
		state := wheretopark.State{
			LastUpdated: lastUpdate.In(defaultLocation).Format(time.RFC3339),
			AvailableSpots: map[string]uint{
				wheretopark.SpotTypeCar: vendor.Record.TotalCapacity - uint(vendor.Record.OccupiedSpaces),
			},
		}
		states[id] = state
	}
	return states, nil
}

func NewProvider() (wheretopark.Provider, error) {
	return Provider{}, nil

}
