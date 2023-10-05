package glasgow

import (
	"context"
	"fmt"
	"time"
	"wheretopark/collector/client"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/currency"

	geojson "github.com/paulmach/go.geojson"
)

type Source struct{}

var (
	// https://developer.glasgow.gov.uk/api-details#api=55c36a318b3a0306f0009483&operation=563cea91aab82f1168298575
	DATA_URL = wheretopark.MustParseURL("https://api.glasgow.gov.uk/datextraffic/carparks?format=json")
	API_KEY  = "ccaa1e24db6e4a9bb791f99433cc7ab7"
)

func (s Source) ParkingLots(ctx context.Context) (<-chan map[wheretopark.ID]wheretopark.ParkingLot, error) {
	vendor, err := client.Get[Response](DATA_URL, nil)
	if err != nil {
		return nil, err
	}

	parkingLots := make(map[wheretopark.ID]wheretopark.ParkingLot)
	for _, vendor := range vendor.LogicalModel.PayloadPublication.SituationItems {
		latitude, longitude :=
			vendor.Record.GroupOfLocations.LocationContainedInGroup.PointByCoordinates.PointByCoordinates.Latitude,
			vendor.Record.GroupOfLocations.LocationContainedInGroup.PointByCoordinates.PointByCoordinates.Longitude
		id := wheretopark.CoordinateToID(latitude, longitude)
		configuration, exists := configuration.ParkingLots[vendor.Record.ID]
		if !exists {
			log.Ctx(ctx).
				Warn().
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
			LastUpdated: lastUpdate.In(defaultTimezone),
			AvailableSpots: map[string]uint{
				wheretopark.SpotTypeCar: vendor.Record.TotalCapacity - uint(vendor.Record.OccupiedSpaces),
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
