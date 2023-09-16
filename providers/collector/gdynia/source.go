package gdynia

import (
	"context"
	"time"
	wheretopark "wheretopark/go"
	"wheretopark/providers/collector/client"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/currency"
)

var (
	METADATA_URL = wheretopark.MustParseURL("http://api.zdiz.gdynia.pl/ri/rest/parkings")
	STATE_URL    = wheretopark.MustParseURL("http://api.zdiz.gdynia.pl/ri/rest/parking_places")
)

type Source struct {
	mapping map[int]wheretopark.ID
}

func (s Source) Metadata(ctx context.Context) (map[wheretopark.ID]wheretopark.Metadata, error) {
	vendor, err := client.Get[Metadata](METADATA_URL, nil)
	if err != nil {
		return nil, err
	}

	metadatas := make(map[wheretopark.ID]wheretopark.Metadata)
	for _, vendor := range vendor.ParkingLots {
		configuration, exists := configuration.ParkingLots[vendor.ID]
		if !exists {
			log.Ctx(ctx).
				Warn().
				Int("id", vendor.ID).
				Str("name", vendor.Name).
				Msg("missing configuration")
			continue
		}
		id := wheretopark.GeometryToID(vendor.Location)
		metadata := wheretopark.Metadata{
			LastUpdated:    configuration.LastUpdated,
			Name:           vendor.Name,
			Address:        vendor.Address,
			Geometry:       vendor.Location,
			Resources:      configuration.Resources,
			TotalSpots:     configuration.TotalSpots,
			MaxDimensions:  configuration.MaxDimensions,
			Features:       configuration.Features,
			PaymentMethods: configuration.PaymentMethods,
			Comment:        configuration.Comment,
			Currency:       currency.PLN,
			Timezone:       defaultTimezone,
			Rules:          configuration.Rules,
		}
		metadatas[id] = metadata
		s.mapping[vendor.ID] = id
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
		id, exists := s.mapping[vendor.ParkingID]
		if !exists {
			log.Ctx(ctx).
				Debug().
				Int("id", vendor.ID).
				Msg("no mapping")
			continue
		}
		lastUpdate, err := time.ParseInLocation("2006-01-02 15:04:05", vendor.InsertTime, defaultTimezone)
		if err != nil {
			log.Ctx(ctx).
				Error().
				Err(err).
				Msg("failed to parse time")
			continue
		}
		state := wheretopark.State{
			LastUpdated: lastUpdate.UTC(),
			AvailableSpots: map[wheretopark.SpotType]uint{
				wheretopark.SpotTypeCar: uint(vendor.FreePlaces),
			},
		}
		states[id] = state
	}
	return states, nil
}

func New() Source {
	return Source{
		mapping: make(map[int]wheretopark.ID),
	}
}
